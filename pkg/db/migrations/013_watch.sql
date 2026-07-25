CREATE TABLE IF NOT EXISTS watch_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    status TEXT NOT NULL,
    changes_detected INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_watch_runs_status ON watch_runs(status);
