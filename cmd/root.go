package cmd

import (
	"errors"
	"fmt"
	"strings"
)

type commandHandler func(args []string) error

var commands = map[string]commandHandler{
	"init":     runInit,
	"scan":     runScan,
	"report":   runReport,
	"hotspots": runHotspots,
	"budget":   runBudget,
	"explain":  runExplain,
	"config":   runConfig,
	"version":  runVersion,
	"help":     runHelp,
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
  init       交互式初始化配置（待实现）
  scan       拉取数据并缓存（待实现）
  report     生成成本报告（待实现）
  hotspots   生成热区排行（待实现）
  budget     预算告警（待实现）
  explain    优化建议（待实现）
  config     查看/编辑配置（待实现）
  version    打印版本信息
  help       显示帮助信息`)
	return nil
}

var errNotImplemented = errors.New("not implemented yet; see 技术文档.MD for MVP scope")

