# CICost

GitHub Actions cost, waste, and policy governance CLI for engineers and FinOps teams.

[![CI](https://github.com/peter941221/CICost/actions/workflows/ci.yml/badge.svg)](https://github.com/peter941221/CICost/actions/workflows/ci.yml)
[![Release](https://github.com/peter941221/CICost/actions/workflows/release.yml/badge.svg)](https://github.com/peter941221/CICost/actions/workflows/release.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24%2B-00ADD8?logo=go)](go.mod)

## CLI Demo (GIF)

![CICost CLI Demo](https://github.com/peter941221/CICost/releases/download/v0.2.0/cicost-cli-demo-v4.gif)

This GIF shows a no-network, no-error walkthrough:
- `cicost version`
- `cicost help`
- `cicost policy lint --policy .cicost.policy.yml.example`
- `cicost policy explain`
- `cicost config show`

## Why CICost

- Turn GitHub Actions usage into cost in USD using pricing snapshots.
- Reconcile estimate vs actual billing and persist calibration factor.
- Enforce cost guardrails with policy checks in CI.
- Generate data-backed optimization suggestions (with patch snippets).
- Aggregate multi-repo visibility for org-level decisions.

## Current Capability (v0.2.0)

- [x] Pricing v2 (`pricing_snapshots`, `effective_from`, legacy fallback)
- [x] Reconcile (`--actual-usd`, CSV import, optional calibration apply)
- [x] Policy Gate (`policy lint/check/explain`, error rule => exit code `3`)
- [x] Suggestion Engine (`text|yaml`, patch artifact export)
- [x] Org Report (parallel multi-repo aggregation, partial-failure support)
- [x] Quality gates (`go test ./...`, `go test -race ./...`, `go vet ./...`)

## How It Works

```text
GitHub Actions API / CSV Billing
            |
            v
      [cicost scan]
            |
            v
   SQLite Local Store (runs/jobs + v2 tables)
            |
            +-------------------+
            |                   |
            v                   v
 [report/hotspots/budget]   [reconcile/policy/suggest/org-report]
            |                   |
            +---------+---------+
                      |
                      v
          table / md / json / csv / yaml
```

## Quick Start

1. Prerequisites:
   - Go `1.24+` (CI validated on `1.26.x`)
   - GitHub auth via `gh auth login` or `GITHUB_TOKEN`
2. Build and verify:

```bash
go build -o cicost .
./cicost version
```

3. Run the core workflow:

```bash
./cicost scan --repo owner/repo --days 30
./cicost report --repo owner/repo --format table
./cicost reconcile --repo owner/repo --month 2026-02 --actual-usd 123.45 --apply-calibration
./cicost report --repo owner/repo --calibrated --format json
./cicost policy lint --policy .cicost.policy.yml
./cicost policy check --repo owner/repo --days 30
./cicost suggest --repo owner/repo --format yaml --output patches/
./cicost org-report --repos repos.txt --days 30 --format md
```

## Command Map

| Command | Purpose | Key Flags |
|---|---|---|
| `scan` | Fetch workflow runs/jobs into local SQLite cache | `--repo --days --incremental --full --workers` |
| `report` | Cost + waste report with optional calibration | `--repo --days --format --compare --calibrated` |
| `hotspots` | Rank expensive workflows/jobs/runners/branches | `--group-by --top --sort --format` |
| `budget` | Weekly/monthly budget check and alerting | `--monthly --weekly --notify --webhook-url` |
| `reconcile` | Estimate vs actual billing reconciliation | `--month --source --input --actual-usd --apply-calibration` |
| `policy` | Governance checks (`lint/check/explain`) | `policy check --repo --days --policy` |
| `suggest` | Data-backed optimization suggestions | `--repo --days --format --output` |
| `org-report` | Multi-repo aggregate report | `--repos --days --format --output` |

## Config and Files

- User config: `~/.cicost/config.yml`
- Repo config: `.cicost.yml`
- Policy rules: `.cicost.policy.yml` (sample: `.cicost.policy.yml.example`)
- Pricing defaults: `configs/pricing_default.yml`
- Local data store: platform-specific path resolved by `internal/config`

## Exit Codes (for CI Integration)

- `0`: success
- `1`: generic error
- `2`: budget warning/exceeded
- `3`: policy error rules matched

## Install and Distribution

- Standalone binaries: GitHub Releases (`cicost_*`, checksums included)
- GitHub CLI extension: `gh extension install peter941221/gh-cicost`
- Homebrew formula template: `Formula/cicost.rb`
- Release automation: `.github/workflows/release.yml` + `.goreleaser.yml`
- Extension sync automation: `.github/workflows/sync-gh-extension.yml`

## Documentation

- Operations runbook: [RUNBOOK.md](RUNBOOK.md)
- Release process: [docs/RELEASE.md](docs/RELEASE.md)
- Launch kit: [docs/LAUNCH.md](docs/LAUNCH.md)
- Demo recording script (headless): `docs/scripts/record_cli_gif_headless.ps1`
- Technical specs: [TECHNICAL_SPEC_V1.md](TECHNICAL_SPEC_V1.md), [TECHNICAL_SPEC_V2.md](TECHNICAL_SPEC_V2.md)

## Community and Trust

- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
- Code of conduct: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- Security policy: [SECURITY.md](SECURITY.md)
- License: Apache-2.0 ([LICENSE](LICENSE))

## Project Layout

```text
CICost/
├── cmd/                     # CLI command entrypoints
├── internal/
│   ├── analytics/           # cost/waste/hotspot/budget logic
│   ├── billing/             # billing import adapters
│   ├── policy/              # policy parser/evaluator
│   ├── pricing/             # snapshot loader + resolver
│   ├── reconcile/           # estimate-vs-actual calibration
│   ├── suggest/             # recommendation generation
│   └── store/               # SQLite schema + access layer
├── configs/                 # default pricing/config examples
├── docs/                    # release + launch docs
├── MEMORY.md                # project memory log
└── RUNBOOK.md               # operational commands
```
