package cmd

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/output"
	"github.com/peter941221/CICost/internal/pricing"
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

	pricingCfg, _ := pricing.LoadFromFile("configs/pricing_default.yml")
	if pricingCfg.PerMinuteUSD == 0 {
		pricingCfg.PerMinuteUSD = rt.cfg.Pricing.LinuxPerMin
	}
	pricingCfg.FreeTierPerMonth = rt.cfg.FreeTier.MinutesPerMonth
	if rt.cfg.Pricing.WindowsMultiplier > 0 {
		pricingCfg.WindowsMultiplier = rt.cfg.Pricing.WindowsMultiplier
	}
	if rt.cfg.Pricing.MacOSMultiplier > 0 {
		pricingCfg.MacOSMultiplier = rt.cfg.Pricing.MacOSMultiplier
	}

	cost, _ := analytics.CalculateCost(jobs, pricingCfg, 1.0)
	waste := analytics.CalculateWaste(runs, jobs, pricingCfg, cost.TotalCostUSD)
	view := output.ReportView{
		Repo:      repo,
		Start:     start,
		End:       end,
		Days:      *daysFlag,
		TotalRuns: len(runs),
		Cost:      cost,
		Waste:     waste,
	}

	if *compareFlag {
		prevStart := start.AddDate(0, 0, -*daysFlag)
		prevEnd := start.Add(-time.Second)
		prevRuns, _ := st.ListRuns(repo, prevStart, prevEnd)
		prevJobs, _ := st.ListJobs(repo, prevStart, prevEnd)
		if len(prevRuns) > 0 && len(prevJobs) > 0 {
			prevCost, _ := analytics.CalculateCost(prevJobs, pricingCfg, 1.0)
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
