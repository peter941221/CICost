package analytics

import (
	"testing"
	"time"
)

func TestEvaluateBudget(t *testing.T) {
	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)

	ok := EvaluateBudget(now, 30, 100, "monthly")
	if ok.Status != BudgetOK {
		t.Fatalf("expected OK, got %s", ok.Status)
	}

	ex := EvaluateBudget(now, 120, 100, "monthly")
	if ex.Status != BudgetExceeded {
		t.Fatalf("expected EXCEEDED, got %s", ex.Status)
	}

	warn := EvaluateBudget(now, 95, 100, "monthly")
	if warn.Status != BudgetWarning {
		t.Fatalf("expected WARNING, got %s", warn.Status)
	}
}
