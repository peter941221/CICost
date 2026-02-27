# MEMORY

## 2026-02-26

### Progress

- Initialized this directory as an independent Git repository (`main`).
- Bound remote origin to `https://github.com/peter941221/CICost.git`.
- Parsed `TECHNICAL_SPEC_V1.md` and aligned initial scaffold with the documented structure.
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

## 2026-02-26 (Planning Enhancement v2)

### Progress

- Added `TECHNICAL_SPEC_V2.md` as post-v0.1.0 strengthening blueprint.
- Document includes:
  - competitive positioning summary
  - v2 architecture upgrade
  - feature roadmap (pricing v2 / reconcile / policy / suggestion / org-report)
  - detailed acceptance criteria (AC-PRICING/REC/POL/SUG/ORG)
  - test gate and risk-tier matrix
  - release/rollback strategy and done definition

### Validation

- Verified file creation and structure:
  - file exists
  - 300+ lines content
  - AC entries detected and indexed

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

## 2026-02-26 (v2 Implementation Sprint)

### Progress

- Implemented Pricing v2 engine:
  - added `pricing_snapshots` + `effective_from` loader support
  - added snapshot resolver and SKU direct rate pricing
  - kept legacy multiplier fallback with warning
- Upgraded schema to v2 with new tables:
  - `billing_snapshots`
  - `reconcile_results`
  - `policy_runs`
  - `suggestion_history`
- Added `reconcile` command:
  - monthly estimate vs actual reconciliation
  - confidence grading and calibration factor persistence
  - optional `--apply-calibration` for future `report --calibrated`
- Added policy engine and command group:
  - `policy lint/check/explain`
  - expression parser for threshold rules
  - CI gate behavior (`error => exit code 3`)
- Added suggestion engine v2 and `suggest` command:
  - text/yaml output
  - executable patch file export
  - evidence-backed recommendations only
- Added `org-report` command:
  - multi-repo parallel aggregation
  - partial result support when some repos have no data/fail
- Updated docs and examples:
  - README, RUNBOOK, CHANGELOG
  - added `.cicost.policy.yml.example`

### Validation

- `go test ./...` ✅
- `go test -race ./...` ✅
- `go vet ./...` ✅

## 2026-02-27 (Go-to-Market Readiness Assessment)

### Progress

- Performed repository-wide readiness review for promotion/go-to-market.
- Reviewed product, runbook, release, and changelog docs:
  - `README.md`
  - `RUNBOOK.md`
  - `docs/RELEASE.md`
  - `CHANGELOG.md`
- Reviewed CI/CD setup:
  - `.github/workflows/ci.yml`
  - `.github/workflows/release.yml`
  - `.goreleaser.yml`
- Verified repository state and release/tag signal:
  - working tree clean
  - latest tag remains `v0.1.0` while changelog includes `v0.2.0` section
- Identified documentation consistency gaps:
  - Go version requirements are inconsistent across files (`README` 1.26+, `RUNBOOK` 1.23+, `go.mod` 1.24.0)
  - no LICENSE file found at repository root

### Validation

- `go test ./...` ✅
- `go test -race ./...` ✅
- `go vet ./...` ✅
- `go test -cover ./...` ✅ (package-level coverage collected)
- `go run . help` ✅
- `goreleaser check` ✅

## 2026-02-27 (License + GTM 48h Plan Start)

### Progress

- Confirmed license direction: Apache-2.0.
- Added repository root `LICENSE` file using Apache License 2.0 text.
- Updated README to include explicit license declaration.
- Prepared 48-hour promotion hardening plan request output with PR slices and validation commands.

### Validation

- `rg -n "Apache License|Version 2.0" LICENSE` ✅
- `rg -n "^## License|Apache-2.0" README.md` ✅
- `go test ./...` ✅

## 2026-02-27 (Executed All GTM PR Scopes)

### Progress

- Completed PR1 scope (consistency + legal baseline):
  - set CLI default version to `0.2.0`
  - aligned Go version docs to `Go 1.24+` (CI remains `1.26.x`)
  - updated release doc tag example to `v0.2.0`
  - added Apache-2.0 `LICENSE`
- Completed PR2 scope (trust surface docs):
  - added `SECURITY.md`
  - added `CONTRIBUTING.md`
  - added `CODE_OF_CONDUCT.md`
  - linked trust docs from `README.md`
- Completed PR3 scope (quality gate hardening):
  - added `internal/github` tests for client headers, API error parsing, pagination/list mapping, runner inference
  - added `internal/store` tests for cursor missing/update and v2 persistence behaviors
  - added dedicated CI `coverage` job in `.github/workflows/ci.yml` with artifact upload
