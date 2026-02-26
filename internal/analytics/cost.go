package analytics

import (
	"math"
	"strings"

	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/pricing"
)

func CalculateCost(jobs []model.Job, cfg pricing.Config, completeness float64) (model.CostResult, map[int64]float64) {
	result := model.CostResult{
		ByOS:             map[string]model.OSCost{},
		DataCompleteness: completeness,
		Disclaimer:       "Estimate only. Free tier is shared across repositories in an account/org.",
	}
	jobCost := make(map[int64]float64, len(jobs))
	preFreeTotalCost := 0.0

	for _, job := range jobs {
		if strings.TrimSpace(job.Status) != "completed" {
			continue
		}
		if job.IsSelfHosted {
			continue
		}
		rawMin := round2(float64(job.DurationSec) / 60)
		if rawMin < 0 {
			rawMin = 0
		}
		billable := pricing.BillableMinutes(job.DurationSec, job.RunnerOS, cfg)
		cost := billable * cfg.PerMinuteUSD

		result.TotalMinutes += rawMin
		result.BillableMinutes += billable
		jobCost[job.ID] = cost
		preFreeTotalCost += cost

		osCost := result.ByOS[job.RunnerOS]
		osCost.OS = job.RunnerOS
		osCost.Minutes += rawMin
		osCost.CostUSD += cost
		if osCost.Multiplier == 0 {
			switch job.RunnerOS {
			case "Windows":
				osCost.Multiplier = cfg.WindowsMultiplier
			case "macOS":
				osCost.Multiplier = cfg.MacOSMultiplier
			default:
				osCost.Multiplier = 1
			}
		}
		result.ByOS[job.RunnerOS] = osCost
	}

	charged := pricing.ChargedMinutes(result.BillableMinutes, cfg.FreeTierPerMonth, cfg.AlreadyUsedThisMon)
	result.FreeTierUsed = math.Min(cfg.FreeTierPerMonth, result.BillableMinutes)
	result.TotalCostUSD = round2(charged * cfg.PerMinuteUSD)
	if result.BillableMinutes <= 0 {
		result.TotalCostUSD = 0
	}

	for k, osCost := range result.ByOS {
		if preFreeTotalCost > 0 {
			osCost.Percentage = round2((osCost.CostUSD / preFreeTotalCost) * 100)
		}
		osCost.CostUSD = round2(osCost.CostUSD)
		osCost.Minutes = round2(osCost.Minutes)
		result.ByOS[k] = osCost
	}
	result.TotalMinutes = round2(result.TotalMinutes)
	result.BillableMinutes = round2(result.BillableMinutes)
	return result, jobCost
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
