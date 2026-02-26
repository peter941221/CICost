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
	fmt.Println(`CICost - GitHub Actions 成本与浪费热区分析 CLI

Usage:
  cicost <command> [flags]

Commands:
  init       交互式初始化配置
  scan       拉取 GitHub Actions 数据并缓存到 SQLite
  report     生成综合成本报告（table/md/json/csv）
  reconcile  估算值与实际账单对账并生成校准系数
  policy     策略门禁（check/lint/explain）
  suggest    生成可执行优化建议（text/yaml + patch）
  org-report 多仓聚合报告（md/json，支持 partial result）
  hotspots   热区排行（workflow/job/runner/branch）
  budget     预算告警（stdout/webhook/file）
  explain    生成可执行优化建议
  config     show/edit 配置
  version    打印版本信息
  help       显示帮助信息`)
	return nil
}

var errNotImplemented = errors.New("not implemented yet; see 技术文档.MD for MVP scope")

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
