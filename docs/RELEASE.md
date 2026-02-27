# Release Guide

## Prerequisites

1. CI on `main` is green.
2. `GITHUB_TOKEN` permission for release workflow is enabled (default repo token is enough).
3. Optional: create tap repo `peter941221/homebrew-tap` for Homebrew distribution.

## Tag Release

```bash
git checkout main
git pull
git tag v0.2.0
git push origin v0.2.0
```

Release workflow (`.github/workflows/release.yml`) will run GoReleaser and publish:

- `cicost_<version>_<os>_<arch>.tar.gz|zip`
- `gh-cicost_<version>_<os>_<arch>.tar.gz|zip`
- `checksums.txt`

## Install Methods

### Standalone binary

```bash
curl -L https://github.com/peter941221/CICost/releases/latest/download/cicost_<version>_linux_amd64.tar.gz | tar xz
```

### GitHub CLI extension

```bash
gh extension install peter941221/CICost
gh cicost version
```

`gh cicost` uses the `gh-cicost` release artifacts.

### Homebrew

1. Copy `Formula/cicost.rb` to your tap repository.
2. Replace `sha256` placeholders using `checksums.txt`.
3. Publish formula:

```bash
brew tap peter941221/tap
brew install cicost
```

## Validate Local Before Tag

```bash
go test ./...
go test -race ./...
go vet ./...
goreleaser check
```
