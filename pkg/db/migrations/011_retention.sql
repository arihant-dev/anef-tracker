CREATE TABLE IF NOT EXISTS retention_policy (
    resource TEXT PRIMARY KEY,
    keep_days INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL
);

INSERT OR IGNORE INTO retention_policy (resource, keep_days, created_at) VALUES ('snapshots', 0, CURRENT_TIMESTAMP);
INSERT OR IGNORE INTO retention_policy (resource, keep_days, created_at) VALUES ('events', 0, CURRENT_TIMESTAMP);
INSERT OR IGNORE INTO retention_policy (resource, keep_days, created_at) VALUES ('http_logs', 365, CURRENT_TIMESTAMP);
INSERT OR IGNORE INTO retention_policy (resource, keep_days, created_at) VALUES ('raw_payloads', 365, CURRENT_TIMESTAMP);
