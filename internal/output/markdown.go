package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/peter941221/CICost/internal/model"
)

func RenderReportMarkdown(v ReportView) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# CICost Report: %s\n\n", v.Repo)
	fmt.Fprintf(&b, "- Period: `%s ~ %s` (%d days, UTC)\n", v.Start.Format("2006-01-02"), v.End.Format("2006-01-02"), v.Days)
	fmt.Fprintf(&b, "- Total Runs: `%d`\n\n", v.TotalRuns)

	fmt.Fprintf(&b, "## Summary\n\n")
	fmt.Fprintf(&b, "| Metric | Value |\n|---|---:|\n")
	fmt.Fprintf(&b, "| Total Minutes (raw) | %.2f |\n", v.Cost.TotalMinutes)
	fmt.Fprintf(&b, "| Total Minutes (billable) | %.2f |\n", v.Cost.BillableMinutes)
	fmt.Fprintf(&b, "| Estimated Cost (USD) | $%.2f |\n", v.Cost.TotalCostUSD)
	fmt.Fprintf(&b, "| Free Tier Used | %.2f min |\n\n", v.Cost.FreeTierUsed)

	fmt.Fprintf(&b, "## Waste\n\n")
	fmt.Fprintf(&b, "| Metric | Value |\n|---|---:|\n")
	fmt.Fprintf(&b, "| Fail Rate | %.1f%% |\n", v.Waste.FailRate*100)
	fmt.Fprintf(&b, "| Rerun Waste | $%.2f |\n", v.Waste.RerunWasteUSD)
	fmt.Fprintf(&b, "| Cancel Waste | $%.2f |\n", v.Waste.CancelWasteUSD)
	fmt.Fprintf(&b, "| Total Waste | $%.2f |\n", v.Waste.TotalWasteUSD)
	fmt.Fprintf(&b, "| Waste Percentage | %.1f%% |\n\n", v.Waste.WastePercentage)

	fmt.Fprintf(&b, "## By OS\n\n")
	fmt.Fprintf(&b, "| OS | Minutes | Cost(USD) | Cost %% |\n|---|---:|---:|---:|\n")
	keys := make([]string, 0, len(v.Cost.ByOS))
	for k := range v.Cost.ByOS {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		osv := v.Cost.ByOS[k]
		fmt.Fprintf(&b, "| %s | %.2f | %.2f | %.2f%% |\n", osv.OS, osv.Minutes, osv.CostUSD, osv.Percentage)
	}
	fmt.Fprintf(&b, "\n> Data completeness: %.1f%%\n", v.Cost.DataCompleteness*100)
	fmt.Fprintf(&b, "> Disclaimer: %s\n", v.Cost.Disclaimer)
	return b.String()
}

func RenderHotspotsMarkdown(repo string, days int, groupBy string, entries []model.HotspotEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Top %d %s hotspots for `%s` (last %d days)\n\n", len(entries), groupBy, repo, days)
	fmt.Fprintf(&b, "| Rank | Name | Minutes | Cost(USD) | Cost %% | Fail %% |\n|---:|---|---:|---:|---:|---:|\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "| %d | %s | %.2f | %.2f | %.2f%% | %.2f%% |\n", e.Rank, e.Name, e.Minutes, e.CostUSD, e.CostPct, e.FailRate)
	}
	return b.String()
}
