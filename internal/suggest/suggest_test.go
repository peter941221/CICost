package suggest

import (
	"testing"

	"github.com/peter941221/CICost/internal/model"
)

func TestGenerateProducesExecutableSuggestions(t *testing.T) {
	input := Inputs{
		Cost: model.CostResult{
			TotalCostUSD: 120,
			ByOS: map[string]model.OSCost{
				"macOS": {OS: "macOS", CostUSD: 50},
				"Linux": {OS: "Linux", CostUSD: 70},
			},
		},
		Waste: model.WasteMetrics{
			CancelWasteUSD:  18,
			WastePercentage: 30,
		},
		Hotspots: []model.HotspotEntry{
			{Name: "ci-build", CostUSD: 60, CostPct: 50, FailRate: 25},
		},
		Runs: []model.WorkflowRun{
			{ID: 1, Event: "push"},
			{ID: 2, Event: "push"},
			{ID: 3, Event: "push"},
			{ID: 4, Event: "push"},
			{ID: 5, Event: "push"},
			{ID: 6, Event: "push"},
			{ID: 7, Event: "push"},
			{ID: 8, Event: "push"},
			{ID: 9, Event: "push"},
			{ID: 10, Event: "push"},
			{ID: 11, Event: "push"},
			{ID: 12, Event: "push"},
			{ID: 13, Event: "push"},
			{ID: 14, Event: "push"},
			{ID: 15, Event: "push"},
			{ID: 16, Event: "push"},
			{ID: 17, Event: "push"},
			{ID: 18, Event: "push"},
			{ID: 19, Event: "push"},
			{ID: 20, Event: "push"},
			{ID: 21, Event: "push"},
		},
	}
	suggestions := Generate(input)
	if len(suggestions) < 3 {
		t.Fatalf("expected at least 3 suggestions, got %d", len(suggestions))
	}
	required := map[string]bool{
		"concurrency": false,
		"cache":       false,
		"paths":       false,
	}
	for _, s := range suggestions {
		if s.Problem == "" || s.CurrentData == "" || s.Patch == "" || len(s.Evidence) == 0 {
			t.Fatalf("suggestion %s missing required fields", s.Type)
		}
		if s.EstimatedSavingUSD <= 0 {
			t.Fatalf("suggestion %s expected positive estimated saving", s.Type)
		}
		if _, ok := required[s.Type]; ok {
			required[s.Type] = true
		}
	}
	for typ, ok := range required {
		if !ok {
			t.Fatalf("missing required suggestion type: %s", typ)
		}
	}
}
