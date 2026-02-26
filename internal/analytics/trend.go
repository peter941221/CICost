package analytics

import "github.com/peter941221/CICost/internal/model"

type Trend struct {
	Direction string  `json:"direction"`
	DeltaUSD  float64 `json:"delta_usd"`
	DeltaPct  float64 `json:"delta_pct"`
}

func CompareCost(current, previous model.CostResult) Trend {
	if previous.TotalCostUSD <= 0 {
		return Trend{Direction: "N/A", DeltaUSD: round2(current.TotalCostUSD), DeltaPct: 0}
	}
	delta := current.TotalCostUSD - previous.TotalCostUSD
	dir := "→"
	if delta > 0.01 {
		dir = "↑"
	} else if delta < -0.01 {
		dir = "↓"
	}
	return Trend{
		Direction: dir,
		DeltaUSD:  round2(delta),
		DeltaPct:  round2((delta / previous.TotalCostUSD) * 100),
	}
}
