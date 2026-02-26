package pricing

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBillableMinutes(t *testing.T) {
	cfg := Config{
		WindowsMultiplier: 2,
		MacOSMultiplier:   10,
	}
	tests := []struct {
		name     string
		sec      int
		os       string
		expected float64
	}{
		{"linux_59s", 59, "Linux", 1},
		{"linux_60s", 60, "Linux", 1},
		{"linux_61s", 61, "Linux", 2},
		{"win_200s", 200, "Windows", 8},
		{"mac_125s", 125, "macOS", 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BillableMinutes(tt.sec, tt.os, cfg)
			if got != tt.expected {
				t.Fatalf("expected %.2f, got %.2f", tt.expected, got)
			}
		})
	}
}

func TestChargedMinutes(t *testing.T) {
	got := ChargedMinutes(2500, 2000, 0)
	if got != 500 {
		t.Fatalf("expected 500, got %.2f", got)
	}
	got = ChargedMinutes(1000, 2000, 0)
	if got != 0 {
		t.Fatalf("expected 0, got %.2f", got)
	}
}

func TestLoadPricingSnapshotsAndSelect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pricing.yml")
	content := `version: "2026.02"
pricing_snapshots:
  - version: "2026.01"
    effective_from: "2026-01-01"
    skus:
      linux: 0.008
      windows: 0.016
      macos: 0.08
  - version: "2026.02"
    effective_from: "2026-02-15"
    skus:
      linux: 0.0075
      windows: 0.015
      macos: 0.075
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Snapshots) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(cfg.Snapshots))
	}

	s1, err := SelectSnapshot(cfg, time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if s1.Version != "2026.01" {
		t.Fatalf("expected version 2026.01, got %s", s1.Version)
	}

	s2, err := SelectSnapshot(cfg, time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if s2.Version != "2026.02" {
		t.Fatalf("expected version 2026.02, got %s", s2.Version)
	}
}

func TestResolveRateBySKU(t *testing.T) {
	cfg := Config{
		Snapshots: []Snapshot{
			{
				Version:       "2026.02",
				EffectiveFrom: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
				SKUs: map[string]float64{
					"linux":   0.008,
					"windows": 0.016,
					"macos":   0.08,
				},
			},
		},
	}
	rate, sku, source, snap, err := ResolveRate(cfg, time.Date(2026, 2, 26, 0, 0, 0, 0, time.UTC), "macOS", "")
	if err != nil {
		t.Fatal(err)
	}
	if source != PricingSourceSKU {
		t.Fatalf("expected sku source, got %s", source)
	}
	if sku != "macos" {
		t.Fatalf("expected sku macos, got %s", sku)
	}
	if snap.Version != "2026.02" {
		t.Fatalf("expected snapshot 2026.02, got %s", snap.Version)
	}
	if rate != 0.08 {
		t.Fatalf("expected rate 0.08, got %f", rate)
	}
}

func TestResolveRateMissingSnapshot(t *testing.T) {
	cfg := Config{
		Snapshots: []Snapshot{
			{
				Version:       "2026.03",
				EffectiveFrom: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
				SKUs:          map[string]float64{"linux": 0.008},
			},
		},
	}
	_, _, _, _, err := ResolveRate(cfg, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), "Linux", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
