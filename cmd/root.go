package cmd

import (
	"errors"
	"fmt"
	"strings"
)

type commandHandler func(args []string) error

var commands = map[string]commandHandler{
	"init":       runInit,
	"scan":       runScan,
	"report":     runReport,
	"reconcile":  runReconcile,
	"policy":     runPolicy,
	"hotspots":   runHotspots,
	"budget":     runBudget,
	"suggest":    runSuggest,
	"org-report": runOrgReport,
	"explain":    runExplain,
	"config":     runConfig,
	"version":    runVersion,
	"help":       runHelp,
}

// Execute routes subcommands and returns a structured error for main() to print.
func Execute(args []string) error {
	if len(args) == 0 {
		return runHelp(nil)
	}

	name := strings.ToLower(args[0])
	handler, ok := commands[name]
	if !ok {
		_ = runHelp(nil)
		return fmt.Errorf("unknown command %q", name)
	}
	if err := handler(args[1:]); err != nil {
		return err
	}
	return nil
}

func runHelp(_ []string) error {
	fmt.Println(`CICost - GitHub Actions cost and waste hotspot analysis CLI

Usage:
  cicost <command> [flags]

Commands:
  init       Interactive config initialization
  scan       Fetch GitHub Actions data and cache to SQLite
  report     Generate cost report (table/md/json/csv)
  reconcile  Reconcile estimate vs actual billing and persist calibration
  policy     Policy gate (check/lint/explain)
  suggest    Generate actionable optimization suggestions (text/yaml + patch)
  org-report Multi-repo aggregate report (md/json, partial result supported)
  hotspots   Hotspot ranking (workflow/job/runner/branch)
  budget     Budget alerting (stdout/webhook/file)
  explain    Generate optimization suggestions
  config     show/edit config
  version    Print version info
  help       Show help`)
	return nil
}

var errNotImplemented = errors.New("not implemented yet; see TECHNICAL_SPEC_V1.md for MVP scope")

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.Code)
	}
	return e.Err.Error()
}

func (e ExitError) Unwrap() error { return e.Err }

func withExit(code int, err error) error {
	return ExitError{Code: code, Err: err}
}

func ExitCode(err error) int {
	var ex ExitError
	if errors.As(err, &ex) && ex.Code > 0 {
		return ex.Code
	}
	return 1
}
