# RUNBOOK

## Environment Setup

1. Install Go 1.23+ (required).
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
```

## First Milestone (D1)

1. ✅ `internal/auth` full chain:
   - CLI token
   - `GITHUB_TOKEN` / `GH_TOKEN`
   - `gh auth token`
   - user config fallback
2. ✅ `internal/github` runs/jobs endpoint fetch + pagination.
3. ✅ `cicost scan --repo owner/repo --days 30` 输出入库摘要。

## Validation Commands

```powershell
git status --short
git remote -v
go test ./...
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

git tag v0.1.0
git push origin v0.1.0
```

Then GitHub Actions `release.yml` publishes binaries and checksums.
