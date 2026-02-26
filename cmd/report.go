package cmd

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/output"
	"github.com/peter941221/CICost/internal/store"
)

func runReport(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "目标仓库，格式 owner/repo")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "时间窗口天数")
	formatFlag := fs.String("format", rt.cfg.Output.Format, "输出格式 table|md|json|csv")
	outputFlag := fs.String("output", "", "输出文件")
	compareFlag := fs.Bool("compare", false, "与上一周期对比")
	calibratedFlag := fs.Bool("calibrated", false, "使用最近对账结果进行校准")
	if err := fs.Parse(args); err != nil {
		return err
	}

	repo, err := pickRepo(*repoFlag, rt.cfg)
	if err != nil {
		return err
	}
	dbPath, err := config.DBPath()
	if err != nil {
		return err
	}
	st, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer st.Close()

	start, end := calcPeriod(*daysFlag)
	runs, err := st.ListRuns(repo, start, end)
	if err != nil {
		return err
	}
	jobs, err := st.ListJobs(repo, start, end)
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		return fmt.Errorf("no data in local store for %s, run `cicost scan --repo %s` first", repo, repo)
	}

	pricingCfg, err := loadPricingConfig(rt)
	if err != nil {
		return err
	}

	cost, _, pricingMeta, err := analytics.CalculateCostDetailed(jobs, pricingCfg, 1.0)
	if err != nil {
		return err
	}
	calibrationFactor := 1.0
	calibrated := false
	if *calibratedFlag {
		if rec, ok, err := st.GetLatestReconcile(repo); err == nil && ok && rec.CalibrationFactor > 0 {
			calibrationFactor = rec.CalibrationFactor
			calibrated = true
		}
	}
	if calibrated && calibrationFactor != 1 {
		cost.TotalCostUSD = round2(cost.TotalCostUSD * calibrationFactor)
		for k, osCost := range cost.ByOS {
			osCost.CostUSD = round2(osCost.CostUSD * calibrationFactor)
			cost.ByOS[k] = osCost
		}
	}
	waste := analytics.CalculateWaste(runs, jobs, pricingCfg, cost.TotalCostUSD)
	if calibrated && calibrationFactor != 1 {
		waste.RerunWasteUSD = round2(waste.RerunWasteUSD * calibrationFactor)
		waste.CancelWasteUSD = round2(waste.CancelWasteUSD * calibrationFactor)
		waste.TotalWasteUSD = round2(waste.TotalWasteUSD * calibrationFactor)
	}
	view := output.ReportView{
		Repo:                   repo,
		Start:                  start,
		End:                    end,
		Days:                   *daysFlag,
		TotalRuns:              len(runs),
		Cost:                   cost,
		Waste:                  waste,
		PricingSnapshotVersion: pricingMeta.PricingSnapshotVersion,
		PricingEffectiveFrom:   pricingMeta.PricingEffectiveFrom.Format("2006-01-02"),
		PricingSource:          pricingMeta.PricingSource,
		Calibrated:             calibrated,
		CalibrationFactor:      calibrationFactor,
	}

	if *compareFlag {
		prevStart := start.AddDate(0, 0, -*daysFlag)
		prevEnd := start.Add(-time.Second)
		prevRuns, _ := st.ListRuns(repo, prevStart, prevEnd)
		prevJobs, _ := st.ListJobs(repo, prevStart, prevEnd)
		if len(prevRuns) > 0 && len(prevJobs) > 0 {
			prevCost, _, _, _ := analytics.CalculateCostDetailed(prevJobs, pricingCfg, 1.0)
			trend := analytics.CompareCost(cost, prevCost)
			fmt.Printf("Compare previous %d days: %s %.2f USD (%.2f%%)\n", *daysFlag, trend.Direction, trend.DeltaUSD, trend.DeltaPct)
		} else {
			fmt.Printf("Compare previous %d days: N/A (no prior data)\n", *daysFlag)
		}
	}

	var out string
	switch strings.ToLower(*formatFlag) {
	case "md", "markdown":
		out = output.RenderReportMarkdown(view)
	case "json":
		out, err = output.RenderReportJSON(view, version)
		if err != nil {
			return err
		}
	case "csv":
		out, err = output.RenderReportCSV(view)
		if err != nil {
			return err
		}
	default:
		out = output.RenderReportTable(view)
	}
	if err := writeOutput(*outputFlag, out); err != nil {
		return err
	}
	return nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
