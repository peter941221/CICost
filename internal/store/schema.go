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
);

CREATE TABLE IF NOT EXISTS billing_snapshots (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    repo            TEXT NOT NULL,
    period          TEXT NOT NULL,
    actual_cost_usd REAL NOT NULL,
    source          TEXT NOT NULL,
    fetched_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(repo, period, source)
);

CREATE TABLE IF NOT EXISTS reconcile_results (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    repo                  TEXT NOT NULL,
    period                TEXT NOT NULL,
    estimated_cost_usd    REAL NOT NULL,
    actual_cost_usd       REAL NOT NULL,
    delta_ratio           REAL NOT NULL,
    calibration_factor    REAL NOT NULL,
    confidence            TEXT NOT NULL,
    created_at            TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS policy_runs (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    repo              TEXT NOT NULL,
    period_start      TEXT NOT NULL,
    period_end        TEXT NOT NULL,
    rule_id           TEXT NOT NULL,
    severity          TEXT NOT NULL,
    matched           INTEGER NOT NULL DEFAULT 0,
    evidence_key      TEXT NOT NULL,
    evidence_value    REAL NOT NULL,
    expression        TEXT NOT NULL,
    created_at        TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS suggestion_history (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    repo                  TEXT NOT NULL,
    period_start          TEXT NOT NULL,
    period_end            TEXT NOT NULL,
    suggestion_type       TEXT NOT NULL,
    title                 TEXT NOT NULL,
    estimated_saving_usd  REAL NOT NULL,
    evidence_json         TEXT NOT NULL,
    created_at            TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_runs_repo_created ON workflow_runs(repo, created_at);
CREATE INDEX IF NOT EXISTS idx_runs_repo_workflow ON workflow_runs(repo, workflow_name);
CREATE INDEX IF NOT EXISTS idx_jobs_run_id ON jobs(run_id, run_attempt);
CREATE INDEX IF NOT EXISTS idx_jobs_repo_runner ON jobs(repo, runner_os);
CREATE INDEX IF NOT EXISTS idx_billing_repo_period ON billing_snapshots(repo, period);
CREATE INDEX IF NOT EXISTS idx_reconcile_repo_period ON reconcile_results(repo, period);
CREATE INDEX IF NOT EXISTS idx_policy_repo_created ON policy_runs(repo, created_at);
CREATE INDEX IF NOT EXISTS idx_suggestion_repo_created ON suggestion_history(repo, created_at);
`