- Completed PR4 scope (release + launch assets):
  - added `docs/LAUNCH.md` (positioning/demo/checklist/launch template)
  - linked launch kit from `README.md`

### Validation

- `go test ./...` ✅
- `go test -race ./...` ✅
- `go vet ./...` ✅
- `go test -cover ./...` ✅
  - `internal/github`: `79.5%` (previously `8.3%`)
  - `internal/store`: `32.5%` (previously `13.8%`)
- `goreleaser check` ✅
- `go run . version` ✅ (`cicost 0.2.0`)

### Release Prep

- Created four PR-style commits for GTM execution scopes.
- Created local release tag `v0.2.0` on commit `3269de2`.
- Re-ran `go test ./...` after commits/tag: ✅
- Pushed `main` and `v0.2.0` to remote `origin`.

## 2026-02-27 (Install Matrix Smoke Validation)

### Progress

- Ran distribution install smoke checks after `v0.2.0` release publish.
- Verified standalone Windows binary path end-to-end:
  - download release zip
  - verify checksum from `checksums.txt`
  - extract and execute `cicost.exe version`
- Verified GitHub CLI extension install path:
  - `gh extension install peter941221/CICost` currently fails because extension repo naming must start with `gh-`.
  - URL-based install fallback also fails with the same naming rule.
- Verified Homebrew path readiness signal:
  - Homebrew unavailable in this Windows environment (`brew` command missing).
  - `Formula/cicost.rb` still contains template placeholders and outdated `version "0.1.0"`.

### Validation

- `gh --version` ✅ (`2.86.0`)
- `gh auth status` ✅ (authenticated as `peter941221`)
- standalone binary checksum compare ✅
- standalone binary run ✅ (`cicost 0.2.0`)
- `gh extension install peter941221/CICost` ❌ (repo name rule)
- `gh extension install https://github.com/peter941221/CICost` ❌ (repo name rule)

## 2026-02-27 (Extension Repo + English-Only Sweep)

### Progress

- Created GitHub extension repository: `https://github.com/peter941221/gh-cicost`.
- Published `v0.2.0` release in `gh-cicost` with extension binaries.
- Added additional release assets using `gh`-recognized naming (`gh-cicost_v0.2.0_<os>-<arch>[.exe]`) to make installation resolvable.
- Verified `gh extension install peter941221/gh-cicost` works and `gh cicost version` executes successfully.
- Completed English-only sweep in CICost:
  - translated all remaining Chinese CLI flag/help text to English
  - replaced README with full English version
  - replaced non-English spec files with `TECHNICAL_SPEC_V1.md` and `TECHNICAL_SPEC_V2.md`
  - removed non-English filenames from repository
  - updated docs to use `gh extension install peter941221/gh-cicost`

### Validation

- `go test ./...` ✅
- `go test -race ./...` ✅
- `go vet ./...` ✅
- `go run . help` ✅ (English output verified)
- `rg -nP "[\\p{Han}]" -S .` ✅ (no matches)
- `gh extension install peter941221/gh-cicost` ✅
- `gh cicost version` ✅

## 2026-02-27 (Homebrew + Extension Sync Automation)

### Progress

- Updated `Formula/cicost.rb` to release-ready state for `v0.2.0`:
  - set formula `version` to `0.2.0`
  - added `license "Apache-2.0"`
  - replaced all placeholder `sha256` values with real checksums for macOS/Linux artifacts
- Added automated extension mirroring workflow:
  - new workflow `.github/workflows/sync-gh-extension.yml`
  - trigger on `release: published` and manual `workflow_dispatch`
  - downloads `gh-cicost_*` artifacts from `peter941221/CICost` release
  - generates `gh`-recognized binary names (for direct extension resolution)
  - upserts release assets into `peter941221/gh-cicost`
- Updated docs to match the automation and release process:
  - `README.md` now references extension sync pipeline
  - `docs/RELEASE.md` now documents required secret `GH_CICOST_REPO_TOKEN` and sync workflow behavior

### Validation

- `go test ./...` ✅
- `go vet ./...` ✅
- `rg -n "REPLACE_WITH_REAL_SHA256|0.1.0" Formula/cicost.rb` ✅ (no matches)
- workflow file review completed: `.github/workflows/sync-gh-extension.yml` ✅
- configured repo secret: `GH_CICOST_REPO_TOKEN` ✅
- workflow dispatch test:
  - first run failed due missing token (`run_id=22469418997`) ❌
  - rerun passed after secret setup (`run_id=22469443908`) ✅

