# MEMORY

## 2026-02-26

### Progress

- Initialized this directory as an independent Git repository (`main`).
- Bound remote origin to `https://github.com/peter941221/CICost.git`.
- Parsed `技术文档.MD` and aligned initial scaffold with the documented structure.
- Added first-pass Go CLI skeleton and internal package placeholders.
- Added baseline project files: `.gitignore`, `Makefile`, `README.md`, `.cicost.yml.example`, `configs/pricing_default.yml`.
- Created initial commit: `22a1f73 chore: bootstrap cicost repository scaffold`.
- Pushed branch `main` to remote `origin/main`.
- Installed Go toolchain via Scoop (`go1.26.0`) and removed local environment blocker.
- Implemented MVP command chain:
  - `init/config` (config bootstrap + merged config show/edit)
  - `scan` (GitHub API runs/jobs fetch + pagination + worker pool + SQLite upsert + sync cursor)
  - `report` (table/md/json/csv output + compare)
  - `hotspots` (workflow/job/runner/branch ranking)
  - `budget` (weekly/monthly threshold, projection, webhook/file/stdout, exit code 2)
  - `explain` (rule-based optimization suggestions)
- Implemented core modules:
  - `internal/config` merged configuration hierarchy
  - `internal/auth` token chain (CLI/env/gh/config)
  - `internal/github` REST client + pagination parsing
  - `internal/store` SQLite schema + cursor + query interfaces
  - `internal/analytics` cost/waste/hotspots/budget/trend
  - `internal/output` table/markdown/json/csv rendering
- Added CI workflow `.github/workflows/ci.yml` with `go vet`, `go test -race`, and cross-platform builds.
- Added unit tests for pricing, budget, waste, pagination, and sync cursor.
- Added command-level integration test (`cmd/integration_test.go`) for report + budget path.
- Added `CICOST_GITHUB_API_BASE_URL` override to support mock server / enterprise API endpoints.
- Added release/distribution pipeline:
  - `.goreleaser.yml` for dual binaries (`cicost`, `gh-cicost`) + checksums.
  - `.github/workflows/release.yml` for tag-triggered release publishing.
  - `cmd/gh-cicost/main.go` as GitHub CLI extension entrypoint.
  - `Formula/cicost.rb` template for Homebrew tap publishing.
  - `docs/RELEASE.md` with tagging and distribution runbook.
- Installed local `goreleaser` and validated release config via `goreleaser check`.
- Prepared release candidate `v0.1.0`:
  - set default CLI version to `0.1.0`
  - added `CHANGELOG.md` first release notes
  - linked changelog in README distribution section
- Released `v0.1.0`:
  - pushed tag `v0.1.0`
  - release workflow completed successfully
  - GitHub Release published with cicost + gh-cicost artifacts and checksums

## 2026-02-26 (Post-v0.1.0 Product Strategy)

### Competitive Research Summary

- Benchmarked key alternatives and adjacent products:
  - GitHub native billing/runner pricing docs (source-of-truth for rates and quotas)
  - Datadog CI Visibility (pipeline observability + committer-based pricing)
  - BuildPulse (flaky tests + CI speed/cost claims)
  - BuildJet / Depot / RunsOn (runner-layer cost optimization vendors)
  - Infracost (FinOps-in-PR for IaC, policy and guardrail style)

### Product Gap Highlights

- Current CICost pricing defaults are legacy (`linux_per_min: 0.008` with multipliers), while current GitHub docs show SKU-based rates (e.g., Linux 2-core `0.006`, Windows 2-core `0.010`, macOS `0.062`).
- CICost currently estimates from local workflow/job durations only; lacks org-level calibration against real billing exports.
- Suggestion engine is useful but still shallow compared to “policy + guardrail + automation” products.

### Strengthening Priorities

1. Pricing accuracy v2:
   - switch from multiplier model to SKU-based rate table
   - support dated pricing snapshots and explicit effective dates
2. Confidence & calibration:
   - add optional reconciliation mode against GitHub billing endpoints/exports
   - surface confidence bands in report
3. Guardrails:
   - add policy engine (budget by repo/team/workflow, fail PR on threshold)
4. Monetizable differentiation:
   - “waste prevention mode” with actionable YAML patch suggestions and PR comments
- Completed live end-to-end smoke run against public repo `cli/cli`:
  - `scan`, `report`, `hotspots`, `explain`, `budget` all executable.

### Decisions

- Chose independent nested repo strategy to avoid contaminating parent monorepo changes.
- Chose `modernc.org/sqlite` (pure Go) to avoid CGO friction on Windows.
- Kept CLI parser on Go stdlib `flag` for lightweight MVP delivery.

### Next Actions

1. Add recorded fixture-based integration tests for scan/report deterministic CI.
2. Improve `scan` incremental window and data completeness accounting.
3. Add GoReleaser packaging and gh extension adapter.
