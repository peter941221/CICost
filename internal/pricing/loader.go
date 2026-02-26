package pricing

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrPricingLoaderNotImplemented = errors.New("pricing loader not implemented")

type fileConfig struct {
	Version          string             `yaml:"version"`
	EffectiveFrom    string             `yaml:"effective_from"`
	PerMinuteUSD     float64            `yaml:"per_minute_usd"`
	Multipliers      map[string]float64 `yaml:"multipliers"`
	LargerRunners    map[string]float64 `yaml:"larger_runners"`
	FreeTiers        map[string]float64 `yaml:"free_tiers"`
	PricingSnapshots []snapshotFile     `yaml:"pricing_snapshots"`
}

type snapshotFile struct {
	Version       string             `yaml:"version"`
	EffectiveFrom string             `yaml:"effective_from"`
	SKUs          map[string]float64 `yaml:"skus"`
}

func LoadFromFile(path string) (Config, error) {
	var raw fileConfig
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return Config{}, fmt.Errorf("parse pricing file failed: %w", err)
	}
	cfg := Config{
		Version:             raw.Version,
		PerMinuteUSD:        raw.PerMinuteUSD,
		FreeTierByPlan:      raw.FreeTiers,
		LargerRunnersPerMin: raw.LargerRunners,
		WindowsMultiplier:   2,
		MacOSMultiplier:     10,
	}
	if strings.TrimSpace(raw.EffectiveFrom) != "" {
		t, err := time.Parse("2006-01-02", raw.EffectiveFrom)
		if err != nil {
			return Config{}, fmt.Errorf("invalid effective_from %q: %w", raw.EffectiveFrom, err)
		}
		cfg.EffectiveFrom = t.UTC()
	}
	if v := raw.Multipliers["Windows"]; v > 0 {
		cfg.WindowsMultiplier = v
	}
	if v := raw.Multipliers["macOS"]; v > 0 {
		cfg.MacOSMultiplier = v
	}
	if cfg.PerMinuteUSD <= 0 {
		cfg.PerMinuteUSD = 0.008
	}

	if len(raw.PricingSnapshots) > 0 {
		cfg.Snapshots = make([]Snapshot, 0, len(raw.PricingSnapshots))
		for _, s := range raw.PricingSnapshots {
			eff, err := time.Parse("2006-01-02", strings.TrimSpace(s.EffectiveFrom))
			if err != nil {
				return Config{}, fmt.Errorf("invalid pricing_snapshot.effective_from %q: %w", s.EffectiveFrom, err)
			}
			if len(s.SKUs) == 0 {
				return Config{}, fmt.Errorf("pricing_snapshot %q has empty skus", s.Version)
			}
			skus := make(map[string]float64, len(s.SKUs))
			for k, v := range s.SKUs {
				if v <= 0 {
					continue
				}
				skus[normalizeSKU(k)] = v
			}
			if len(skus) == 0 {
				return Config{}, fmt.Errorf("pricing_snapshot %q has no valid sku rates", s.Version)
			}
			cfg.Snapshots = append(cfg.Snapshots, Snapshot{
				Version:       strings.TrimSpace(s.Version),
				EffectiveFrom: eff.UTC(),
				SKUs:          skus,
			})
		}
		sort.Slice(cfg.Snapshots, func(i, j int) bool {
			return cfg.Snapshots[i].EffectiveFrom.Before(cfg.Snapshots[j].EffectiveFrom)
		})
		latest := cfg.Snapshots[len(cfg.Snapshots)-1]
		cfg.Version = latest.Version
		cfg.EffectiveFrom = latest.EffectiveFrom
	}
	return cfg, nil
}
