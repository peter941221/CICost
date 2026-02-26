package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLintAndEvaluate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".cicost.policy.yml")
	content := `rules:
  - id: budget_monthly
    when: monthly_cost_usd > 200
    severity: error
  - id: waste_ratio
    when: waste_percentage > 25
    severity: warn
actions:
  on_error: fail_ci
  on_warn: comment_pr
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := Lint(cfg); err != nil {
		t.Fatal(err)
	}

	findings, err := Evaluate(cfg, map[string]float64{
		"monthly_cost_usd": 300,
		"waste_percentage": 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "budget_monthly" {
		t.Fatalf("expected budget_monthly, got %s", findings[0].RuleID)
	}
	if findings[0].Severity != SeverityError {
		t.Fatalf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestLintInvalidExpression(t *testing.T) {
	cfg := Config{
		Rules: []Rule{
			{ID: "bad", When: "monthly_cost_usd >> 200", Severity: SeverityError},
		},
	}
	if err := Lint(cfg); err == nil {
		t.Fatal("expected lint error")
	}
}
