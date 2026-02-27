# Contributing to CICost

Thanks for helping improve CICost.

## Development Setup

1. Install Go 1.24+.
2. Install GitHub CLI (`gh`) if you need live API testing.
3. Clone and run baseline checks:

```bash
go test ./...
go test -race ./...
go vet ./...
go run . help
```

## Branch and Commit Guidelines

1. Create a branch from `main`.
2. Keep pull requests focused and small.
3. Use conventional prefixes when possible:
   - `feat:`
   - `fix:`
   - `docs:`
   - `test:`
   - `chore:`

## Pull Request Checklist

1. Add or update tests for behavior changes.
2. Update docs and examples if CLI behavior changes.
3. Ensure validation commands pass locally:

```bash
go test ./...
go test -race ./...
go vet ./...
```

4. Include risk notes in the PR description:
   - what changed
   - what was validated
   - residual risk

## Issue Reporting

When filing an issue, include:

1. Operating system and architecture.
2. CICost version (`cicost version`).
3. Command used and full error output.
4. Minimal reproduction data.
