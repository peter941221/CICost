package suggest

import (
	"fmt"
	"strings"

	"github.com/peter941221/CICost/internal/model"
)

type Inputs struct {
	Repo     string
	Runs     []model.WorkflowRun
	Jobs     []model.Job
	Cost     model.CostResult
	Waste    model.WasteMetrics
	Hotspots []model.HotspotEntry
}

type Suggestion struct {
	Type               string         `json:"type" yaml:"type"`
	Title              string         `json:"title" yaml:"title"`
	Problem            string         `json:"problem" yaml:"problem"`
	CurrentData        string         `json:"current_data" yaml:"current_data"`
	EstimatedSavingUSD float64        `json:"estimated_saving_usd" yaml:"estimated_saving_usd"`
	Patch              string         `json:"patch" yaml:"patch"`
	Evidence           map[string]any `json:"evidence" yaml:"evidence"`
}

func Generate(in Inputs) []Suggestion {
	out := make([]Suggestion, 0, 4)

	if in.Waste.CancelWasteUSD > 0 {
		out = append(out, Suggestion{
			Type:               "concurrency",
			Title:              "Enable cancel-in-progress concurrency",
			Problem:            "Cancelled runs are consuming avoidable CI spend.",
			CurrentData:        fmt.Sprintf("cancel_waste_usd=%.2f, cancelled_runs=%d", in.Waste.CancelWasteUSD, in.Waste.CancelledRuns),
			EstimatedSavingUSD: round2(in.Waste.CancelWasteUSD * 0.8),
			Patch: `concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true`,
			Evidence: evidence(firstWorkflow(in.Hotspots), "", pct(in.Waste.FailRate), in.Waste.CancelWasteUSD),
		})
	}

	if len(in.Hotspots) > 0 && in.Hotspots[0].CostUSD > 0 {
		top := in.Hotspots[0]
		out = append(out, Suggestion{
			Type:               "cache",
			Title:              "Add dependency cache to hottest workflow",
			Problem:            "High-cost workflow repeats dependency resolution.",
			CurrentData:        fmt.Sprintf("workflow=%s, workflow_cost_usd=%.2f, fail_rate=%.2f%%", top.Name, top.CostUSD, top.FailRate),
			EstimatedSavingUSD: round2(top.CostUSD * 0.15),
			Patch: `- uses: actions/cache@v4
  with:
    path: |
      ~/.npm
      ~/.cache/pip
    key: ${{ runner.os }}-${{ hashFiles('**/package-lock.json', '**/requirements.txt') }}`,
			Evidence: evidence(top.Name, "", top.FailRate, top.CostUSD),
		})
	}

	pushCount := 0
	for _, r := range in.Runs {
		if strings.EqualFold(r.Event, "push") {
			pushCount++
		}
	}
	if pushCount >= 20 && len(in.Runs) > 0 && in.Cost.TotalCostUSD > 0 {
		avg := in.Cost.TotalCostUSD / float64(len(in.Runs))
		out = append(out, Suggestion{
			Type:               "paths",
			Title:              "Add paths filters to reduce unnecessary runs",
			Problem:            "Frequent push-triggered workflows are likely over-triggering.",
			CurrentData:        fmt.Sprintf("push_runs=%d, avg_cost_per_run=%.2f", pushCount, avg),
			EstimatedSavingUSD: round2(in.Cost.TotalCostUSD * 0.10),
			Patch: `on:
  push:
    paths:
      - "src/**"
      - ".github/workflows/**"`,
			Evidence: evidence(firstWorkflow(in.Hotspots), "", pct(in.Waste.FailRate), in.Cost.TotalCostUSD),
		})
	}

	if in.Cost.TotalCostUSD > 0 {
		if mac, ok := in.Cost.ByOS["macOS"]; ok && mac.CostUSD > 0 {
			share := mac.CostUSD / in.Cost.TotalCostUSD
			if share >= 0.2 {
				out = append(out, Suggestion{
					Type:               "runner_migration",
					Title:              "Evaluate migration from macOS to Linux runner",
					Problem:            "macOS spend share is high.",
					CurrentData:        fmt.Sprintf("macos_cost_usd=%.2f, total_cost_usd=%.2f, share=%.1f%%", mac.CostUSD, in.Cost.TotalCostUSD, share*100),
					EstimatedSavingUSD: round2(mac.CostUSD * 0.6),
					Patch: `jobs:
  build:
    runs-on: ubuntu-latest`,
					Evidence: evidence(firstWorkflow(in.Hotspots), "", pct(in.Waste.FailRate), mac.CostUSD),
				})
			}
		}
	}

	filtered := make([]Suggestion, 0, len(out))
	for _, s := range out {
		if s.EstimatedSavingUSD <= 0 || len(s.Evidence) == 0 {
			continue
		}
		filtered = append(filtered, s)
	}
	return filtered
}

func evidence(workflow, job string, failRate, cost float64) map[string]any {
	return map[string]any{
		"workflow":  workflow,
		"job":       job,
		"fail_rate": round2(failRate),
		"cost":      round2(cost),
	}
}

func firstWorkflow(entries []model.HotspotEntry) string {
	if len(entries) == 0 {
		return ""
	}
	return entries[0].Name
}

func pct(v float64) float64 {
	return v * 100
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
