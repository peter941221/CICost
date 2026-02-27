package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/peter941221/CICost/internal/analytics"
	"github.com/peter941221/CICost/internal/config"
	"github.com/peter941221/CICost/internal/model"
	"github.com/peter941221/CICost/internal/store"
	"github.com/peter941221/CICost/internal/suggest"
)

func runSuggest(args []string) error {
	rt, err := newRuntimeContext()
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("suggest", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "Target repository in owner/repo format")
	daysFlag := fs.Int("days", rt.cfg.Scan.Days, "Time window in days")
	formatFlag := fs.String("format", "text", "Output format: text|yaml")
	outputFlag := fs.String("output", "", "Output path (directory or file)")
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
		fmt.Println("No suggestions generated: insufficient local data.")
		return nil
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
	hotspots := analytics.CalculateHotspots(runs, jobs, pcfg, analytics.HotspotOptions{
		GroupBy: "workflow",
		TopN:    5,
		SortBy:  "cost",
	})

	suggestions := suggest.Generate(suggest.Inputs{
		Repo:     repo,
		Runs:     runs,
		Jobs:     jobs,
		Cost:     cost,
		Waste:    waste,
		Hotspots: hotspots,
	})
	if len(suggestions) == 0 {
		fmt.Println("No data-backed suggestions found.")
		return nil
	}

	for _, s := range suggestions {
		evidence, _ := json.Marshal(s.Evidence)
		if err := st.InsertSuggestionHistory(model.SuggestionRecord{
			Repo:               repo,
			PeriodStart:        start,
			PeriodEnd:          end,
			SuggestionType:     s.Type,
			Title:              s.Title,
			EstimatedSavingUSD: s.EstimatedSavingUSD,
			EvidenceJSON:       string(evidence),
		}); err != nil {
			return err
		}
	}

	format := strings.ToLower(strings.TrimSpace(*formatFlag))
	switch format {
	case "yaml":
		b, err := yaml.Marshal(suggestions)
		if err != nil {
			return err
		}
		if *outputFlag != "" {
			if err := writeSuggestArtifacts(*outputFlag, suggestions, string(b)); err != nil {
				return err
			}
		} else {
			fmt.Print(string(b))
		}
	default:
		var lines []string
		lines = append(lines, fmt.Sprintf("CICost Suggestions for %s", repo))
		for i, s := range suggestions {
			lines = append(lines, fmt.Sprintf("%d. [%s] %s", i+1, strings.ToUpper(s.Type), s.Title))
			lines = append(lines, fmt.Sprintf("   Problem: %s", s.Problem))
			lines = append(lines, fmt.Sprintf("   Current data: %s", s.CurrentData))
			lines = append(lines, fmt.Sprintf("   Estimated saving: $%.2f", s.EstimatedSavingUSD))
			lines = append(lines, "   Patch snippet:")
			for _, ln := range strings.Split(s.Patch, "\n") {
				lines = append(lines, "     "+ln)
			}
		}
		out := strings.Join(lines, "\n") + "\n"
		if err := writeOutput(*outputFlag, out); err != nil {
			return err
		}
	}

	return nil
}

func writeSuggestArtifacts(path string, suggestions []suggest.Suggestion, yamlContent string) error {
	target := strings.TrimSpace(path)
	isFile := strings.HasSuffix(strings.ToLower(target), ".yml") || strings.HasSuffix(strings.ToLower(target), ".yaml")
	if isFile {
		return writeOutput(target, yamlContent)
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}
	summaryPath := filepath.Join(target, "suggestions.yaml")
	if err := writeOutput(summaryPath, yamlContent); err != nil {
		return err
	}
	for i, s := range suggestions {
		name := fmt.Sprintf("%02d_%s.patch.yml", i+1, s.Type)
		if err := writeOutput(filepath.Join(target, name), s.Patch+"\n"); err != nil {
			return err
		}
	}
	return nil
}
