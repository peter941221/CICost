# CICost Launch Kit (v0.2.0)

## Positioning

- Problem: CI cost visibility is fragmented and optimization is hard to operationalize.
- Product: CICost provides cost attribution, policy gates, reconciliation, and executable optimization suggestions.
- Audience: engineering teams using GitHub Actions and seeking quick FinOps guardrails.

## 30-Second Demo Flow

```bash
go run . scan --repo owner/repo --days 30
go run . report --repo owner/repo --format table
go run . policy check --repo owner/repo --days 30
go run . suggest --repo owner/repo --format yaml --output patches/
```

## Release Day Checklist

1. Validate release prerequisites:
   - `go test ./...`
   - `go test -race ./...`
   - `go vet ./...`
   - `goreleaser check`
2. Tag release:
   - `git tag v0.2.0`
   - `git push origin v0.2.0`
3. Verify install paths:
   - binary download from GitHub Releases
   - `gh extension install peter941221/gh-cicost`
   - Homebrew tap formula checksum update
4. Publish launch message:
   - one short post with problem/solution/demo/install links

## Launch Post Template

Title:
`CICost v0.2.0: GitHub Actions cost guardrails with reconciliation and policy gates`

Body:

1. Why we built it: CI cost blind spots and reactive optimization are expensive.
2. What is new in v0.2.0:
   - pricing snapshots
   - reconcile + calibration
   - policy gates
   - actionable YAML suggestions
   - org-level reports
3. Demo commands (copy/paste block).
4. Install:
   - Releases binary
   - GitHub CLI extension
   - Homebrew
5. Feedback request:
   - open issues
   - feature suggestions