## 2026-02-27 (User Acquisition Target Search + Outreach Preparation)

### Progress

- Searched GitHub for active repositories related to GitHub Actions and CI tooling.
- Curated a high-fit outreach shortlist with owner type, recency, and channel readiness:
  - `rhysd/actionlint`
  - `zizmorcore/zizmor`
  - `actions/actions-runner-controller`
  - `github-aws-runners/terraform-aws-github-runner`
  - `dorny/paths-filter`
  - `tj-actions/changed-files`
  - `softprops/action-gh-release`
  - `JamesIves/github-pages-deploy-action`
  - `shivammathur/setup-php`
  - `docker/build-push-action`
  - plus additional secondary targets (SamKirkland, mxschmitt, gradle, pulumi, goreleaser)
- Prepared practical outreach structure:
  - channel selection logic (discussion > issue > DM)
  - 7-day execution cadence
  - copy-ready outreach message templates

### Validation

- `gh search repos "topic:github-actions stars:50..5000 pushed:>=2025-12-01 archived:false"` ✅
- repository metadata verification via `gh api repos/<owner>/<repo>` ✅
- fields verified: `pushed_at`, `owner.type`, `has_discussions`, `has_issues`, `open_issues_count` ✅

## 2026-02-27 (Outreach Sent - First 10 Targets)

### Progress

- User confirmed that the first 10 personalized outreach messages were sent.
- Campaign execution moved from preparation to active outreach phase.
- Current focus shifts to:
  - response collection
  - 48-hour follow-up for no-response targets
  - scheduling short feedback sessions for responders

### Validation

- Source of truth: user confirmation in-session ("already sent").
- Note: send actions were performed by user side and were not programmatically verified from this environment.

## 2026-02-27 (Memory Restore + README Optimization)

### Progress

- Restored project context from `MEMORY.md`, `RUNBOOK.md`, and command source files under `cmd/`.
- Reworked `README.md` into a release-facing structure:
  - added CI/release/license/Go badges
  - clarified value proposition and current capability (v0.2.0)
  - upgraded architecture section with clearer end-to-end data flow
  - added quick-start flow with verified command set
  - added command map table with key flags
  - documented CI-relevant exit codes (`0/1/2/3`)
  - consolidated docs/community/trust entry points
- Kept distribution paths aligned with current automation:
  - release workflow
  - GoReleaser config
  - extension sync workflow

### Validation

- `go run . version` ✅
- `go run . help` ✅
- link presence check in README (`rg` + `Test-Path`) ✅
- `go test ./...` ✅

## 2026-02-27 (CLI GIF Recording + README Embed)

### Progress

- Added reproducible CLI demo scripts:
  - `docs/scripts/demo_session.ps1` (plays command sequence for recording)
  - `docs/scripts/record_cli_gif.ps1` (build + record + encode gif)
- Recorded a real desktop CLI session with `ffmpeg` and generated:
  - `docs/assets/cicost-cli-demo.gif`
- Embedded GIF in `README.md` under `CLI Demo (GIF)` section.
- Added script reference in README documentation list for future re-recording.

### Validation

- `powershell -File docs/scripts/record_cli_gif.ps1` ✅
- generated file check: `docs/assets/cicost-cli-demo.gif` ✅
- `go test ./...` ✅

## 2026-02-27 (README GIF Update Pushed)

### Progress

- Committed documentation/demo updates with:
  - commit: `1b201f2`
  - message: `docs: add CLI demo gif and README improvements`
- Pushed commit to remote `origin/main`.

### Validation

- `git push origin main` ✅
- remote update range: `2f6b849..1b201f2` ✅

## 2026-02-27 (Headless GIF Re-record)

### Progress

- Re-recorded CLI demo GIF with headless terminal rendering to avoid desktop interference.
- Installed required tooling for headless flow:
  - `vhs` (via `go install`)
  - `ttyd` (via `scoop install ttyd`)
- Added reusable tape-based recording assets:
  - `docs/scripts/cicost_manual_typing.tape`
  - `docs/scripts/record_cli_gif_headless.ps1`
- Removed aborted desktop-injection attempt scripts:
  - deleted `docs/scripts/manual_type_session.ps1`
  - deleted `docs/scripts/record_cli_gif_manual.ps1`
- Updated README demo script reference to the new headless recorder.

### Validation

- `vhs docs/scripts/cicost_manual_typing.tape` ✅
- `docs/assets/cicost-cli-demo.gif` regenerated ✅
- `go test ./...` ✅
