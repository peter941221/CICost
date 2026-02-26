# Changelog

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

