package reconcile

import "testing"

func TestBuildResultAndConfidence(t *testing.T) {
	res := BuildResult("owner/repo", "2026-02", 110, 100)
	if res.DeltaRatio != 0.1 {
		t.Fatalf("expected delta 0.1, got %f", res.DeltaRatio)
	}
	if res.CalibrationFactor <= 0 || res.CalibrationFactor >= 1 {
		t.Fatalf("expected factor between 0 and 1, got %f", res.CalibrationFactor)
	}
	if res.Confidence != "medium" {
		t.Fatalf("expected medium confidence, got %s", res.Confidence)
	}
}

func TestConfidenceBands(t *testing.T) {
	if got := Confidence(0.03); got != "high" {
		t.Fatalf("expected high, got %s", got)
	}
	if got := Confidence(-0.10); got != "medium" {
		t.Fatalf("expected medium, got %s", got)
	}
	if got := Confidence(0.20); got != "low" {
		t.Fatalf("expected low, got %s", got)
	}
}
