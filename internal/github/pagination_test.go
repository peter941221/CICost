package github

import (
	"net/http"
	"testing"
)

func TestNextPageURL(t *testing.T) {
	h := http.Header{}
	h.Set("Link", `<https://api.github.com/resource?page=2>; rel="next", <https://api.github.com/resource?page=3>; rel="last"`)
	got := NextPageURL(h)
	if got != "https://api.github.com/resource?page=2" {
		t.Fatalf("unexpected next url: %s", got)
	}
}

func TestNextPageURLEmpty(t *testing.T) {
	h := http.Header{}
	if got := NextPageURL(h); got != "" {
		t.Fatalf("expected empty, got %s", got)
	}
}
