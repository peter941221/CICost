package cmd

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/billing"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/reconcile"
	"github.com/peter941221/CICost/internal/store"
)

func runReconcile(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("reconcile", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "目标仓库，格式 owner/repo")
	monthFlag := fs.String("month", time.Now().UTC().Format("2006-01"), "月份，格式 YYYY-MM")
	sourceFlag := fs.String("source", "csv", "账单来源 csv|github")
	inputFlag := fs.String("input", "", "账单 CSV 文件（repo,period,actual_cost_usd）")
	actualFlag := fs.Float64("actual-usd", 0, "实际账单金额，优先级高于 --input")
	applyFlag := fs.Bool("apply-calibration", false, "将校准系数应用到后续 --calibrated 报告")
	if err := fs.Parse(args); err != nil {
		return err
	}

	repo, err := pickRepo(*repoFlag, rt.cfg)
	if err != nil {
		return err
	}
	monthStart, err := time.Parse("2006-01", strings.TrimSpace(*monthFlag))
	if err != nil {
		return fmt.Errorf("invalid month %q, expected YYYY-MM", *monthFlag)
	}
	period := monthStart.Format("2006-01")
	start := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	dbPath, err := config.DBPath()
	if err != nil {
		return err
	}
	st, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer st.Close()

	jobs, err := st.ListJobs(repo, start, end)
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return fmt.Errorf("no jobs found for %s in %s", repo, period)
	}

	pcfg, err := loadPricingConfig(rt)
	if err != nil {
		return err
	}

	cost, _, _, err := analytics.CalculateCostDetailed(jobs, pcfg, 1.0)
	if err != nil {
		return err
	}
	estimated := cost.TotalCostUSD

	actual := *actualFlag
	src := strings.ToLower(strings.TrimSpace(*sourceFlag))
	if actual <= 0 {
		switch src {
		case "csv":
			if strings.TrimSpace(*inputFlag) == "" {
				return fmt.Errorf("--input is required when --source=csv and --actual-usd is not provided")
			}
			actual, err = billing.LoadActualFromCSV(*inputFlag, repo, period)
			if err != nil {
				return err
			}
		case "github":
			return fmt.Errorf("github billing source is not available yet; use --source csv or --actual-usd")
		default:
			return fmt.Errorf("unsupported source %q, expected csv|github", src)
		}
	}
	if actual <= 0 {
		return fmt.Errorf("actual cost must be > 0")
	}

	if err := st.UpsertBillingSnapshot(model.BillingSnapshot{
		Repo:          repo,
		Period:        period,
		ActualCostUSD: actual,
		Source:        src,
		FetchedAt:     time.Now().UTC(),
	}); err != nil {
		return err
	}

	res := reconcile.BuildResult(repo, period, estimated, actual)
	if !*applyFlag {
		res.CalibrationFactor = 1
	}
	if err := st.InsertReconcileResult(res); err != nil {
		return err
	}

	fmt.Printf("Reconcile Result: %s %s\n", repo, period)
	fmt.Printf("  Estimate  : $%.2f\n", res.EstimatedCostUSD)
	fmt.Printf("  Actual    : $%.2f\n", res.ActualCostUSD)
	fmt.Printf("  Delta     : %.2f%%\n", res.DeltaRatio*100)
	fmt.Printf("  Factor    : %.4f\n", res.CalibrationFactor)
	fmt.Printf("  Confidence: %s\n", res.Confidence)
	if *applyFlag {
		fmt.Println("  Calibration: enabled for future `report --calibrated`")
	} else {
		fmt.Println("  Calibration: disabled (use --apply-calibration to enable)")
	}
	return nil
}
