CREATE TABLE IF NOT EXISTS http_replays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_request_id INTEGER,
    timestamp TIMESTAMP,
    status_code INTEGER,
    response_hash TEXT,
    matched BOOLEAN
);
