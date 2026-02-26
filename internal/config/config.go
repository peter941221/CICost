package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repos []string `yaml:"repos"`
	Auth  struct {
		Token string `yaml:"token"`
	} `yaml:"auth"`
	Scan struct {
		Days        int  `yaml:"days"`
		Workers     int  `yaml:"workers"`
		Incremental bool `yaml:"incremental"`
	} `yaml:"scan"`
	Pricing struct {
		Source            string  `yaml:"source"`
		LinuxPerMin       float64 `yaml:"linux_per_min"`
		WindowsMultiplier float64 `yaml:"windows_multiplier"`
		MacOSMultiplier   float64 `yaml:"macos_multiplier"`
		Currency          string  `yaml:"currency"`
	} `yaml:"pricing"`
	FreeTier struct {
		Plan            string  `yaml:"plan"`
		MinutesPerMonth float64 `yaml:"minutes_per_month"`
	} `yaml:"free_tier"`
	Budget struct {
		Monthly    float64 `yaml:"monthly"`
		Weekly     float64 `yaml:"weekly"`
		Notify     string  `yaml:"notify"`
		WebhookURL string  `yaml:"webhook_url"`
	} `yaml:"budget"`
	Output struct {
		Format string `yaml:"format"`
		Color  string `yaml:"color"`
	} `yaml:"output"`
	Ignore struct {
		Workflows []string `yaml:"workflows"`
	} `yaml:"ignore"`
}

var ErrNoHome = errors.New("unable to resolve user home dir")

func Default() Config {
	var c Config
	c.Scan.Days = 30
	c.Scan.Workers = 4
	c.Scan.Incremental = true
	c.Pricing.Source = "default"
	c.Pricing.LinuxPerMin = 0.008
	c.Pricing.WindowsMultiplier = 2
	c.Pricing.MacOSMultiplier = 10
	c.Pricing.Currency = "USD"
	c.FreeTier.Plan = "free"
	c.FreeTier.MinutesPerMonth = 2000
	c.Budget.Monthly = 100
	c.Budget.Weekly = 0
	c.Budget.Notify = "stdout"
	c.Output.Format = "table"
	c.Output.Color = "auto"
	return c
}

func UserConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", ErrNoHome
	}
	return filepath.Join(home, ".cicost", "config.yml"), nil
}

func DBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", ErrNoHome
	}
	return filepath.Join(home, ".cicost", "data", "cicost.db"), nil
}

func LoadMerged(repoConfigPath string) (Config, error) {
	cfg := Default()

	userCfgPath, err := UserConfigPath()
	if err == nil {
		if loaded, ok, err2 := loadFromFile(userCfgPath); err2 != nil {
			return Config{}, err2
		} else if ok {
			merge(&cfg, loaded)
		}
	}

	if repoConfigPath != "" {
		if loaded, ok, err2 := loadFromFile(repoConfigPath); err2 != nil {
			return Config{}, err2
		} else if ok {
			merge(&cfg, loaded)
		}
	}

	mergeEnv(&cfg)
	return cfg, nil
}

func SaveUserConfig(cfg Config) (string, error) {
	p, err := UserConfigPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(p, out, 0o600); err != nil {
		return "", err
	}
	return p, nil
}

func Expand(value string) string {
	return os.ExpandEnv(value)
}

func loadFromFile(path string) (Config, bool, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, false, nil
		}
		return cfg, false, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, false, fmt.Errorf("invalid yaml in %s: %w", path, err)
	}
	return cfg, true, nil
}

func merge(dst *Config, src Config) {
	if len(src.Repos) > 0 {
		dst.Repos = src.Repos
	}
	if src.Auth.Token != "" {
		dst.Auth.Token = src.Auth.Token
	}
	if src.Scan.Days > 0 {
		dst.Scan.Days = src.Scan.Days
	}
	if src.Scan.Workers > 0 {
		dst.Scan.Workers = src.Scan.Workers
	}
	if src.Scan.Incremental {
		dst.Scan.Incremental = true
	}
	if src.Pricing.Source != "" {
		dst.Pricing.Source = src.Pricing.Source
	}
	if src.Pricing.LinuxPerMin > 0 {
		dst.Pricing.LinuxPerMin = src.Pricing.LinuxPerMin
	}
	if src.Pricing.WindowsMultiplier > 0 {
		dst.Pricing.WindowsMultiplier = src.Pricing.WindowsMultiplier
	}
	if src.Pricing.MacOSMultiplier > 0 {
		dst.Pricing.MacOSMultiplier = src.Pricing.MacOSMultiplier
	}
	if src.Pricing.Currency != "" {
		dst.Pricing.Currency = src.Pricing.Currency
	}
	if src.FreeTier.Plan != "" {
		dst.FreeTier.Plan = src.FreeTier.Plan
	}
	if src.FreeTier.MinutesPerMonth > 0 {
		dst.FreeTier.MinutesPerMonth = src.FreeTier.MinutesPerMonth
	}
	if src.Budget.Monthly > 0 {
		dst.Budget.Monthly = src.Budget.Monthly
	}
	if src.Budget.Weekly > 0 {
		dst.Budget.Weekly = src.Budget.Weekly
	}
	if src.Budget.Notify != "" {
		dst.Budget.Notify = src.Budget.Notify
	}
	if src.Budget.WebhookURL != "" {
		dst.Budget.WebhookURL = src.Budget.WebhookURL
	}
	if src.Output.Format != "" {
		dst.Output.Format = src.Output.Format
	}
	if src.Output.Color != "" {
		dst.Output.Color = src.Output.Color
	}
	if len(src.Ignore.Workflows) > 0 {
		dst.Ignore.Workflows = src.Ignore.Workflows
	}
}

func mergeEnv(cfg *Config) {
	if v := os.Getenv("CICOST_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Scan.Days = n
		}
	}
	if v := os.Getenv("CICOST_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Scan.Workers = n
		}
	}
	if v := os.Getenv("CICOST_INCREMENTAL"); v != "" {
		cfg.Scan.Incremental = strings.EqualFold(v, "true") || v == "1"
	}
	if v := os.Getenv("CICOST_FORMAT"); v != "" {
		cfg.Output.Format = strings.ToLower(v)
	}
	if v := os.Getenv("CICOST_MONTHLY_BUDGET"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			cfg.Budget.Monthly = f
		}
	}
}
