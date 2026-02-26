package pricing

import "testing"

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
