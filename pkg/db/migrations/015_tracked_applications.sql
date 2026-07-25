CREATE TABLE IF NOT EXISTS tracked_applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    anef_id TEXT,
    type TEXT NOT NULL DEFAULT 'RESIDENCE_PERMIT',
    status TEXT,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles(id)
);

INSERT INTO tracked_applications (id, profile_id, name, anef_id, type, status, created_at)
SELECT 1, 1, 'Default Application', NULL, 'RESIDENCE_PERMIT', NULL, datetime('now')
WHERE NOT EXISTS (SELECT 1 FROM tracked_applications WHERE id = 1);

CREATE INDEX IF NOT EXISTS idx_tracked_app_profile ON tracked_applications(profile_id);
