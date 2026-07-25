CREATE TABLE IF NOT EXISTS schema_fields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    endpoint TEXT,
    field_path TEXT,
    field_type TEXT,
    first_seen TIMESTAMP,
    last_seen TIMESTAMP,
    occurrences INTEGER DEFAULT 1,
    confidence REAL DEFAULT 1.0,
    UNIQUE(endpoint, field_path)
);
