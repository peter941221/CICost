package analytics

import (
	"testing"

	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/pricing"
)

func TestCalculateWaste(t *testing.T) {
	runs := []model.WorkflowRun{
		{ID: 100, RunAttempt: 1, Conclusion: "failure"},
		{ID: 100, RunAttempt: 2, Conclusion: "success"},
		{ID: 200, RunAttempt: 1, Conclusion: "cancelled"},
	}
	jobs := []model.Job{
		{ID: 1, RunID: 100, RunAttempt: 1, Status: "completed", DurationSec: 120, RunnerOS: "Linux"},
		{ID: 2, RunID: 100, RunAttempt: 2, Status: "completed", DurationSec: 120, RunnerOS: "Linux"},
		{ID: 3, RunID: 200, RunAttempt: 1, Status: "completed", DurationSec: 60, RunnerOS: "Linux"},
	}
	cfg := pricing.Config{PerMinuteUSD: 0.008, WindowsMultiplier: 2, MacOSMultiplier: 10}
	got := CalculateWaste(runs, jobs, cfg, 1.0)
	if got.RerunCount != 1 {
		t.Fatalf("expected rerun count 1, got %d", got.RerunCount)
	}
	if got.RerunWasteUSD <= 0 {
		t.Fatalf("expected rerun waste > 0")
	}
	if got.CancelWasteUSD <= 0 {
		t.Fatalf("expected cancel waste > 0")
	}
	if got.TotalRuns != 2 { // latest attempts deduped by run id
		t.Fatalf("expected total runs 2, got %d", got.TotalRuns)
	}
}
