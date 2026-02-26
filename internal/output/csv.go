package output

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sort"
)

func RenderReportCSV(v ReportView) (string, error) {
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)

	rows := [][]string{
		{"section", "metric", "value"},
		{"summary", "repo", v.Repo},
		{"summary", "period_start", v.Start.Format("2006-01-02")},
		{"summary", "period_end", v.End.Format("2006-01-02")},
		{"summary", "days", fmt.Sprintf("%d", v.Days)},
		{"summary", "total_runs", fmt.Sprintf("%d", v.TotalRuns)},
		{"summary", "total_minutes_raw", fmt.Sprintf("%.2f", v.Cost.TotalMinutes)},
		{"summary", "total_minutes_billable", fmt.Sprintf("%.2f", v.Cost.BillableMinutes)},
		{"summary", "estimated_cost_usd", fmt.Sprintf("%.2f", v.Cost.TotalCostUSD)},
		{"pricing", "pricing_source", v.PricingSource},
		{"pricing", "pricing_snapshot_version", v.PricingSnapshotVersion},
		{"pricing", "pricing_effective_from", v.PricingEffectiveFrom},
		{"pricing", "calibrated", fmt.Sprintf("%t", v.Calibrated)},
		{"pricing", "calibration_factor", fmt.Sprintf("%.4f", v.CalibrationFactor)},
		{"waste", "fail_rate_pct", fmt.Sprintf("%.2f", v.Waste.FailRate*100)},
		{"waste", "rerun_waste_usd", fmt.Sprintf("%.2f", v.Waste.RerunWasteUSD)},
		{"waste", "cancel_waste_usd", fmt.Sprintf("%.2f", v.Waste.CancelWasteUSD)},
		{"waste", "total_waste_usd", fmt.Sprintf("%.2f", v.Waste.TotalWasteUSD)},
	}
	keys := make([]string, 0, len(v.Cost.ByOS))
	for k := range v.Cost.ByOS {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		osv := v.Cost.ByOS[k]
		rows = append(rows,
			[]string{"by_os", osv.OS + "_minutes", fmt.Sprintf("%.2f", osv.Minutes)},
			[]string{"by_os", osv.OS + "_cost_usd", fmt.Sprintf("%.2f", osv.CostUSD)},
			[]string{"by_os", osv.OS + "_cost_pct", fmt.Sprintf("%.2f", osv.Percentage)},
		)
	}
	if err := w.WriteAll(rows); err != nil {
		return "", err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
