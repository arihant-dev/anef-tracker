CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id TEXT,
    event_type TEXT,
    severity TEXT,
    confidence REAL DEFAULT 1.0,
    field_path TEXT,
    old_val TEXT,
    new_val TEXT,
    created_at TIMESTAMP
);
