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
	repoFlag := fs.String("repo", "", "Target repository in owner/repo format")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "Time window in days")
	groupByFlag := fs.String("group-by", "workflow", "Group by: workflow|job|runner|branch")
	topFlag := fs.Int("top", 10, "Show top N entries")
	sortFlag := fs.String("sort", "cost", "Sort by: cost|minutes|fail_rate")
	formatFlag := fs.String("format", "table", "Output format: table|md|json")
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
