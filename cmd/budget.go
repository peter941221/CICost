package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/pricing"
	"github.com/peter941221/CICost/internal/store"
)

func runBudget(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("budget", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "目标仓库，格式 owner/repo")
	monthlyFlag := fs.Float64("monthly", 0, "月预算阈值（USD）")
	weeklyFlag := fs.Float64("weekly", 0, "周预算阈值（USD）")
	notifyFlag := fs.String("notify", rt.cfg.Budget.Notify, "通知方式 stdout|webhook|file")
	webhookFlag := fs.String("webhook-url", rt.cfg.Budget.WebhookURL, "Webhook URL")
	outputFlag := fs.String("output", "", "文件输出路径")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := pickRepo(*repoFlag, rt.cfg)
	if err != nil {
		return err
	}

	checkType := "monthly"
	threshold := rt.cfg.Budget.Monthly
	if *monthlyFlag > 0 {
		threshold = *monthlyFlag
	}
	if *weeklyFlag > 0 {
		checkType = "weekly"
		threshold = *weeklyFlag
	}
	if threshold <= 0 {
		return fmt.Errorf("threshold is required: set --monthly or --weekly")
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

	now := time.Now().UTC()
	start, _ := analytics.PeriodBounds(now, checkType)
	runs, err := st.ListRuns(repo, start, now)
	if err != nil {
		return err
	}
	jobs, err := st.ListJobs(repo, start, now)
	if err != nil {
		return err
	}
	pcfg, _ := pricing.LoadFromFile("configs/pricing_default.yml")
	if pcfg.PerMinuteUSD == 0 {
		pcfg.PerMinuteUSD = rt.cfg.Pricing.LinuxPerMin
	}
	if pcfg.WindowsMultiplier == 0 {
		pcfg.WindowsMultiplier = rt.cfg.Pricing.WindowsMultiplier
	}
	if pcfg.MacOSMultiplier == 0 {
		pcfg.MacOSMultiplier = rt.cfg.Pricing.MacOSMultiplier
	}
	pcfg.FreeTierPerMonth = rt.cfg.FreeTier.MinutesPerMonth
	cost, _ := analytics.CalculateCost(jobs, pcfg, 1.0)
	result := analytics.EvaluateBudget(now, cost.TotalCostUSD, threshold, checkType)
	top := analytics.CalculateHotspots(runs, jobs, pcfg, analytics.HotspotOptions{
		GroupBy: "workflow",
		TopN:    3,
		SortBy:  "cost",
	})

	_ = st.InsertBudgetCheck(repo, checkType, start, now, threshold, cost.TotalCostUSD, result.Status != analytics.BudgetOK)

	msg := fmt.Sprintf("Budget %s for %s\n  Period   : %s ~ %s\n  Budget   : $%.2f\n  Actual   : $%.2f\n  Projected: $%.2f\n",
		strings.ToUpper(string(result.Status)), repo, start.Format("2006-01-02"), now.Format("2006-01-02"),
		result.ThresholdUSD, result.ActualUSD, result.ProjectedUSD)
	if len(top) > 0 {
		msg += "  Top Contributors:\n"
		for i, e := range top {
			msg += fmt.Sprintf("    %d. %s $%.2f\n", i+1, e.Name, e.CostUSD)
		}
	}

	switch strings.ToLower(*notifyFlag) {
	case "file":
		target := *outputFlag
		if target == "" {
			target = "budget.txt"
		}
		if err := writeOutput(target, msg); err != nil {
			return err
		}
	case "webhook":
		if strings.TrimSpace(*webhookFlag) == "" {
			return fmt.Errorf("webhook-url is required when notify=webhook")
		}
		payload := map[string]any{
			"status":           string(result.Status),
			"repo":             repo,
			"budget_usd":       result.ThresholdUSD,
			"actual_usd":       result.ActualUSD,
			"projected_usd":    result.ProjectedUSD,
			"top_contributors": top,
		}
		b, _ := json.Marshal(payload)
		resp, err := http.Post(*webhookFlag, "application/json", bytes.NewReader(b))
		if err != nil {
			fmt.Printf("WARN: webhook send failed: %v\n", err)
		} else {
			_ = resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				fmt.Printf("WARN: webhook returned status %d\n", resp.StatusCode)
			}
		}
		fmt.Print(msg)
	default:
		fmt.Print(msg)
	}

	if result.Status == analytics.BudgetExceeded || result.Status == analytics.BudgetWarning {
		return withExit(2, fmt.Errorf("budget %s", result.Status))
	}
	return nil
}
