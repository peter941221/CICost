package store

import (
	"database/sql"
	"encoding/json"

	"github.com/peter941221/CICost/internal/model"
)

func (s *Store) UpsertBillingSnapshot(snapshot model.BillingSnapshot) error {
	_, err := s.db.Exec(`
INSERT INTO billing_snapshots (repo, period, actual_cost_usd, source, fetched_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(repo, period, source) DO UPDATE SET
	actual_cost_usd=excluded.actual_cost_usd,
	fetched_at=excluded.fetched_at`,
		snapshot.Repo,
		snapshot.Period,
		snapshot.ActualCostUSD,
		snapshot.Source,
		asRFC3339(snapshot.FetchedAt),
	)
	return err
}

func (s *Store) GetBillingSnapshot(repo, period string) (model.BillingSnapshot, bool, error) {
	row := s.db.QueryRow(`
SELECT repo, period, actual_cost_usd, source, fetched_at
FROM billing_snapshots
WHERE repo = ? AND period = ?
ORDER BY fetched_at DESC
LIMIT 1`, repo, period)
	var out model.BillingSnapshot
	var fetchedAt sql.NullString
	if err := row.Scan(&out.Repo, &out.Period, &out.ActualCostUSD, &out.Source, &fetchedAt); err != nil {
		if err == sql.ErrNoRows {
			return model.BillingSnapshot{}, false, nil
		}
		return model.BillingSnapshot{}, false, err
	}
	out.FetchedAt = parseRFC3339(fetchedAt.String)
	return out, true, nil
}

func (s *Store) InsertReconcileResult(res model.ReconcileResult) error {
	_, err := s.db.Exec(`
INSERT INTO reconcile_results (repo, period, estimated_cost_usd, actual_cost_usd, delta_ratio, calibration_factor, confidence)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		res.Repo,
		res.Period,
		res.EstimatedCostUSD,
		res.ActualCostUSD,
		res.DeltaRatio,
		res.CalibrationFactor,
		res.Confidence,
	)
	return err
}

func (s *Store) GetLatestReconcile(repo string) (model.ReconcileResult, bool, error) {
	row := s.db.QueryRow(`
SELECT repo, period, estimated_cost_usd, actual_cost_usd, delta_ratio, calibration_factor, confidence, created_at
FROM reconcile_results
WHERE repo = ?
ORDER BY created_at DESC
LIMIT 1`, repo)
	var out model.ReconcileResult
	var createdAt sql.NullString
	if err := row.Scan(&out.Repo, &out.Period, &out.EstimatedCostUSD, &out.ActualCostUSD, &out.DeltaRatio, &out.CalibrationFactor, &out.Confidence, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return model.ReconcileResult{}, false, nil
		}
		return model.ReconcileResult{}, false, err
	}
	out.CreatedAt = parseRFC3339(createdAt.String)
	return out, true, nil
}

func (s *Store) InsertPolicyRun(run model.PolicyRun) error {
	matched := 0
	if run.Matched {
		matched = 1
	}
	_, err := s.db.Exec(`
INSERT INTO policy_runs (repo, period_start, period_end, rule_id, severity, matched, evidence_key, evidence_value, expression)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		run.Repo,
		asRFC3339(run.PeriodStart),
		asRFC3339(run.PeriodEnd),
		run.RuleID,
		run.Severity,
		matched,
		run.EvidenceKey,
		run.EvidenceValue,
		run.Expression,
	)
	return err
}

func (s *Store) InsertSuggestionHistory(rec model.SuggestionRecord) error {
	evidence := rec.EvidenceJSON
	if evidence == "" {
		evidence = "{}"
	}
	if !json.Valid([]byte(evidence)) {
		evidence = "{}"
	}
	_, err := s.db.Exec(`
INSERT INTO suggestion_history (repo, period_start, period_end, suggestion_type, title, estimated_saving_usd, evidence_json)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.Repo,
		asRFC3339(rec.PeriodStart),
		asRFC3339(rec.PeriodEnd),
		rec.SuggestionType,
		rec.Title,
		rec.EstimatedSavingUSD,
		evidence,
	)
	return err
}
