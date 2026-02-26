package analytics

import (
	"sort"
	"strconv"

	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/pricing"
)

func CalculateWaste(runs []model.WorkflowRun, jobs []model.Job, cfg pricing.Config, totalCost float64) model.WasteMetrics {
	type runAgg struct {
		maxAttempt int
		status     string
		conclusion string
	}
	runState := map[int64]runAgg{}
	for _, r := range runs {
		cur := runState[r.ID]
		if r.RunAttempt >= cur.maxAttempt {
			cur.maxAttempt = r.RunAttempt
			cur.status = r.Status
			cur.conclusion = r.Conclusion
		}
		runState[r.ID] = cur
	}

	latestRuns := make([]model.WorkflowRun, 0, len(runState))
	for _, r := range runs {
		if state, ok := runState[r.ID]; ok && r.RunAttempt == state.maxAttempt {
			latestRuns = append(latestRuns, r)
		}
	}
	sort.Slice(latestRuns, func(i, j int) bool { return latestRuns[i].ID < latestRuns[j].ID })

	costByRunAttempt := map[string]float64{}
	minByRunAttempt := map[string]float64{}
	for _, j := range jobs {
		if j.IsSelfHosted || j.Status != "completed" {
			continue
		}
		key := runAttemptKey(j.RunID, j.RunAttempt)
		billable := pricing.BillableMinutes(j.DurationSec, j.RunnerOS, cfg)
		costByRunAttempt[key] += billable * cfg.PerMinuteUSD
		minByRunAttempt[key] += billable
	}

	var m model.WasteMetrics
	m.TotalRuns = len(runState)
	for _, r := range latestRuns {
		if r.Conclusion == "failure" {
			m.FailedRuns++
		}
		if r.Conclusion == "cancelled" {
			key := runAttemptKey(r.ID, r.RunAttempt)
			m.CancelWasteUSD += costByRunAttempt[key]
			m.CancelWasteMin += minByRunAttempt[key]
			m.CancelledRuns++
		}
	}
	if m.TotalRuns > 0 {
		m.FailRate = float64(m.FailedRuns) / float64(m.TotalRuns)
	}

	attemptsByRun := map[int64]int{}
	for _, r := range runs {
		if r.RunAttempt > attemptsByRun[r.ID] {
			attemptsByRun[r.ID] = r.RunAttempt
		}
	}
	for runID, maxAttempt := range attemptsByRun {
		if maxAttempt <= 1 {
			continue
		}
		m.RerunCount++
		for a := 1; a < maxAttempt; a++ {
			key := runAttemptKey(runID, a)
			m.RerunWasteUSD += costByRunAttempt[key]
			m.RerunWasteMin += minByRunAttempt[key]
		}
	}

	m.TotalWasteUSD = round2(m.RerunWasteUSD + m.CancelWasteUSD)
	if totalCost > 0 {
		m.WastePercentage = round2((m.TotalWasteUSD / totalCost) * 100)
	}
	m.RerunWasteUSD = round2(m.RerunWasteUSD)
	m.CancelWasteUSD = round2(m.CancelWasteUSD)
	m.RerunWasteMin = round2(m.RerunWasteMin)
	m.CancelWasteMin = round2(m.CancelWasteMin)
	return m
}

func runAttemptKey(runID int64, attempt int) string {
	return strconv.FormatInt(runID, 10) + "#" + strconv.Itoa(attempt)
}
