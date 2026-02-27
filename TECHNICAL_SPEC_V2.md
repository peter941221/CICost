# CICost Technical Specification v2

## Goal

Upgrade from a useful estimator to an operational CI FinOps tool with reconciliation, policy gates, and actionable suggestions.

## Key Enhancements

- Pricing v2:
  - Snapshot-based SKU rates
  - Effective date selection for historical accuracy
- Reconciliation:
  - Estimate vs actual billing comparison
  - Calibration factor and confidence band
- Policy gate:
  - Policy-as-code with lint/check/explain
  - Exit code `3` for `error` findings
- Suggestion engine:
  - Evidence-backed suggestions
  - YAML/patch export
- Org report:
  - Multi-repo aggregation
  - Partial-result mode

## Data Model Additions

- `billing_snapshots`
- `reconcile_results`
- `policy_runs`
- `suggestion_history`

## Quality Gate

- Unit + integration tests for critical logic
- Race and vet checks in CI
- Release validation with GoReleaser

## Rollout Notes

- Keep calibration opt-in (`--apply-calibration`)
- Keep policy defaults safe (`warn` before stricter `error` rollout)
- Preserve partial results when one repository fails
