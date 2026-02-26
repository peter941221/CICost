package pricing

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	PricingSourceSKU    = "sku_direct"
	PricingSourceLegacy = "legacy_multiplier"
)

type Config struct {
	Version             string
	EffectiveFrom       time.Time
	PerMinuteUSD        float64
	WindowsMultiplier   float64
	MacOSMultiplier     float64
	FreeTierPerMonth    float64
	AlreadyUsedThisMon  float64
	FreeTierByPlan      map[string]float64
	LargerRunnersPerMin map[string]float64
	Snapshots           []Snapshot
}

type Snapshot struct {
	Version       string
	EffectiveFrom time.Time
	SKUs          map[string]float64
}

type JobPrice struct {
	BillableMinutes float64
	RatePerMin      float64
	CostUSD         float64
	Source          string
	SKU             string
	Snapshot        Snapshot
}

func ceilMinutes(durationSec int) float64 {
	if durationSec <= 0 {
		return 0
	}
	return math.Ceil(float64(durationSec) / 60)
}

func RawBillableMinutes(durationSec int) float64 {
	return ceilMinutes(durationSec)
}

func LegacyMultiplier(runnerOS string, cfg Config) float64 {
	switch runnerOS {
	case "Windows":
		if cfg.WindowsMultiplier > 0 {
			return cfg.WindowsMultiplier
		}
		return 2
	case "macOS":
		if cfg.MacOSMultiplier > 0 {
			return cfg.MacOSMultiplier
		}
		return 10
	default:
		return 1
	}
}

func BillableMinutes(durationSec int, runnerOS string, cfg Config) float64 {
	raw := RawBillableMinutes(durationSec)
	return raw * LegacyMultiplier(runnerOS, cfg)
}

func SelectSnapshot(cfg Config, at time.Time) (Snapshot, error) {
	if len(cfg.Snapshots) == 0 {
		return Snapshot{}, fmt.Errorf("pricing snapshots not configured")
	}
	ref := at.UTC()
	var chosen *Snapshot
	for i := range cfg.Snapshots {
		s := cfg.Snapshots[i]
		if s.EffectiveFrom.IsZero() {
			continue
		}
		if s.EffectiveFrom.After(ref) {
			continue
		}
		if chosen == nil || s.EffectiveFrom.After(chosen.EffectiveFrom) {
			ss := s
			chosen = &ss
		}
	}
	if chosen == nil {
		earliest := cfg.Snapshots[0].EffectiveFrom.Format("2006-01-02")
		return Snapshot{}, fmt.Errorf("no pricing snapshot matched %s; earliest available effective_from is %s (add an older pricing_snapshots entry)", ref.Format("2006-01-02"), earliest)
	}
	return *chosen, nil
}

func ResolveRate(cfg Config, at time.Time, runnerOS, runnerName string) (rate float64, sku string, source string, snapshot Snapshot, err error) {
	if len(cfg.Snapshots) > 0 {
		snapshot, err = SelectSnapshot(cfg, at)
		if err != nil {
			return 0, "", "", Snapshot{}, err
		}
		candidates := skuCandidates(runnerOS, runnerName)
		for _, k := range candidates {
			if v, ok := snapshot.SKUs[k]; ok && v > 0 {
				return v, k, PricingSourceSKU, snapshot, nil
			}
		}
		return 0, "", "", Snapshot{}, fmt.Errorf("no SKU rate matched runner_os=%q runner_name=%q under snapshot %s", runnerOS, runnerName, snapshot.Version)
	}

	base := cfg.PerMinuteUSD
	if base <= 0 {
		base = 0.008
	}
	mult := LegacyMultiplier(runnerOS, cfg)
	return base * mult, normalizeSKU(runnerOS), PricingSourceLegacy, Snapshot{
		Version:       cfg.Version,
		EffectiveFrom: cfg.EffectiveFrom,
	}, nil
}

func skuCandidates(runnerOS, runnerName string) []string {
	out := []string{}
	appendIf := func(v string) {
		if v == "" {
			return
		}
		for _, existing := range out {
			if existing == v {
				return
			}
		}
		out = append(out, v)
	}
	appendIf(normalizeSKU(runnerName))
	appendIf(normalizeSKU(runnerOS))

	switch strings.ToLower(strings.TrimSpace(runnerOS)) {
	case "linux":
		appendIf("ubuntu-latest")
	case "windows":
		appendIf("windows-latest")
	case "macos":
		appendIf("macos-latest")
	}
	if len(out) == 0 {
		out = append(out, "linux")
	}
	return out
}

func normalizeSKU(v string) string {
	s := strings.TrimSpace(strings.ToLower(v))
	s = strings.ReplaceAll(s, "_", "-")
	switch s {
	case "macos", "macos-latest", "mac-os":
		return "macos"
	case "windows", "windows-latest":
		return "windows"
	case "linux", "ubuntu", "ubuntu-latest":
		return "linux"
	}
	return s
}

func PriceJob(durationSec int, runnerOS, runnerName string, startedAt time.Time, cfg Config) (JobPrice, error) {
	ref := startedAt
	if ref.IsZero() {
		ref = time.Now().UTC()
	}
	rate, sku, source, snapshot, err := ResolveRate(cfg, ref, runnerOS, runnerName)
	if err != nil {
		return JobPrice{}, err
	}
	billable := RawBillableMinutes(durationSec)
	if source == PricingSourceLegacy {
		billable = BillableMinutes(durationSec, runnerOS, cfg)
	}
	return JobPrice{
		BillableMinutes: billable,
		RatePerMin:      rate,
		CostUSD:         billable * rate,
		Source:          source,
		SKU:             sku,
		Snapshot:        snapshot,
	}, nil
}
