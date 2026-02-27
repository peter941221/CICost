package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListJobsForRunUsesAttemptFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/actions/runs/77/attempts/1/jobs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"total_count":1,"jobs":[{"id":7001,"name":"unit-test","status":"completed","conclusion":"success","started_at":"2026-02-26T10:00:00Z","completed_at":"2026-02-26T10:01:30Z","runner_name":"linux-host","runner_group_name":"default","labels":["self-hosted","linux"]}]}`))
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	}
	jobs, calls, err := c.ListJobsForRun(context.Background(), "owner", "repo", 77, 0)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 api call, got %d", calls)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].RunAttempt != 1 {
		t.Fatalf("expected fallback run attempt 1, got %d", jobs[0].RunAttempt)
	}
	if !jobs[0].IsSelfHosted {
		t.Fatalf("expected self hosted job")
	}
	if jobs[0].RunnerOS != "Linux" {
		t.Fatalf("expected Linux, got %s", jobs[0].RunnerOS)
	}
	if jobs[0].DurationSec != 90 {
		t.Fatalf("expected duration 90 sec, got %d", jobs[0].DurationSec)
	}
}

func TestGuessRunnerOSAndSelfHostedHelpers(t *testing.T) {
	if !hasSelfHosted([]string{"linux", "SELF-HOSTED"}) {
		t.Fatalf("expected case-insensitive self-hosted detection")
	}
	if got := guessRunnerOS([]string{"macos-14"}, "runner-x"); got != "macOS" {
		t.Fatalf("expected macOS, got %s", got)
	}
	if got := guessRunnerOS([]string{}, "win-runner-1"); got != "Windows" {
		t.Fatalf("expected Windows fallback, got %s", got)
	}
	if got := guessRunnerOS([]string{}, "unknown"); got != "Linux" {
		t.Fatalf("expected Linux default, got %s", got)
	}
}
