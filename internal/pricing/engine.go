package pricing

import "math"

type Config struct {
	Version             string
	PerMinuteUSD        float64
	WindowsMultiplier   float64
	MacOSMultiplier     float64
	FreeTierPerMonth    float64
	AlreadyUsedThisMon  float64
	FreeTierByPlan      map[string]float64
	LargerRunnersPerMin map[string]float64
}

func ceilMinutes(durationSec int) float64 {
	if durationSec <= 0 {
		return 0
	}
	return math.Ceil(float64(durationSec) / 60)
}

func BillableMinutes(durationSec int, runnerOS string, cfg Config) float64 {
	raw := ceilMinutes(durationSec)
	switch runnerOS {
	case "Windows":
		if cfg.WindowsMultiplier > 0 {
			return raw * cfg.WindowsMultiplier
		}
		return raw * 2
	case "macOS":
		if cfg.MacOSMultiplier > 0 {
			return raw * cfg.MacOSMultiplier
		}
		return raw * 10
	default:
		return raw
	}
}
