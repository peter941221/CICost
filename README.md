# CICost

GitHub Actions 成本分析与治理 CLI。

## Status (2026-02-26)

- [x] Pricing v2：SKU 定价 + `effective_from` 版本快照
- [x] Reconcile：估算值与实际账单对账 + 校准系数
- [x] Policy Gate：`policy lint/check/explain` + `error => exit code 3`
- [x] Suggest：可执行优化建议（`text|yaml` + patch 文件导出）
- [x] Org Report：多仓聚合，支持 partial result
- [x] 测试门禁：`go test ./...` / `go test -race ./...` / `go vet ./...`

## Architecture

```text
User / CI
   |
   v
Commands: scan/report/reconcile/policy/suggest/org-report
   |
   +--> Config/Auth
   +--> GitHub Data Ingestion
   +--> SQLite Store (runs/jobs/billing/reconcile/policy/suggest)
   +--> Analytics + Policy Engine
   +--> Output (table/md/json/csv/yaml)
```

## Quick Start

1. 安装 Go 1.24+（CI 使用 Go 1.26.x 验证）
2. 认证（任选其一）:
   - `gh auth login`
   - `set GITHUB_TOKEN=ghp_xxx` (Windows)
3. 执行核心流程:

```bash
go run . scan --repo owner/repo --days 30
go run . report --repo owner/repo --format table
go run . reconcile --repo owner/repo --month 2026-02 --actual-usd 123.45 --apply-calibration
go run . report --repo owner/repo --calibrated --format json
go run . policy lint --policy .cicost.policy.yml
go run . policy check --repo owner/repo --days 30
go run . suggest --repo owner/repo --format yaml --output patches/
go run . org-report --repos repos.txt --days 30 --format md
```

## Config

- 用户级: `~/.cicost/config.yml`
- 仓库级: `.cicost.yml`
- 策略文件: `.cicost.policy.yml`（示例：`.cicost.policy.yml.example`）
- 定价文件: `configs/pricing_default.yml`（支持 `pricing_snapshots`）

## Distribution

- Standalone binary: GitHub Releases artifacts (`cicost_*`)
- GitHub CLI extension: `gh extension install peter941221/CICost`
- Homebrew: `Formula/cicost.rb`
- Release pipeline: `.github/workflows/release.yml` + `.goreleaser.yml`
- Release runbook: `docs/RELEASE.md`
- Launch kit: `docs/LAUNCH.md`

## License

Apache-2.0. See `LICENSE`.

## Community

- Contribution guide: `CONTRIBUTING.md`
- Code of conduct: `CODE_OF_CONDUCT.md`
- Security policy: `SECURITY.md`

## Project Layout

```text
CICost/
├── cmd/                     # CLI commands
├── internal/
│   ├── analytics/           # cost/waste/hotspot/budget
│   ├── billing/             # billing import adapters
│   ├── policy/              # policy parser + evaluator
│   ├── pricing/             # pricing snapshots + resolver
│   ├── reconcile/           # estimate-vs-actual calibration
│   ├── suggest/             # actionable patch suggestions
│   └── store/               # SQLite schema v2 + access layer
├── configs/
├── 技术文档.MD
├── 技术文档2.MD
├── MEMORY.md
└── RUNBOOK.md
```
