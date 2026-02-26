package analytics

import (
	"math"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/pricing"
)

type CostPricingMeta struct {
	PricingSource          string
	PricingSnapshotVersion string
	PricingEffectiveFrom   time.Time
}

func CalculateCostDetailed(jobs []model.Job, cfg pricing.Config, completeness float64) (model.CostResult, map[int64]float64, CostPricingMeta, error) {
	result := model.CostResult{
		ByOS:             map[string]model.OSCost{},
		DataCompleteness: completeness,
		Disclaimer:       "Estimate only. Free tier is shared across repositories in an account/org.",
	}
	jobCost := make(map[int64]float64, len(jobs))
	preFreeTotalCost := 0.0
	meta := CostPricingMeta{}
	firstMeta := true

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
		quote, err := pricing.PriceJob(job.DurationSec, job.RunnerOS, job.RunnerName, job.StartedAt, cfg)
		if err != nil {
			return model.CostResult{}, nil, CostPricingMeta{}, err
		}
		billable := quote.BillableMinutes
		cost := quote.CostUSD

		result.TotalMinutes += rawMin
		result.BillableMinutes += billable
		jobCost[job.ID] = cost
		preFreeTotalCost += cost

		if firstMeta {
			meta.PricingSource = quote.Source
			meta.PricingSnapshotVersion = quote.Snapshot.Version
			meta.PricingEffectiveFrom = quote.Snapshot.EffectiveFrom
			firstMeta = false
		} else {
			if meta.PricingSource != quote.Source {
				meta.PricingSource = "mixed"
			}
			if meta.PricingSnapshotVersion != quote.Snapshot.Version {
				meta.PricingSnapshotVersion = "mixed"
				meta.PricingEffectiveFrom = time.Time{}
			}
		}

		osCost := result.ByOS[job.RunnerOS]
		osCost.OS = job.RunnerOS
		osCost.Minutes += rawMin
		osCost.CostUSD += cost
		if osCost.Multiplier == 0 {
			if quote.Source == pricing.PricingSourceLegacy {
				osCost.Multiplier = pricing.LegacyMultiplier(job.RunnerOS, cfg)
			} else {
				osCost.Multiplier = 1
			}
		}
		result.ByOS[job.RunnerOS] = osCost
	}

	charged := pricing.ChargedMinutes(result.BillableMinutes, cfg.FreeTierPerMonth, cfg.AlreadyUsedThisMon)
	result.FreeTierUsed = math.Min(cfg.FreeTierPerMonth, result.BillableMinutes)
	if result.BillableMinutes > 0 && preFreeTotalCost > 0 {
		effectiveRate := preFreeTotalCost / result.BillableMinutes
		result.TotalCostUSD = round2(charged * effectiveRate)
	}
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
	return result, jobCost, meta, nil
}

func CalculateCost(jobs []model.Job, cfg pricing.Config, completeness float64) (model.CostResult, map[int64]float64) {
	res, costs, _, err := CalculateCostDetailed(jobs, cfg, completeness)
	if err != nil {
		return model.CostResult{
			ByOS:             map[string]model.OSCost{},
			DataCompleteness: completeness,
			Disclaimer:       "Estimate failed due to pricing configuration error.",
		}, map[int64]float64{}
	}
	return res, costs
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
