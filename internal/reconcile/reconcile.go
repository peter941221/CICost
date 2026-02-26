package reconcile

import (
	"math"
	"time"

	"github.com/peter941221/CICost/internal/model"
)

func BuildResult(repo, period string, estimate, actual float64) model.ReconcileResult {
	delta := 0.0
	if actual != 0 {
		delta = (estimate - actual) / actual
	}
	factor := 1.0
	if estimate > 0 {
		factor = actual / estimate
	}
	if factor <= 0 {
		factor = 1
	}
	return model.ReconcileResult{
		Repo:              repo,
		Period:            period,
		EstimatedCostUSD:  round2(estimate),
		ActualCostUSD:     round2(actual),
		DeltaRatio:        round4(delta),
		CalibrationFactor: round4(factor),
		Confidence:        Confidence(delta),
		CreatedAt:         time.Now().UTC(),
	}
}

func Confidence(delta float64) string {
	absDelta := math.Abs(delta)
	switch {
	case absDelta <= 0.05:
		return "high"
	case absDelta <= 0.15:
		return "medium"
	default:
		return "low"
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
