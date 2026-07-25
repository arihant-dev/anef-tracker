CREATE TABLE IF NOT EXISTS applications (
    id TEXT PRIMARY KEY,
    user_login TEXT,
    numero_demande TEXT,
    legal_category TEXT,
    status_code TEXT,
    status_label TEXT,
    version INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS snapshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id TEXT,
    snapshot_dir TEXT,
    created_at TIMESTAMP
);
