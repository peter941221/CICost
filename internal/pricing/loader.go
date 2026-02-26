package pricing

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrPricingLoaderNotImplemented = errors.New("pricing loader not implemented")

type fileConfig struct {
	Version       string             `yaml:"version"`
	PerMinuteUSD  float64            `yaml:"per_minute_usd"`
	Multipliers   map[string]float64 `yaml:"multipliers"`
	LargerRunners map[string]float64 `yaml:"larger_runners"`
	FreeTiers     map[string]float64 `yaml:"free_tiers"`
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
	if v := raw.Multipliers["Windows"]; v > 0 {
		cfg.WindowsMultiplier = v
	}
	if v := raw.Multipliers["macOS"]; v > 0 {
		cfg.MacOSMultiplier = v
	}
	if cfg.PerMinuteUSD <= 0 {
		cfg.PerMinuteUSD = 0.008
	}
	return cfg, nil
}
