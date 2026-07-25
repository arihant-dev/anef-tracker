CREATE TABLE IF NOT EXISTS evidence_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_type TEXT NOT NULL,
    snapshot_id TEXT,
    event_id INTEGER,
    http_log_id INTEGER,
    payload_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_evidence_source_type ON evidence_records(source_type);
CREATE INDEX IF NOT EXISTS idx_evidence_snapshot ON evidence_records(snapshot_id);
CREATE INDEX IF NOT EXISTS idx_evidence_hash ON evidence_records(payload_hash);
