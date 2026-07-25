CREATE TABLE IF NOT EXISTS audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    profile_id INTEGER,
    metadata TEXT,
    evidence_id INTEGER,
    hash_algorithm TEXT NOT NULL DEFAULT 'SHA256',
    entry_hash TEXT NOT NULL,
    previous_hash TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_log(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_profile ON audit_log(profile_id);
