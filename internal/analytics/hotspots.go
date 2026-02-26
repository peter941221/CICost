package analytics

import (
	"sort"
	"strings"

	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/pricing"
)

type HotspotOptions struct {
	GroupBy string
	TopN    int
	SortBy  string
}

func CalculateHotspots(runs []model.WorkflowRun, jobs []model.Job, cfg pricing.Config, opts HotspotOptions) []model.HotspotEntry {
	if opts.TopN <= 0 {
		opts.TopN = 10
	}
	if opts.GroupBy == "" {
		opts.GroupBy = "workflow"
	}
	if opts.SortBy == "" {
		opts.SortBy = "cost"
	}

	runByIDAttempt := map[string]model.WorkflowRun{}
	for _, r := range runs {
		runByIDAttempt[runAttemptKey(r.ID, r.RunAttempt)] = r
	}

	type agg struct {
		name      string
		groupType string
		minutes   float64
		cost      float64
		runIDs    map[int64]struct{}
		failRuns  int
		totalRuns int
		durSum    float64
		jobCount  int
	}
	aggs := map[string]*agg{}
	totalCost := 0.0
	for _, j := range jobs {
		if j.IsSelfHosted || j.Status != "completed" {
			continue
		}
		key := runAttemptKey(j.RunID, j.RunAttempt)
		run := runByIDAttempt[key]
		groupName := hotspotGroupName(opts.GroupBy, run, j)
		if groupName == "" {
			groupName = "unknown"
		}
		a := aggs[groupName]
		if a == nil {
			a = &agg{name: groupName, groupType: opts.GroupBy, runIDs: map[int64]struct{}{}}
			aggs[groupName] = a
		}
		billable := pricing.BillableMinutes(j.DurationSec, j.RunnerOS, cfg)
		cost := billable * cfg.PerMinuteUSD
		a.minutes += billable
		a.cost += cost
		a.runIDs[j.RunID] = struct{}{}
		a.durSum += float64(j.DurationSec)
		a.jobCount++
		totalCost += cost
	}

	for _, r := range runs {
		name := hotspotGroupName(opts.GroupBy, r, model.Job{})
		if a, ok := aggs[name]; ok {
			if r.Conclusion == "failure" {
				a.failRuns++
			}
			a.totalRuns++
		}
	}

	out := make([]model.HotspotEntry, 0, len(aggs))
	for _, a := range aggs {
		entry := model.HotspotEntry{
			Name:      a.name,
			GroupType: a.groupType,
			Minutes:   round2(a.minutes),
			CostUSD:   round2(a.cost),
			RunCount:  len(a.runIDs),
		}
		if totalCost > 0 {
			entry.CostPct = round2((a.cost / totalCost) * 100)
		}
		if a.totalRuns > 0 {
			entry.FailRate = round2((float64(a.failRuns) / float64(a.totalRuns)) * 100)
		}
		if a.jobCount > 0 {
			entry.AvgDuration = round2(a.durSum / float64(a.jobCount))
		}
		out = append(out, entry)
	}

	sort.Slice(out, func(i, j int) bool {
		switch opts.SortBy {
		case "minutes":
			return out[i].Minutes > out[j].Minutes
		case "fail_rate":
			return out[i].FailRate > out[j].FailRate
		default:
			return out[i].CostUSD > out[j].CostUSD
		}
	})

	if len(out) > opts.TopN {
		out = out[:opts.TopN]
	}
	for i := range out {
		out[i].Rank = i + 1
	}
	return out
}

func hotspotGroupName(groupBy string, run model.WorkflowRun, job model.Job) string {
	switch strings.ToLower(groupBy) {
	case "job":
		if run.WorkflowName != "" && job.Name != "" {
			return run.WorkflowName + " / " + job.Name
		}
		return job.Name
	case "runner":
		if job.IsSelfHosted {
			return "self-hosted"
		}
		if job.RunnerOS != "" {
			return job.RunnerOS
		}
		return "unknown"
	case "branch":
		if run.HeadBranch == "" {
			return "(no-branch)"
		}
		return run.HeadBranch
	default:
		if run.WorkflowName == "" {
			return "(unknown-workflow)"
		}
		return run.WorkflowName
	}
}
