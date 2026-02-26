package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/pricing"
	"github.com/peter941221/CICost/internal/store"
)

func runExplain(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "目标仓库，格式 owner/repo")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "时间窗口天数")
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
	if len(runs) == 0 || len(jobs) == 0 {
		fmt.Println("No significant optimization opportunities found (insufficient local data).")
		return nil
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
	cost, _ := analytics.CalculateCost(jobs, pcfg, 1.0)
	waste := analytics.CalculateWaste(runs, jobs, pcfg, cost.TotalCostUSD)
	hotspots := analytics.CalculateHotspots(runs, jobs, pcfg, analytics.HotspotOptions{
		GroupBy: "workflow",
		TopN:    5,
		SortBy:  "cost",
	})

	var tips []string
	if mac, ok := cost.ByOS["macOS"]; ok && cost.TotalCostUSD > 0 && (mac.CostUSD/cost.TotalCostUSD) > 0.2 {
		saving := mac.CostUSD * 0.9
		tips = append(tips, fmt.Sprintf("HIGH 迁移 macOS 重负载任务到 Linux，预计可节省约 $%.2f /周期。", saving))
	}
	if waste.FailRate > 0.15 {
		tips = append(tips, fmt.Sprintf("MEDIUM 失败率为 %.1f%%，建议优先排查 flaky tests 与不稳定依赖。", waste.FailRate*100))
	}
	if waste.CancelWasteUSD > 0 {
		tips = append(tips, fmt.Sprintf("MEDIUM 检测到取消浪费 $%.2f，建议给 PR workflow 增加 concurrency + cancel-in-progress。", waste.CancelWasteUSD))
	}
	if len(hotspots) > 0 {
		sort.Slice(hotspots, func(i, j int) bool { return hotspots[i].CostUSD > hotspots[j].CostUSD })
		top := hotspots[0]
		if top.FailRate > 20 {
			tips = append(tips, fmt.Sprintf("HIGH 工作流 `%s` 成本最高且失败率 %.1f%%，优先治理这个热区。", top.Name, top.FailRate))
		} else {
			tips = append(tips, fmt.Sprintf("LOW 工作流 `%s` 占成本 %.1f%%，可先做缓存与路径过滤优化。", top.Name, top.CostPct))
		}
	}
	if len(tips) == 0 {
		fmt.Println("No significant optimization opportunities found.")
		return nil
	}

	fmt.Printf("CICost Recommendations for %s\n", repo)
	fmt.Println(strings.Repeat("=", 50))
	for i, tip := range tips {
		fmt.Printf("%d. %s\n", i+1, tip)
	}
	return nil
}
