# CICost

GitHub Actions cost analysis and governance CLI.

## Status (2026-02-27)

- [x] Pricing v2: SKU-based pricing + `effective_from` snapshots
- [x] Reconcile: estimate vs actual billing with calibration factor
- [x] Policy Gate: `policy lint/check/explain` + `error => exit code 3`
- [x] Suggest: actionable optimization output (`text|yaml` + patch export)
- [x] Org Report: multi-repo aggregation with partial-result support
- [x] Quality gate: `go test ./...` / `go test -race ./...` / `go vet ./...`

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

1. Install Go 1.24+ (CI validated on Go 1.26.x)
2. Authenticate (one option):
   - `gh auth login`
   - `set GITHUB_TOKEN=ghp_xxx` (Windows)
3. Run the core flow:

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

- User-level: `~/.cicost/config.yml`
- Repo-level: `.cicost.yml`
- Policy file: `.cicost.policy.yml` (example: `.cicost.policy.yml.example`)
- Pricing file: `configs/pricing_default.yml` (supports `pricing_snapshots`)

## Distribution

- Standalone binary: GitHub Releases artifacts (`cicost_*`)
- GitHub CLI extension: `gh extension install peter941221/gh-cicost`
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
├── TECHNICAL_SPEC_V1.md
├── TECHNICAL_SPEC_V2.md
├── MEMORY.md
└── RUNBOOK.md
```
