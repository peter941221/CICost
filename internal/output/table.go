package output

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/model"
)

type ReportView struct {
	Repo      string               `json:"repo"`
	Start     time.Time            `json:"start"`
	End       time.Time            `json:"end"`
	Days      int                  `json:"days"`
	TotalRuns int                  `json:"total_runs"`
	Cost      model.CostResult     `json:"cost"`
	Waste     model.WasteMetrics   `json:"waste"`
	Hotspots  []model.HotspotEntry `json:"hotspots,omitempty"`
}

func RenderReportTable(v ReportView) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CICost Report: %s\n", v.Repo)
	fmt.Fprintf(&b, "Period: %s ~ %s (%d days, UTC)\n\n", v.Start.Format("2006-01-02"), v.End.Format("2006-01-02"), v.Days)
	fmt.Fprintf(&b, "SUMMARY\n")
	fmt.Fprintf(&b, "  Total Runs: %d\n", v.TotalRuns)
	fmt.Fprintf(&b, "  Total Minutes (raw): %.2f\n", v.Cost.TotalMinutes)
	fmt.Fprintf(&b, "  Total Minutes (billable): %.2f\n", v.Cost.BillableMinutes)
	fmt.Fprintf(&b, "  Estimated Cost: $%.2f\n", v.Cost.TotalCostUSD)
	fmt.Fprintf(&b, "  Free Tier Used: %.2f min\n\n", v.Cost.FreeTierUsed)

	fmt.Fprintf(&b, "WASTE\n")
	fmt.Fprintf(&b, "  Fail Rate: %.1f%%\n", v.Waste.FailRate*100)
	fmt.Fprintf(&b, "  Rerun Waste: $%.2f\n", v.Waste.RerunWasteUSD)
	fmt.Fprintf(&b, "  Cancel Waste: $%.2f\n", v.Waste.CancelWasteUSD)
	fmt.Fprintf(&b, "  Total Waste: $%.2f (%.1f%% of total)\n\n", v.Waste.TotalWasteUSD, v.Waste.WastePercentage)

	fmt.Fprintf(&b, "BY OS\n")
	keys := make([]string, 0, len(v.Cost.ByOS))
	for k := range v.Cost.ByOS {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		osv := v.Cost.ByOS[k]
		fmt.Fprintf(&b, "  %-8s minutes=%8.2f cost=$%8.2f pct=%6.2f%%\n", osv.OS, osv.Minutes, osv.CostUSD, osv.Percentage)
	}
	fmt.Fprintf(&b, "\nData completeness: %.1f%%\n", v.Cost.DataCompleteness*100)
	fmt.Fprintf(&b, "Disclaimer: %s\n", v.Cost.Disclaimer)
	return b.String()
}

func RenderHotspotsTable(repo string, days int, groupBy string, entries []model.HotspotEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Top %d %s hotspots for %s (last %d days)\n", len(entries), groupBy, repo, days)
	fmt.Fprintf(&b, "Rank  Name                                   Minutes    Cost($)   Cost%%   Fail%%\n")
	fmt.Fprintf(&b, "----  -------------------------------------  ---------  --------  ------  ------\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "%-4d  %-37.37s  %9.2f  %8.2f  %6.2f  %6.2f\n",
			e.Rank, e.Name, e.Minutes, e.CostUSD, e.CostPct, e.FailRate)
	}
	return b.String()
}
