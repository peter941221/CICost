package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/store"
)

type orgRepoSummary struct {
	Repo          string  `json:"repo"`
	TotalRuns     int     `json:"total_runs"`
	TotalCostUSD  float64 `json:"total_cost_usd"`
	TotalWasteUSD float64 `json:"total_waste_usd"`
}

type orgHotspot struct {
	Repo    string  `json:"repo"`
	Name    string  `json:"name"`
	CostUSD float64 `json:"cost_usd"`
}

type orgReportPayload struct {
	GeneratedAt  time.Time         `json:"generated_at"`
	Days         int               `json:"days"`
	TotalRepos   int               `json:"total_repos"`
	SuccessRepos int               `json:"success_repos"`
	FailedRepos  int               `json:"failed_repos"`
	TotalCostUSD float64           `json:"total_cost_usd"`
	RepoRankings []orgRepoSummary  `json:"repo_rankings"`
	TopHotspots  []orgHotspot      `json:"top_hotspots"`
	Failures     map[string]string `json:"failures,omitempty"`
}

func runOrgReport(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("org-report", flag.ContinueOnError)
	reposFlag := fs.String("repos", "", "Repository list file (one owner/repo per line)")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "Time window in days")
	formatFlag := fs.String("format", "md", "Output format: md|json")
	outputFlag := fs.String("output", "", "Output file path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*reposFlag) == "" {
		return fmt.Errorf("--repos is required")
	}

	repos, err := readRepoList(*reposFlag)
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		return fmt.Errorf("no repositories found in %s", *reposFlag)
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

	pcfg, err := loadPricingConfig(rt)
	if err != nil {
		return err
	}
	start, end := calcPeriod(*daysFlag)

	type repoResult struct {
		repo     string
		summary  orgRepoSummary
		hotspots []orgHotspot
		err      error
	}
	in := make(chan string)
	out := make(chan repoResult)
	workers := 4
	if len(repos) < workers {
		workers = len(repos)
	}
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range in {
				runs, err := st.ListRuns(repo, start, end)
				if err != nil {
					out <- repoResult{repo: repo, err: err}
					continue
				}
				jobs, err := st.ListJobs(repo, start, end)
				if err != nil {
					out <- repoResult{repo: repo, err: err}
					continue
				}
				if len(runs) == 0 || len(jobs) == 0 {
					out <- repoResult{repo: repo, err: fmt.Errorf("no local data")}
					continue
				}
				cost, _, _, err := analytics.CalculateCostDetailed(jobs, pcfg, 1.0)
				if err != nil {
					out <- repoResult{repo: repo, err: err}
					continue
				}
				waste := analytics.CalculateWaste(runs, jobs, pcfg, cost.TotalCostUSD)
				entries := analytics.CalculateHotspots(runs, jobs, pcfg, analytics.HotspotOptions{
					GroupBy: "workflow",
					TopN:    3,
					SortBy:  "cost",
				})
				hs := make([]orgHotspot, 0, len(entries))
				for _, e := range entries {
					hs = append(hs, orgHotspot{Repo: repo, Name: e.Name, CostUSD: e.CostUSD})
				}
				out <- repoResult{
					repo: repo,
					summary: orgRepoSummary{
						Repo:          repo,
						TotalRuns:     len(runs),
						TotalCostUSD:  cost.TotalCostUSD,
						TotalWasteUSD: waste.TotalWasteUSD,
					},
					hotspots: hs,
				}
			}
		}()
	}

	go func() {
		for _, r := range repos {
			in <- r
		}
		close(in)
		wg.Wait()
		close(out)
	}()

	report := orgReportPayload{
		GeneratedAt: time.Now().UTC(),
		Days:        *daysFlag,
		TotalRepos:  len(repos),
		Failures:    map[string]string{},
	}

	for r := range out {
		if r.err != nil {
			report.Failures[r.repo] = r.err.Error()
			continue
		}
		report.RepoRankings = append(report.RepoRankings, r.summary)
		report.TopHotspots = append(report.TopHotspots, r.hotspots...)
		report.TotalCostUSD += r.summary.TotalCostUSD
	}
	report.SuccessRepos = len(report.RepoRankings)
	report.FailedRepos = len(report.Failures)
	sort.Slice(report.RepoRankings, func(i, j int) bool {
		return report.RepoRankings[i].TotalCostUSD > report.RepoRankings[j].TotalCostUSD
	})
	sort.Slice(report.TopHotspots, func(i, j int) bool {
		return report.TopHotspots[i].CostUSD > report.TopHotspots[j].CostUSD
	})
	if len(report.TopHotspots) > 10 {
		report.TopHotspots = report.TopHotspots[:10]
	}
	report.TotalCostUSD = round2(report.TotalCostUSD)

	var rendered string
	switch strings.ToLower(strings.TrimSpace(*formatFlag)) {
	case "json":
		b, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		rendered = string(b)
	default:
		rendered = renderOrgReportMarkdown(report)
	}
	return writeOutput(*outputFlag, rendered)
}

func readRepoList(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		s := strings.TrimSpace(ln)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

func renderOrgReportMarkdown(report orgReportPayload) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# CICost Org Report\n\n")
	fmt.Fprintf(&b, "- GeneratedAt: `%s`\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- Days: `%d`\n", report.Days)
	fmt.Fprintf(&b, "- Total repos: `%d` (success=%d, failed=%d)\n", report.TotalRepos, report.SuccessRepos, report.FailedRepos)
	fmt.Fprintf(&b, "- Total cost: `$%.2f`\n\n", report.TotalCostUSD)

	fmt.Fprintf(&b, "## Repo Ranking\n\n")
	fmt.Fprintf(&b, "| Repo | Runs | Cost(USD) | Waste(USD) |\n|---|---:|---:|---:|\n")
	for _, r := range report.RepoRankings {
		fmt.Fprintf(&b, "| %s | %d | %.2f | %.2f |\n", r.Repo, r.TotalRuns, r.TotalCostUSD, r.TotalWasteUSD)
	}
	fmt.Fprintf(&b, "\n## Top Hotspots\n\n")
	fmt.Fprintf(&b, "| Repo | Workflow | Cost(USD) |\n|---|---|---:|\n")
	for _, h := range report.TopHotspots {
		fmt.Fprintf(&b, "| %s | %s | %.2f |\n", h.Repo, h.Name, h.CostUSD)
	}
	if report.FailedRepos > 0 {
		fmt.Fprintf(&b, "\n## Partial Failures\n\n")
		for repo, reason := range report.Failures {
			fmt.Fprintf(&b, "- %s: %s\n", repo, reason)
		}
	}
	return b.String()
}
