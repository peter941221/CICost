# CICost

GitHub Actions cost intelligence CLI for engineering and FinOps teams.

[![CI](https://github.com/peter941221/CICost/actions/workflows/ci.yml/badge.svg)](https://github.com/peter941221/CICost/actions/workflows/ci.yml)
[![Release](https://github.com/peter941221/CICost/actions/workflows/release.yml/badge.svg)](https://github.com/peter941221/CICost/actions/workflows/release.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24%2B-00ADD8?logo=go)](go.mod)

> Turn CI activity into cost visibility, hotspots, and policy decisions before billing surprises happen.

## About

CICost helps you answer three questions fast:

- Where is GitHub Actions spend going?
- Which workflows are wasting minutes?
- Are we still inside budget policy guardrails?

### GitHub Repository About (Copy/Paste)

- Description: `GitHub Actions cost intelligence CLI: scan usage, estimate spend, find hotspots, and enforce CI budget policies.`
- Topics: `github-actions`, `finops`, `cost-optimization`, `ci-cd`, `devops`, `golang`, `cli`

## CLI Demo

![CICost CLI Demo](https://github.com/peter941221/CICost/releases/download/v0.2.0/cicost-cli-demo-v8.gif)

Narrated in English subtitles:

- scan recent workflow data
- review cost summary
- find top expensive workflows
- run policy checks
- decide next action

## Why Teams Use CICost

| Need | CICost gives you | Outcome |
|---|---|---|
| CI cost visibility | Usage-to-USD estimation with pricing snapshots | No more blind spots |
| Waste discovery | Hotspots and waste metrics | Faster optimization loops |
| Budget control | Policy checks in CI (`exit code 3` on error rules) | Enforceable governance |
| Multi-repo oversight | Org-level aggregation (`org-report`) | Better portfolio decisions |

## 30-Second Flow

```text
[scan GitHub Actions data]
          |
          v
[local SQLite cache]
          |
          v
[report + hotspots]
          |
          v
[policy check]
          |
          v
[action: optimize or enforce]
```

## Quick Start

1. Prerequisites

- Go `1.24+` (CI validated on `1.26.x`)
- GitHub auth via `gh auth login` or `GITHUB_TOKEN`

2. Build once

```bash
go build -o cicost .
./cicost version
```

3. Run core workflow

```bash
./cicost scan --repo owner/repo --days 30
./cicost report --repo owner/repo --format table
./cicost hotspots --repo owner/repo --days 30 --group-by workflow --top 5 --sort cost
./cicost policy check --repo owner/repo --days 30 --policy .cicost.policy.yml
```

4. Optional advanced flow

```bash
./cicost reconcile --repo owner/repo --month 2026-02 --actual-usd 123.45 --apply-calibration
./cicost report --repo owner/repo --calibrated --format json
./cicost suggest --repo owner/repo --format yaml --output patches/
./cicost org-report --repos repos.txt --days 30 --format md
```

## Command Map

| Command | What it does | Common flags |
|---|---|---|
| `scan` | Pull runs/jobs into local cache | `--repo --days --incremental --full --workers` |
| `report` | Cost and waste report | `--repo --days --format --compare --calibrated` |
| `hotspots` | Rank costly workflows/jobs/runners/branches | `--group-by --top --sort --format` |
| `budget` | Budget check and notifications | `--monthly --weekly --notify --webhook-url` |
| `reconcile` | Estimate vs actual calibration | `--month --source --input --actual-usd --apply-calibration` |
| `policy` | Lint/check/explain budget policies | `policy check --repo --days --policy` |
| `suggest` | Data-backed optimization suggestions | `--repo --days --format --output` |
| `org-report` | Multi-repo summary | `--repos --days --format --output` |

## Output and CI Exit Codes

- `0`: success
- `1`: generic error
- `2`: budget warning/exceeded
- `3`: policy error rule matched

## Current Capability (v0.2.0)

- [x] Pricing v2 (`pricing_snapshots`, `effective_from`, legacy fallback)
- [x] Reconcile (`--actual-usd`, CSV import, optional calibration apply)
- [x] Policy Gate (`policy lint/check/explain`, error rule => exit code `3`)
- [x] Suggestion Engine (`text|yaml`, patch artifact export)
- [x] Org Report (parallel multi-repo aggregation, partial-failure support)
- [x] Quality gates (`go test ./...`, `go test -race ./...`, `go vet ./...`)

## How It Works

```text
GitHub Actions API / Billing Input
               |
               v
      +----------------------+
      |      cicost scan     |
      +----------+-----------+
                 |
                 v
      +----------------------+
      |   SQLite local data  |
      | runs/jobs/reconcile  |
      +-----+-----------+----+
            |           |
            v           v
   +--------+---+   +---+-------------------+
   | report/hotspots | policy/suggest/org   |
   +--------+---+   +---+-------------------+
            \           /
             \         /
              v       v
         table / md / json / csv / yaml
```

## Configuration

- User config: `~/.cicost/config.yml`
- Repo config: `.cicost.yml`
- Policy rules: `.cicost.policy.yml` (sample: `.cicost.policy.yml.example`)
- Pricing defaults: `configs/pricing_default.yml`
- Local DB path: resolved by `internal/config`

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
- Demo recorder (headless): `docs/scripts/record_cli_gif_headless.ps1`
- Subtitle renderer: `docs/scripts/render_subtitled_gif_v7.ps1`
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
