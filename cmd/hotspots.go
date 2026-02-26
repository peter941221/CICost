package cmd

import (
	"flag"
	"strings"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/output"
	"github.com/peter941221/CICost/internal/store"
)

func runHotspots(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("hotspots", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "目标仓库，格式 owner/repo")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "时间窗口天数")
	groupByFlag := fs.String("group-by", "workflow", "维度 workflow|job|runner|branch")
	topFlag := fs.Int("top", 10, "显示前 N 条")
	sortFlag := fs.String("sort", "cost", "排序 cost|minutes|fail_rate")
	formatFlag := fs.String("format", "table", "输出格式 table|md|json")
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
	pcfg, err := loadPricingConfig(rt)
	if err != nil {
		return err
	}
	entries := analytics.CalculateHotspots(runs, jobs, pcfg, analytics.HotspotOptions{
		GroupBy: *groupByFlag,
		TopN:    *topFlag,
		SortBy:  *sortFlag,
	})

	switch strings.ToLower(*formatFlag) {
	case "json":
		view := output.ReportView{Repo: repo, Start: start, End: end, Days: *daysFlag, Hotspots: entries}
		s, err := output.RenderReportJSON(view, version)
		if err != nil {
			return err
		}
		return writeOutput("", s)
	case "md":
		md := output.RenderHotspotsMarkdown(repo, *daysFlag, *groupByFlag, entries)
		return writeOutput("", md)
	default:
		tbl := output.RenderHotspotsTable(repo, *daysFlag, *groupByFlag, entries)
		return writeOutput("", tbl)
	}
}
