package cmd

import (
	"encoding/json"
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

func TestReconcileAndCalibratedReportIntegration(t *testing.T) {
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
			ID:           9,
			Repo:         "owner/repo",
			WorkflowID:   99,
			WorkflowName: "ci",
			HeadBranch:   "main",
			Event:        "push",
			Status:       "completed",
			Conclusion:   "success",
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
			ID:           91,
			RunID:        9,
			RunAttempt:   1,
			Repo:         "owner/repo",
			Name:         "heavy-test",
			Status:       "completed",
			Conclusion:   "success",
			RunnerOS:     "Linux",
			IsSelfHosted: false,
			DurationSec:  300000,
			StartedAt:    now,
			CompletedAt:  now.Add(300000 * time.Second),
		},
	}
	if _, _, err := st.UpsertJobs(jobs); err != nil {
		t.Fatal(err)
	}

	month := now.Format("2006-01")
	if err := runReconcile([]string{"--repo", "owner/repo", "--month", month, "--actual-usd", "5", "--apply-calibration"}); err != nil {
		t.Fatal(err)
	}

	reportFile := filepath.Join(tmp, "report-calibrated.json")
	if err := runReport([]string{"--repo", "owner/repo", "--days", "30", "--format", "json", "--calibrated", "--output", reportFile}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(reportFile)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatal(err)
	}
	report, ok := payload["report"].(map[string]any)
	if !ok {
		t.Fatalf("expected report object")
	}
	cal, ok := report["calibrated"].(bool)
	if !ok || !cal {
		t.Fatalf("expected calibrated=true in report")
	}
	if _, ok := report["calibration_factor"]; !ok {
		t.Fatalf("expected calibration_factor in report")
	}
}

func TestPolicyCheckExitCodes(t *testing.T) {
	tmp := t.TempDir()
	originalHome := os.Getenv("USERPROFILE")
	originalHomeUnix := os.Getenv("HOME")
	originalWD, _ := os.Getwd()
	t.Cleanup(func() {
		_ = os.Setenv("USERPROFILE", originalHome)
		_ = os.Setenv("HOME", originalHomeUnix)
		_ = os.Chdir(originalWD)
	})
	_ = os.Setenv("USERPROFILE", tmp)
	_ = os.Setenv("HOME", tmp)
	_ = os.Chdir(tmp)

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
			ID:           301,
			Repo:         "owner/repo",
			WorkflowID:   1,
			WorkflowName: "build",
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
			ID:          302,
			RunID:       301,
			RunAttempt:  1,
			Repo:        "owner/repo",
			Name:        "test",
			Status:      "completed",
			Conclusion:  "failure",
			RunnerOS:    "Linux",
			DurationSec: 400000,
			StartedAt:   now,
			CompletedAt: now.Add(400000 * time.Second),
		},
	}
	if _, _, err := st.UpsertJobs(jobs); err != nil {
		t.Fatal(err)
	}

	policyPath := filepath.Join(tmp, ".cicost.policy.yml")
	content := `rules:
  - id: over_budget
    when: monthly_cost_usd > 1
    severity: error
`
	if err := os.WriteFile(policyPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err = runPolicy([]string{"check", "--repo", "owner/repo", "--days", "30", "--policy", policyPath})
	if err == nil {
		t.Fatal("expected exit error code 3")
	}
	var ex ExitError
	if !errors.As(err, &ex) || ex.Code != 3 {
		t.Fatalf("expected exit code 3, got %v", err)
	}

	content = `rules:
  - id: warning_only
    when: waste_percentage > 1
    severity: warn
`
	if err := os.WriteFile(policyPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runPolicy([]string{"check", "--repo", "owner/repo", "--days", "30", "--policy", policyPath}); err != nil {
		t.Fatalf("expected warn-only policy check to pass, got %v", err)
	}
}

func TestSuggestCommandYAMLOutput(t *testing.T) {
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
			ID:           401,
			Repo:         "owner/repo",
			WorkflowID:   1,
			WorkflowName: "ci",
			Event:        "push",
			Status:       "completed",
			Conclusion:   "cancelled",
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
			ID:          402,
			RunID:       401,
			RunAttempt:  1,
			Repo:        "owner/repo",
			Name:        "build",
			Status:      "completed",
			Conclusion:  "cancelled",
			RunnerOS:    "macOS",
			DurationSec: 260000,
			StartedAt:   now,
			CompletedAt: now.Add(260000 * time.Second),
		},
	}
	if _, _, err := st.UpsertJobs(jobs); err != nil {
		t.Fatal(err)
	}

	patchDir := filepath.Join(tmp, "patches")
	if err := runSuggest([]string{"--repo", "owner/repo", "--days", "30", "--format", "yaml", "--output", patchDir}); err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(patchDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("expected patch files generated")
	}
}

func TestOrgReportPartialResult(t *testing.T) {
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
	if _, _, err := st.UpsertRuns([]model.WorkflowRun{
		{
			ID:           501,
			Repo:         "owner/repo-a",
			WorkflowID:   1,
			WorkflowName: "ci",
			Event:        "push",
			Status:       "completed",
			Conclusion:   "success",
			RunAttempt:   1,
			CreatedAt:    now,
			UpdatedAt:    now,
			RunStartedAt: now,
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := st.UpsertJobs([]model.Job{
		{
			ID:          502,
			RunID:       501,
			RunAttempt:  1,
			Repo:        "owner/repo-a",
			Name:        "build",
			Status:      "completed",
			RunnerOS:    "Linux",
			DurationSec: 400000,
			StartedAt:   now,
			CompletedAt: now.Add(400000 * time.Second),
		},
	}); err != nil {
		t.Fatal(err)
	}

	repoFile := filepath.Join(tmp, "repos.txt")
	repoText := "owner/repo-a\nowner/repo-missing\n"
	if err := os.WriteFile(repoFile, []byte(repoText), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(tmp, "org.json")
	if err := runOrgReport([]string{"--repos", repoFile, "--days", "30", "--format", "json", "--output", out}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["failed_repos"].(float64) != 1 {
		t.Fatalf("expected failed_repos=1")
	}
	if payload["success_repos"].(float64) != 1 {
		t.Fatalf("expected success_repos=1")
	}
}
