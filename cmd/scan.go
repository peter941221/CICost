package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/peter941221/CICost/internal/auth"
	"github.com/peter941221/CICost/internal/config"
	gh "github.com/peter941221/CICost/internal/github"
	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/store"
)

func runScan(args []string) error {
	ctx, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "Target repository in owner/repo format")
	daysFlag := fs.Int("days", ctx.cfg.Scan.Days, "Time window in days")
	incrementalFlag := fs.Bool("incremental", ctx.cfg.Scan.Incremental, "Enable incremental sync")
	fullFlag := fs.Bool("full", false, "Force full sync")
	workersFlag := fs.Int("workers", ctx.cfg.Scan.Workers, "Concurrent worker count")
	tokenFlag := fs.String("token", "", "GitHub token (optional)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	repo, err := pickRepo(*repoFlag, ctx.cfg)
	if err != nil {
		return err
	}
	owner, repoName, err := splitRepo(repo)
	if err != nil {
		return err
	}
	if *workersFlag <= 0 {
		*workersFlag = 4
	}
	if *workersFlag > 8 {
		*workersFlag = 8
	}

	token, err := auth.ResolveToken(*tokenFlag, ctx.cfg.Auth.Token)
	if err != nil {
		return fmt.Errorf("token missing: %w", err)
	}
	dbPath, err := config.DBPath()
	if err != nil {
		return err
	}
	st, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer st.Close()

	start, _ := calcPeriod(*daysFlag)
	if *incrementalFlag && !*fullFlag {
		if cur, ok, err := st.GetCursor(repo); err == nil && ok && !cur.LastCreatedAt.IsZero() {
			start = cur.LastCreatedAt
		}
	}

	client := gh.NewClient(token)
	runList, runCalls, err := client.ListWorkflowRuns(context.Background(), owner, repoName, start)
	if err != nil {
		return err
	}
	for i := range runList {
		runList[i].Repo = repo
	}
	if len(runList) == 0 {
		fmt.Printf("No workflow runs found for %s since %s\n", repo, start.Format(time.RFC3339))
		return nil
	}
	sort.Slice(runList, func(i, j int) bool { return runList[i].CreatedAt.After(runList[j].CreatedAt) })

	newRuns, updatedRuns, err := st.UpsertRuns(runList)
	if err != nil {
		return err
	}

	type result struct {
		jobs     []model.Job
		calls    int
		err      error
		runID    int64
		attempt  int
		workflow string
	}
	sem := make(chan struct{}, *workersFlag)
	runCh := make(chan model.WorkflowRun)
	resCh := make(chan result)
	var wg sync.WaitGroup

	for i := 0; i < *workersFlag; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range runCh {
				sem <- struct{}{}
				jobs, calls, err := client.ListJobsForRun(context.Background(), owner, repoName, r.ID, r.RunAttempt)
				<-sem
				resCh <- result{jobs: jobs, calls: calls, err: err, runID: r.ID, attempt: r.RunAttempt, workflow: r.WorkflowName}
			}
		}()
	}
	go func() {
		for _, r := range runList {
			runCh <- r
		}
		close(runCh)
		wg.Wait()
		close(resCh)
	}()

	allJobs := make([]model.Job, 0, len(runList)*4)
	jobCalls := 0
	failures := 0
	for r := range resCh {
		jobCalls += r.calls
		if r.err != nil {
			failures++
			fmt.Fprintf(os.Stderr, "WARN: jobs fetch failed for run %d attempt %d (%s): %v\n", r.runID, r.attempt, r.workflow, r.err)
			continue
		}
		for i := range r.jobs {
			r.jobs[i].Repo = repo
		}
		allJobs = append(allJobs, r.jobs...)
	}

	newJobs, updatedJobs, err := st.UpsertJobs(allJobs)
	if err != nil {
		return err
	}

	latest := runList[0]
	cursor := store.SyncCursor{
		Repo:          repo,
		LastRunID:     latest.ID,
		LastCreatedAt: latest.CreatedAt,
		LastSyncAt:    time.Now().UTC(),
		TotalRuns:     len(runList),
		TotalJobs:     len(allJobs),
	}
	if err := st.UpsertCursor(cursor); err != nil {
		return err
	}

	fmt.Printf("CICost Scan: %s\n", repo)
	fmt.Printf("  Period     : %s ~ %s\n", start.Format("2006-01-02"), time.Now().UTC().Format("2006-01-02"))
	fmt.Printf("  Runs found : %d (%d new, %d updated)\n", len(runList), newRuns, updatedRuns)
	fmt.Printf("  Jobs found : %d (%d new, %d updated)\n", len(allJobs), newJobs, updatedJobs)
	fmt.Printf("  API calls  : %d (runs=%d jobs=%d)\n", runCalls+jobCalls, runCalls, jobCalls)
	if failures > 0 {
		fmt.Printf("  Partial    : %d runs failed to fetch jobs\n", failures)
	}
	fmt.Printf("  DB path    : %s\n", dbPath)
	return nil
}
