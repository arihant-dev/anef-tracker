CREATE TABLE IF NOT EXISTS status_durations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status_code TEXT UNIQUE,
    samples INTEGER DEFAULT 1,
    min_days REAL,
    median_days REAL,
    max_days REAL,
    last_updated TIMESTAMP
);
