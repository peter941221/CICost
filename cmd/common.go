package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/pricing"
)

type runtimeContext struct {
	cfg config.Config
}

func newRuntimeContext() (runtimeContext, error) {
	cfg, err := config.LoadMerged(".cicost.yml")
	if err != nil {
		return runtimeContext{}, err
	}
	return runtimeContext{cfg: cfg}, nil
}

func pickRepo(input string, cfg config.Config) (string, error) {
	if strings.TrimSpace(input) != "" {
		return strings.TrimSpace(input), nil
	}
	if len(cfg.Repos) > 0 && strings.TrimSpace(cfg.Repos[0]) != "" {
		return strings.TrimSpace(cfg.Repos[0]), nil
	}
	return "", errors.New("repo is required (use --repo owner/repo)")
}

func splitRepo(full string) (owner, repo string, err error) {
	parts := strings.Split(strings.TrimSpace(full), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format %q, expected owner/repo", full)
	}
	return parts[0], parts[1], nil
}

func calcPeriod(days int) (time.Time, time.Time) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -days+1)
	end := now
	return start, end
}

func writeOutput(path string, content string) error {
	if path == "" {
		fmt.Print(content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Println()
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func loadPricingConfig(rt runtimeContext) (pricing.Config, error) {
	candidates := []string{
		filepath.Join("configs", "pricing_default.yml"),
		filepath.Join("..", "configs", "pricing_default.yml"),
	}
	var cfg pricing.Config
	var loaded bool
	for _, p := range candidates {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		c, err := pricing.LoadFromFile(p)
		if err != nil {
			return pricing.Config{}, err
		}
		cfg = c
		loaded = true
		break
	}
	if !loaded {
		cfg = pricing.Config{}
	}
	if cfg.PerMinuteUSD == 0 {
		cfg.PerMinuteUSD = rt.cfg.Pricing.LinuxPerMin
	}
	if rt.cfg.Pricing.WindowsMultiplier > 0 {
		cfg.WindowsMultiplier = rt.cfg.Pricing.WindowsMultiplier
	}
	if rt.cfg.Pricing.MacOSMultiplier > 0 {
		cfg.MacOSMultiplier = rt.cfg.Pricing.MacOSMultiplier
	}
	cfg.FreeTierPerMonth = rt.cfg.FreeTier.MinutesPerMonth
	if len(cfg.Snapshots) == 0 {
		fmt.Fprintln(os.Stderr, "WARN: pricing_snapshots not found; fallback to legacy OS multiplier pricing.")
	}
	return cfg, nil
}
