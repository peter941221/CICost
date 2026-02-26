package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/peter941221/CICost/internal/model"
)

var ErrStoreNotImplemented = errors.New("sqlite store not implemented")

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(SchemaSQL); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) UpsertRuns(runs []model.WorkflowRun) (newCount, updatedCount int, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer rollbackIfNeeded(tx, &err)

	checkStmt, err := tx.Prepare(`SELECT COUNT(1) FROM workflow_runs WHERE id = ? AND run_attempt = ?`)
	if err != nil {
		return 0, 0, err
	}
	defer checkStmt.Close()

	upsertStmt, err := tx.Prepare(`
INSERT INTO workflow_runs (
    id, repo, workflow_id, workflow_name, head_branch, event, status, conclusion, run_attempt,
    run_started_at, updated_at, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id, run_attempt) DO UPDATE SET
    repo=excluded.repo,
    workflow_id=excluded.workflow_id,
    workflow_name=excluded.workflow_name,
    head_branch=excluded.head_branch,
    event=excluded.event,
    status=excluded.status,
    conclusion=excluded.conclusion,
    run_started_at=excluded.run_started_at,
    updated_at=excluded.updated_at,
    created_at=excluded.created_at`)
	if err != nil {
		return 0, 0, err
	}
	defer upsertStmt.Close()

	for _, r := range runs {
		var exists int
		if err = checkStmt.QueryRow(r.ID, r.RunAttempt).Scan(&exists); err != nil {
			return 0, 0, err
		}
		if exists > 0 {
			updatedCount++
		} else {
			newCount++
		}
		_, err = upsertStmt.Exec(
			r.ID, r.Repo, r.WorkflowID, r.WorkflowName, r.HeadBranch, r.Event, r.Status, r.Conclusion, r.RunAttempt,
			asRFC3339(r.RunStartedAt), asRFC3339(r.UpdatedAt), asRFC3339(r.CreatedAt),
		)
		if err != nil {
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}
	return newCount, updatedCount, nil
}

func (s *Store) UpsertJobs(jobs []model.Job) (newCount, updatedCount int, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer rollbackIfNeeded(tx, &err)

	checkStmt, err := tx.Prepare(`SELECT COUNT(1) FROM jobs WHERE id = ? AND run_attempt = ?`)
	if err != nil {
		return 0, 0, err
	}
	defer checkStmt.Close()

	upsertStmt, err := tx.Prepare(`
INSERT INTO jobs (
    id, run_id, run_attempt, repo, name, status, conclusion, started_at, completed_at,
    runner_os, runner_name, runner_group, is_self_hosted, duration_sec
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id, run_attempt) DO UPDATE SET
    run_id=excluded.run_id,
    repo=excluded.repo,
    name=excluded.name,
    status=excluded.status,
    conclusion=excluded.conclusion,
    started_at=excluded.started_at,
    completed_at=excluded.completed_at,
    runner_os=excluded.runner_os,
    runner_name=excluded.runner_name,
    runner_group=excluded.runner_group,
    is_self_hosted=excluded.is_self_hosted,
    duration_sec=excluded.duration_sec`)
	if err != nil {
		return 0, 0, err
	}
	defer upsertStmt.Close()

	for _, j := range jobs {
		var exists int
		if err = checkStmt.QueryRow(j.ID, j.RunAttempt).Scan(&exists); err != nil {
			return 0, 0, err
		}
		if exists > 0 {
			updatedCount++
		} else {
			newCount++
		}
		selfHosted := 0
		if j.IsSelfHosted {
			selfHosted = 1
		}
		_, err = upsertStmt.Exec(
			j.ID, j.RunID, j.RunAttempt, j.Repo, j.Name, j.Status, j.Conclusion, asRFC3339(j.StartedAt), asRFC3339(j.CompletedAt),
			j.RunnerOS, j.RunnerName, j.RunnerGroup, selfHosted, j.DurationSec,
		)
		if err != nil {
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}
	return newCount, updatedCount, nil
}

func (s *Store) ListRuns(repo string, start, end time.Time) ([]model.WorkflowRun, error) {
	rows, err := s.db.Query(`
SELECT id, repo, workflow_id, workflow_name, head_branch, event, status, conclusion, run_attempt,
       run_started_at, updated_at, created_at
FROM workflow_runs
WHERE repo = ? AND created_at >= ? AND created_at <= ?
ORDER BY created_at DESC`, repo, start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.WorkflowRun, 0, 256)
	for rows.Next() {
		var r model.WorkflowRun
		var runStartedAt, updatedAt, createdAt sql.NullString
		if err := rows.Scan(&r.ID, &r.Repo, &r.WorkflowID, &r.WorkflowName, &r.HeadBranch, &r.Event, &r.Status, &r.Conclusion, &r.RunAttempt,
			&runStartedAt, &updatedAt, &createdAt); err != nil {
			return nil, err
		}
		r.RunStartedAt = parseRFC3339(runStartedAt.String)
		r.UpdatedAt = parseRFC3339(updatedAt.String)
		r.CreatedAt = parseRFC3339(createdAt.String)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) ListJobs(repo string, start, end time.Time) ([]model.Job, error) {
	rows, err := s.db.Query(`
SELECT j.id, j.run_id, j.run_attempt, j.repo, j.name, j.status, j.conclusion, j.started_at, j.completed_at,
       j.runner_os, j.runner_name, j.runner_group, j.is_self_hosted, j.duration_sec
FROM jobs j
JOIN workflow_runs r ON r.id = j.run_id AND r.run_attempt = j.run_attempt
WHERE j.repo = ? AND r.created_at >= ? AND r.created_at <= ?
ORDER BY r.created_at DESC`, repo, start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Job, 0, 512)
	for rows.Next() {
		var j model.Job
		var startedAt, completedAt sql.NullString
		var selfHosted int
		if err := rows.Scan(&j.ID, &j.RunID, &j.RunAttempt, &j.Repo, &j.Name, &j.Status, &j.Conclusion, &startedAt, &completedAt,
			&j.RunnerOS, &j.RunnerName, &j.RunnerGroup, &selfHosted, &j.DurationSec); err != nil {
			return nil, err
		}
		j.StartedAt = parseRFC3339(startedAt.String)
		j.CompletedAt = parseRFC3339(completedAt.String)
		j.IsSelfHosted = selfHosted == 1
		out = append(out, j)
	}
	return out, rows.Err()
}

func (s *Store) CountRuns(repo string, start, end time.Time) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM workflow_runs WHERE repo = ? AND created_at >= ? AND created_at <= ?`,
		repo, start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339)).Scan(&n)
	return n, err
}

func (s *Store) InsertBudgetCheck(repo, checkType string, periodStart, periodEnd time.Time, threshold, actual float64, exceeded bool) error {
	ex := 0
	if exceeded {
		ex = 1
	}
	_, err := s.db.Exec(`
INSERT INTO budget_checks (repo, check_type, period_start, period_end, threshold_usd, actual_usd, exceeded)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		repo, checkType, periodStart.UTC().Format(time.RFC3339), periodEnd.UTC().Format(time.RFC3339), threshold, actual, ex)
	return err
}

func rollbackIfNeeded(tx *sql.Tx, err *error) {
	if *err != nil {
		_ = tx.Rollback()
	}
}

func asRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func parseRFC3339(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}
	}
	return t
}
