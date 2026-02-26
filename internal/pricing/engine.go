package pricing

import "math"

type Config struct {
	PerMinuteUSD       float64
	WindowsMultiplier  float64
	MacOSMultiplier    float64
	FreeTierPerMonth   float64
	AlreadyUsedThisMon float64
}

func ceilMinutes(durationSec int) float64 {
	if durationSec <= 0 {
		return 0
	}
	return math.Ceil(float64(durationSec) / 60)
}

func BillableMinutes(durationSec int, runnerOS string) float64 {
	raw := ceilMinutes(durationSec)
	switch runnerOS {
	case "Windows":
		return raw * 2
	case "macOS":
		return raw * 10
	default:
		return raw
	}
}

