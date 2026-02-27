package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/peter941221/CICost/internal/model"
)

func TestBillingSnapshotRoundtrip(t *testing.T) {
	db := filepath.Join(t.TempDir(), "cicost.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	snap := model.BillingSnapshot{
		Repo:          "owner/repo",
		Period:        "2026-02",
		ActualCostUSD: 123.45,
		Source:        "csv",
		FetchedAt:     time.Date(2026, 2, 26, 10, 0, 0, 0, time.UTC),
	}
	if err := st.UpsertBillingSnapshot(snap); err != nil {
		t.Fatal(err)
	}

	got, ok, err := st.GetBillingSnapshot("owner/repo", "2026-02")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected snapshot")
	}
	if got.ActualCostUSD != snap.ActualCostUSD || got.Source != "csv" {
		t.Fatalf("unexpected snapshot: %+v", got)
	}
}

func TestReconcileAndPolicyAndSuggestionPersistence(t *testing.T) {
	db := filepath.Join(t.TempDir(), "cicost.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	rec := model.ReconcileResult{
		Repo:              "owner/repo",
		Period:            "2026-02",
		EstimatedCostUSD:  90,
		ActualCostUSD:     100,
		DeltaRatio:        -0.10,
		CalibrationFactor: 1.11,
		Confidence:        "medium",
	}
	if err := st.InsertReconcileResult(rec); err != nil {
		t.Fatal(err)
	}
	gotRec, ok, err := st.GetLatestReconcile("owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected reconcile result")
	}
	if gotRec.Confidence != "medium" || gotRec.CalibrationFactor != 1.11 {
		t.Fatalf("unexpected reconcile result: %+v", gotRec)
	}

	pRun := model.PolicyRun{
		Repo:          "owner/repo",
		PeriodStart:   time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:     time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		RuleID:        "budget_cap",
		Severity:      "error",
		Matched:       true,
		EvidenceKey:   "monthly_cost_usd",
		EvidenceValue: 250,
		Expression:    "monthly_cost_usd > 200",
	}
	if err := st.InsertPolicyRun(pRun); err != nil {
		t.Fatal(err)
	}
	var matched int
	if err := st.db.QueryRow(`SELECT matched FROM policy_runs WHERE repo=? AND rule_id=?`, "owner/repo", "budget_cap").Scan(&matched); err != nil {
		t.Fatal(err)
	}
	if matched != 1 {
		t.Fatalf("expected matched=1, got %d", matched)
	}

	sugg := model.SuggestionRecord{
		Repo:               "owner/repo",
		PeriodStart:        time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:          time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		SuggestionType:     "cache",
		Title:              "Enable dependency cache",
		EstimatedSavingUSD: 19.5,
		EvidenceJSON:       "{invalid",
	}
	if err := st.InsertSuggestionHistory(sugg); err != nil {
		t.Fatal(err)
	}
	var evidence string
	if err := st.db.QueryRow(`SELECT evidence_json FROM suggestion_history WHERE repo=? AND suggestion_type=?`, "owner/repo", "cache").Scan(&evidence); err != nil {
		t.Fatal(err)
	}
	if evidence != "{}" {
		t.Fatalf("expected sanitized evidence json, got %s", evidence)
	}
}
