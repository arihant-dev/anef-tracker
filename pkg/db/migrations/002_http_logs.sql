CREATE TABLE IF NOT EXISTS http_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT,
    url TEXT,
    status_code INTEGER,
    latency_ms INTEGER,
    req_headers TEXT,
    resp_headers TEXT,
    resp_body TEXT,
    created_at TIMESTAMP
);
