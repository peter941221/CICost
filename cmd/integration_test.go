package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/store"
)

func TestReportAndBudgetIntegration(t *testing.T) {
	tmp := t.TempDir()
	originalHome := os.Getenv("USERPROFILE")
	originalHomeUnix := os.Getenv("HOME")
	t.Cleanup(func() {
		_ = os.Setenv("USERPROFILE", originalHome)
		_ = os.Setenv("HOME", originalHomeUnix)
	})
	_ = os.Setenv("USERPROFILE", tmp)
	_ = os.Setenv("HOME", tmp)

	dbPath, err := config.DBPath()
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	now := time.Now().UTC()
	runs := []model.WorkflowRun{
		{
			ID:           1,
			Repo:         "owner/repo",
			WorkflowID:   10,
			WorkflowName: "build-and-test",
			HeadBranch:   "main",
			Event:        "push",
			Status:       "completed",
			Conclusion:   "failure",
			RunAttempt:   1,
			CreatedAt:    now,
			UpdatedAt:    now,
			RunStartedAt: now,
		},
	}
	if _, _, err := st.UpsertRuns(runs); err != nil {
		t.Fatal(err)
	}
	jobs := []model.Job{
		{
			ID:           11,
			RunID:        1,
			RunAttempt:   1,
			Repo:         "owner/repo",
			Name:         "unit-tests",
			Status:       "completed",
			Conclusion:   "failure",
			RunnerOS:     "Linux",
			IsSelfHosted: false,
			DurationSec:  200000,
			StartedAt:    now,
			CompletedAt:  now.Add(200000 * time.Second),
		},
	}
	if _, _, err := st.UpsertJobs(jobs); err != nil {
		t.Fatal(err)
	}

	reportFile := filepath.Join(tmp, "report.json")
	if err := runReport([]string{"--repo", "owner/repo", "--days", "30", "--format", "json", "--output", reportFile}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(reportFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"schema_version": "1.0"`) {
		t.Fatalf("expected schema_version in report json")
	}

	err = runBudget([]string{"--repo", "owner/repo", "--monthly", "1"})
	if err == nil {
		t.Fatalf("expected budget exceed error")
	}
	var ex ExitError
	if !errors.As(err, &ex) || ex.Code != 2 {
		t.Fatalf("expected exit code 2, got %v", err)
	}
}
