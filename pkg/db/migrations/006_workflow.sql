CREATE TABLE IF NOT EXISTS workflow_transitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_status TEXT,
    to_status TEXT,
    count INTEGER DEFAULT 1,
    first_seen TIMESTAMP,
    last_seen TIMESTAMP,
    average_duration REAL,
    confidence REAL DEFAULT 1.0,
    UNIQUE(from_status, to_status)
);
