# CICost

GitHub Actions 成本与浪费热区分析 CLI（MVP）。

## MVP Status (2026-02-26)

- [x] 独立仓库初始化
- [x] 核心命令可执行（init/scan/report/hotspots/budget/explain/config/version）
- [x] SQLite 缓存与增量游标
- [x] 成本/浪费/热区/预算分析基础实现
- [x] 单元测试与 CI 工作流
- [ ] GitHub API 录制 fixture 的端到端集成测试（下一阶段）

## Architecture

```text
User CLI
   |
   v
Command Router (cmd/*)
   |
   +--> Config/Auth
   +--> GitHub API Client
   +--> SQLite Store
   +--> Analytics Engine
   +--> Output Formatter
```

## Quick Start

1. 安装 Go 1.26+
2. 认证（任选其一）:
   - `gh auth login`
   - `set GITHUB_TOKEN=ghp_xxx` (Windows)
3. 拉取并生成报告:

```bash
go run . scan --repo owner/repo --days 30
go run . report --repo owner/repo --format table
go run . hotspots --repo owner/repo --group-by workflow --top 5
go run . budget --repo owner/repo --monthly 100
go run . explain --repo owner/repo
```

## Distribution

- Standalone binary: GitHub Releases artifacts (`cicost_*`).
- GitHub CLI extension: `gh extension install peter941221/CICost` then `gh cicost ...`.
- Homebrew: formula template at `Formula/cicost.rb`.
- Release pipeline: `.github/workflows/release.yml` + `.goreleaser.yml`.
- Release runbook: `docs/RELEASE.md`.

## Config

- 用户级: `~/.cicost/config.yml`
- 仓库级: `.cicost.yml`
- 示例: `.cicost.yml.example`

## Project Layout

```text
CICost/
├── cmd/                # CLI command router
├── internal/
│   ├── analytics/      # cost/waste/hotspot/budget logic
│   ├── auth/           # token resolution
│   ├── github/         # API client and pagination
│   ├── model/          # shared domain types
│   ├── output/         # table/md/json/csv formatters
│   ├── pricing/        # pricing and free tier logic
│   └── store/          # SQLite schema and sync cursor
├── configs/            # pricing defaults
├── testdata/           # fixtures for integration tests
├── 技术文档.MD          # product + technical spec v1.0
├── MEMORY.md           # project memory
└── RUNBOOK.md          # execution handbook
```
