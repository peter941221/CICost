package store

const SchemaSQL = `
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS workflow_runs (
    id              INTEGER NOT NULL,
    repo            TEXT NOT NULL,
    workflow_id     INTEGER NOT NULL,
    workflow_name   TEXT NOT NULL,
    head_branch     TEXT,
    event           TEXT,
    status          TEXT,
    conclusion      TEXT,
    run_attempt     INTEGER DEFAULT 1,
    run_started_at  TEXT,
    updated_at      TEXT,
    created_at      TEXT,
    fetched_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(id, run_attempt)
);

CREATE TABLE IF NOT EXISTS jobs (
    id              INTEGER NOT NULL,
    run_id          INTEGER NOT NULL,
    run_attempt     INTEGER DEFAULT 1,
    repo            TEXT NOT NULL,
    name            TEXT NOT NULL,
    status          TEXT,
    conclusion      TEXT,
    started_at      TEXT,
    completed_at    TEXT,
    runner_os       TEXT,
    runner_name     TEXT,
    runner_group    TEXT,
    is_self_hosted  INTEGER DEFAULT 0,
    duration_sec    INTEGER,
    fetched_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(id, run_attempt)
);

CREATE TABLE IF NOT EXISTS sync_cursors (
    repo            TEXT PRIMARY KEY,
    last_run_id     INTEGER,
    last_created_at TEXT,
    last_sync_at    TEXT,
    total_runs      INTEGER DEFAULT 0,
    total_jobs      INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS budget_checks (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    repo            TEXT NOT NULL,
    check_type      TEXT NOT NULL,
    period_start    TEXT NOT NULL,
    period_end      TEXT NOT NULL,
    threshold_usd   REAL NOT NULL,
    actual_usd      REAL NOT NULL,
    exceeded        INTEGER DEFAULT 0,
    checked_at      TEXT NOT NULL DEFAULT (datetime('now'))
);`

