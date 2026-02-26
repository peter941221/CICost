package store

import (
	"database/sql"
	"time"
)

type SyncCursor struct {
	Repo          string
	LastRunID     int64
	LastCreatedAt time.Time
	LastSyncAt    time.Time
	TotalRuns     int
	TotalJobs     int
}

func (s *Store) GetCursor(repo string) (SyncCursor, bool, error) {
	row := s.db.QueryRow(`SELECT repo, last_run_id, last_created_at, last_sync_at, total_runs, total_jobs FROM sync_cursors WHERE repo = ?`, repo)
	var c SyncCursor
	var createdAt, syncAt sql.NullString
	err := row.Scan(&c.Repo, &c.LastRunID, &createdAt, &syncAt, &c.TotalRuns, &c.TotalJobs)
	if err != nil {
		if err == sql.ErrNoRows {
			return SyncCursor{}, false, nil
		}
		return SyncCursor{}, false, err
	}
	c.LastCreatedAt = parseRFC3339(createdAt.String)
	c.LastSyncAt = parseRFC3339(syncAt.String)
	return c, true, nil
}

func (s *Store) UpsertCursor(c SyncCursor) error {
	_, err := s.db.Exec(`
INSERT INTO sync_cursors (repo, last_run_id, last_created_at, last_sync_at, total_runs, total_jobs)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(repo) DO UPDATE SET
	last_run_id=excluded.last_run_id,
	last_created_at=excluded.last_created_at,
	last_sync_at=excluded.last_sync_at,
	total_runs=excluded.total_runs,
	total_jobs=excluded.total_jobs`,
		c.Repo, c.LastRunID, asRFC3339(c.LastCreatedAt), asRFC3339(c.LastSyncAt), c.TotalRuns, c.TotalJobs)
	return err
}
