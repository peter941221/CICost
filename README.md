# CICost

GitHub Actions 成本与浪费热区分析 CLI（MVP）。

## Current Status

- Date: 2026-02-26
- Stage: Repo initialized + scaffold ready
- Next: Implement `scan -> store -> report` happy path

## Quick Start

1. Install Go 1.23+
2. Build:

```bash
make build
```

3. Run:

```bash
./bin/cicost help
```

## Command Roadmap

- `cicost init`
- `cicost scan --repo owner/repo --days 30`
- `cicost report --repo owner/repo --format table`
- `cicost hotspots --group-by workflow --top 10`
- `cicost budget --monthly 100`
- `cicost explain --repo owner/repo`

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

