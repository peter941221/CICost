package analytics

import "time"

type BudgetStatus string

const (
	BudgetOK       BudgetStatus = "ok"
	BudgetWarning  BudgetStatus = "warning"
	BudgetExceeded BudgetStatus = "exceeded"
)

type BudgetResult struct {
	Status         BudgetStatus `json:"status"`
	ThresholdUSD   float64      `json:"threshold_usd"`
	ActualUSD      float64      `json:"actual_usd"`
	ProjectedUSD   float64      `json:"projected_usd"`
	PercentageUsed float64      `json:"percentage_used"`
	PeriodStart    time.Time    `json:"period_start"`
	PeriodEnd      time.Time    `json:"period_end"`
	CheckType      string       `json:"check_type"`
}

func EvaluateBudget(now time.Time, actual, threshold float64, checkType string) BudgetResult {
	start, end := periodBounds(now, checkType)
	elapsedDays := now.Sub(start).Hours() / 24
	totalDays := end.Sub(start).Hours() / 24
	if elapsedDays < 1 {
		elapsedDays = 1
	}
	projected := actual
	if totalDays > 0 {
		projected = actual * (totalDays / elapsedDays)
	}
	res := BudgetResult{
		Status:         BudgetOK,
		ThresholdUSD:   round2(threshold),
		ActualUSD:      round2(actual),
		ProjectedUSD:   round2(projected),
		PeriodStart:    start,
		PeriodEnd:      end,
		CheckType:      checkType,
		PercentageUsed: 0,
	}
	if threshold > 0 {
		res.PercentageUsed = round2((actual / threshold) * 100)
	}
	if threshold > 0 && actual > threshold {
		res.Status = BudgetExceeded
	} else if threshold > 0 && projected > threshold {
		res.Status = BudgetWarning
	}
	return res
}

func periodBounds(now time.Time, checkType string) (time.Time, time.Time) {
	n := now.UTC()
	if checkType == "weekly" {
		wd := int(n.Weekday())
		if wd == 0 {
			wd = 7
		}
		start := time.Date(n.Year(), n.Month(), n.Day()-(wd-1), 0, 0, 0, 0, time.UTC)
		return start, start.AddDate(0, 0, 7)
	}
	start := time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return start, end
}

func PeriodBounds(now time.Time, checkType string) (time.Time, time.Time) {
	return periodBounds(now, checkType)
}
