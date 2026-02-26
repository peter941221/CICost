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
make build
make test
```

## First Milestone (D1)

1. Implement `internal/auth` full chain:
   - CLI token
   - `GITHUB_TOKEN` / `GH_TOKEN`
   - `gh auth token`
   - user config fallback
2. Implement `internal/github` runs endpoint fetch.
3. Wire `cicost scan --repo owner/repo --days 30` minimal output.

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

