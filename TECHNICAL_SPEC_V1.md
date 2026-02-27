# CICost Technical Specification v1 (MVP)

## Positioning

CICost translates GitHub Actions usage data into cost visibility, waste hotspots, and budget alerts in a CLI-first workflow.

## Scope

- Core commands: `init`, `scan`, `report`, `hotspots`, `budget`, `explain`, `config`, `version`
- Data ingestion: workflow runs + jobs via GitHub API
- Local persistence: SQLite cache with incremental sync cursor
- Output: table, markdown, JSON, CSV

## Architecture

```text
CLI
 |
 +--> Config/Auth
 +--> GitHub API Client
 +--> SQLite Store
 +--> Analytics (cost/waste/hotspots/budget)
 +--> Output Renderer
```

## Cost Model (v1)

- Per-job duration rounded to billable minutes.
- Hosted runner billing only (`self-hosted` excluded).
- Free-tier assumptions applied per config.
- Output marked as estimate.

## MVP Acceptance

- `scan` can pull and persist run/job data.
- `report` provides cost and waste metrics for a selected period.
- `budget` returns exit code `2` on exceeded/projected exceed cases.
- `go test ./...`, `go test -race ./...`, and `go vet ./...` pass.
