# Changelog

## v0.2.0 - 2026-02-26

### Added

- Pricing v2:
  - snapshot-based SKU pricing (`pricing_snapshots`, `effective_from`)
  - automatic historical snapshot selection
  - pricing metadata output (`pricing_snapshot_version`, `effective_from`, `pricing_source`)
- Reconciliation:
  - `cicost reconcile --repo --month --source --apply-calibration`
  - billing snapshot storage and reconcile result storage
  - confidence bands (`high/medium/low`)
  - `report --calibrated` support
- Policy Gate:
  - `cicost policy lint/check/explain`
  - policy-as-code expression parser and evaluator
  - `error` severity exits with code `3`
  - `policy_runs` audit table
- Suggestion Engine v2:
  - `cicost suggest --format text|yaml --output patches/`
  - executable patch snippets (concurrency/cache/paths/runner migration)
  - evidence-backed suggestions only
  - `suggestion_history` audit table
- Org aggregation:
  - `cicost org-report --repos file.txt --days --format`
  - multi-repo parallel aggregation
  - partial-result mode when single repo fails
- Schema v2:
  - `billing_snapshots`
  - `reconcile_results`
  - `policy_runs`
  - `suggestion_history`

### Quality

- Added tests for:
  - pricing snapshots and rate resolution
  - schema v2 tables
  - reconcile + calibrated report integration
  - policy lint/evaluate + exit code behavior
  - suggest generation and YAML output
  - org-report partial result behavior
- Passed:
  - `go test ./...`
  - `go test -race ./...`
  - `go vet ./...`

## v0.1.0 - 2026-02-26

### Added

- MVP CLI commands:
  - `init`, `scan`, `report`, `hotspots`, `budget`, `explain`, `config`, `version`
- GitHub Actions data ingestion:
  - workflow runs and jobs fetch
  - pagination support
  - token resolution chain (CLI/env/gh/config)
- Local persistence:
  - SQLite schema
  - incremental sync cursor
  - budget checks history
- Analytics:
  - cost estimation (OS multipliers + free tier)
  - waste analysis (rerun/cancel/fail rate)
  - hotspots ranking
  - budget projection and threshold status
  - recommendation generation
- Output formats:
  - table / markdown / json / csv
- Quality and delivery:
  - unit + integration tests
  - race/vet checks
  - CI workflow
  - release pipeline with GoReleaser
  - GitHub CLI extension entrypoint (`gh-cicost`)
  - Homebrew formula template

### Notes

- Cost output is an estimate and may differ from official billing.
- Free tier is shared at account/org level; per-repo calculation uses configured assumptions.
