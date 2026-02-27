package github

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientUsesDefaultAndEnvBaseURL(t *testing.T) {
	t.Setenv("CICOST_GITHUB_API_BASE_URL", "")
	c := NewClient("")
	if c.BaseURL != "https://api.github.com" {
		t.Fatalf("expected default github api base, got %s", c.BaseURL)
	}

	t.Setenv("CICOST_GITHUB_API_BASE_URL", "https://example.test")
	c = NewClient("")
	if c.BaseURL != "https://example.test" {
		t.Fatalf("expected env base url, got %s", c.BaseURL)
	}
}

func TestNewRequestSetsHeaders(t *testing.T) {
	c := NewClient("token-123")
	req, err := c.newRequest(context.Background(), http.MethodGet, "https://api.github.com/user")
	if err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("Accept"); got != "application/vnd.github+json" {
		t.Fatalf("unexpected accept header: %s", got)
	}
	if got := req.Header.Get("X-GitHub-Api-Version"); got != "2022-11-28" {
		t.Fatalf("unexpected api version header: %s", got)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer token-123" {
		t.Fatalf("unexpected auth header: %s", got)
	}
}

func TestDoJSONSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := &Client{HTTPClient: srv.Client()}
	req, err := c.newRequest(context.Background(), http.MethodGet, srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		OK bool `json:"ok"`
	}
	resp, err := c.doJSON(req, &payload)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if !payload.OK {
		t.Fatalf("expected payload ok=true")
	}
}

func TestDoJSONReturnsAPIErrorWithMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"resource missing"}`))
	}))
	defer srv.Close()

	c := &Client{HTTPClient: srv.Client()}
	req, err := c.newRequest(context.Background(), http.MethodGet, srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any
	_, err = c.doJSON(req, &payload)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "resource missing" {
		t.Fatalf("unexpected message: %s", apiErr.Message)
	}
}
