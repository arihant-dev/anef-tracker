CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    active BOOLEAN NOT NULL DEFAULT 0
);

INSERT INTO profiles (name, created_at, active)
SELECT 'Default Profile', datetime('now'), 1
WHERE NOT EXISTS (SELECT 1 FROM profiles WHERE id = 1);
