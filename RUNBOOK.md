# RUNBOOK

## Environment Setup

1. Install Go 1.24+ (required, CI validated on Go 1.26.x).
2. Install GitHub CLI (`gh`) and login:

```powershell
gh auth login
```

3. Optional token fallback:

```powershell
$env:GITHUB_TOKEN="ghp_xxx"
```

## Local Development

```powershell
go test ./...
go run . help
go run . scan --repo owner/repo --days 30
go run . report --repo owner/repo --format table
go run . reconcile --repo owner/repo --month 2026-02 --actual-usd 120.50 --apply-calibration
go run . policy lint --policy .cicost.policy.yml
go run . policy check --repo owner/repo --days 30
go run . suggest --repo owner/repo --format yaml --output patches/
go run . org-report --repos repos.txt --days 30 --format md
```

## v2 Commands

1. `reconcile`
   - `--repo owner/repo`
   - `--month YYYY-MM`
   - `--source csv|github`
   - `--actual-usd 123.45`
   - `--apply-calibration`
2. `policy`
   - `lint`
   - `check`
   - `explain`
3. `suggest`
   - `--format text|yaml`
   - `--output patches/`
4. `org-report`
   - `--repos repos.txt`
   - `--format md|json`

## Validation Commands

```powershell
git status --short
git remote -v
go test ./...
go test -race ./...
go vet ./...
go run . help
```

## Release Baseline

```powershell
git add .
git commit -m "chore: bootstrap cicost repository scaffold"
git push -u origin main
```

## Release (Tag-Based)

```powershell
go test ./...
go test -race ./...
go vet ./...
goreleaser check

git tag v0.2.0
git push origin v0.2.0
```

Then GitHub Actions `release.yml` publishes binaries and checksums.
