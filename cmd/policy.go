package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/policy"
	"github.com/peter941221/CICost/internal/store"
)

func runPolicy(args []string) error {
	if len(args) == 0 {
		return runPolicyHelp()
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "check":
		return runPolicyCheck(args[1:])
	case "lint":
		return runPolicyLint(args[1:])
	case "explain":
		return runPolicyExplain(args[1:])
	default:
		return fmt.Errorf("unknown policy subcommand %q, expected check|lint|explain", sub)
	}
}

func runPolicyHelp() error {
	fmt.Println(`Usage:
  cicost policy check --repo owner/repo --days 30 [--policy .cicost.policy.yml]
  cicost policy lint [--policy .cicost.policy.yml]
  cicost policy explain`)
	return nil
}

func runPolicyLint(args []string) error {
	fs := flag.NewFlagSet("policy lint", flag.ContinueOnError)
	policyPath := fs.String("policy", ".cicost.policy.yml", "Policy file path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	p := resolvePolicyPath(*policyPath)
	cfg, err := policy.LoadFromFile(p)
	if err != nil {
		return err
	}
	if err := policy.Lint(cfg); err != nil {
		return err
	}
	fmt.Printf("Policy lint passed: %s\n", p)
	return nil
}

func runPolicyCheck(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("policy check", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "Target repository in owner/repo format")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "Time window in days")
	policyPath := fs.String("policy", ".cicost.policy.yml", "Policy file path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	repo, err := pickRepo(*repoFlag, rt.cfg)
	if err != nil {
		return err
	}
	p := resolvePolicyPath(*policyPath)
	cfg, err := policy.LoadFromFile(p)
	if err != nil {
		return err
	}
	if err := policy.Lint(cfg); err != nil {
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
	pcfg, err := loadPricingConfig(rt)
	if err != nil {
		return err
	}
	cost, _, _, err := analytics.CalculateCostDetailed(jobs, pcfg, 1.0)
	if err != nil {
		return err
	}
	waste := analytics.CalculateWaste(runs, jobs, pcfg, cost.TotalCostUSD)
	metrics := map[string]float64{
		"monthly_cost_usd": cost.TotalCostUSD,
		"total_cost_usd":   cost.TotalCostUSD,
		"waste_percentage": waste.WastePercentage,
		"fail_rate":        waste.FailRate * 100,
		"total_runs":       float64(len(runs)),
	}

	findings, err := policy.Evaluate(cfg, metrics)
	if err != nil {
		return err
	}
	if len(findings) == 0 {
		fmt.Println("Policy check: no rules matched.")
		return nil
	}

	fmt.Printf("Policy check findings for %s (%d days)\n", repo, *daysFlag)
	hasError := false
	for _, f := range findings {
		fmt.Printf("- [%s] %s (evidence: %s=%.4f, when: %s)\n", strings.ToUpper(string(f.Severity)), f.RuleID, f.EvidenceKey, f.EvidenceValue, f.When)
		if err := st.InsertPolicyRun(model.PolicyRun{
			Repo:          repo,
			PeriodStart:   start,
			PeriodEnd:     end,
			RuleID:        f.RuleID,
			Severity:      string(f.Severity),
			Matched:       true,
			EvidenceKey:   f.EvidenceKey,
			EvidenceValue: f.EvidenceValue,
			Expression:    f.When,
			CreatedAt:     time.Now().UTC(),
		}); err != nil {
			return err
		}
		if f.Severity == policy.SeverityError {
			hasError = true
		}
	}

	if hasError {
		return withExit(3, fmt.Errorf("policy check failed: one or more error rules matched"))
	}
	return nil
}

func runPolicyExplain(_ []string) error {
	fmt.Println(`Policy explain
- supported metrics:
  - monthly_cost_usd
  - total_cost_usd
  - waste_percentage
  - fail_rate
  - total_runs
- supported operators: >, >=, <, <=, ==, !=

example:
rules:
  - id: budget_monthly
    when: monthly_cost_usd > 200
    severity: error`)
	return nil
}

func resolvePolicyPath(path string) string {
	if strings.TrimSpace(path) == "" {
		path = ".cicost.policy.yml"
	}
	candidates := []string{path}
	if path == ".cicost.policy.yml" {
		candidates = append(candidates, filepath.Join("..", path))
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return path
}
