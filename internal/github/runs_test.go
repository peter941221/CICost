package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListWorkflowRunsPaginationAndFallbackAttempt(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/actions/runs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			w.Header().Set("Link", fmt.Sprintf(`<%s/repos/owner/repo/actions/runs?page=2>; rel="next"`, srv.URL))
			_, _ = w.Write([]byte(`{"total_count":2,"workflow_runs":[{"id":1,"workflow_id":11,"name":"ci","head_branch":"main","event":"push","status":"completed","conclusion":"success","run_attempt":0,"run_started_at":"2026-02-26T10:00:00Z","updated_at":"2026-02-26T10:10:00Z","created_at":"2026-02-26T10:00:00Z"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"total_count":2,"workflow_runs":[{"id":2,"workflow_id":22,"name":"deploy","head_branch":"main","event":"workflow_dispatch","status":"completed","conclusion":"failure","run_attempt":2,"run_started_at":"2026-02-26T11:00:00Z","updated_at":"2026-02-26T11:10:00Z","created_at":"2026-02-26T11:00:00Z"}]}`))
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	}
	runs, calls, err := c.ListWorkflowRuns(context.Background(), "owner", "repo", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 api calls, got %d", calls)
	}
	if len(runs) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(runs))
	}
	if runs[0].RunAttempt != 1 {
		t.Fatalf("expected fallback attempt 1, got %d", runs[0].RunAttempt)
	}
	if runs[0].Repo != "repo" {
		t.Fatalf("expected repo name mapping to repo, got %s", runs[0].Repo)
	}
	if runs[1].RunAttempt != 2 {
		t.Fatalf("expected run attempt 2, got %d", runs[1].RunAttempt)
	}
}

func TestParseRunID(t *testing.T) {
	got := parseRunID("https://api.github.com/repos/a/b/actions/jobs?run_id=12345")
	if got != 12345 {
		t.Fatalf("expected 12345, got %d", got)
	}
	if parseRunID("not-a-url") != 0 {
		t.Fatalf("expected 0 for invalid url")
	}
}

func TestParseTime(t *testing.T) {
	got := parseTime("2026-02-26T10:00:00Z")
	if got.IsZero() {
		t.Fatalf("expected valid time")
	}
	if got.Location() != time.UTC {
		t.Fatalf("expected utc time")
	}
	if !parseTime("invalid").IsZero() {
		t.Fatalf("expected zero for invalid time")
	}
}
