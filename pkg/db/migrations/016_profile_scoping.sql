ALTER TABLE snapshots ADD COLUMN profile_id INTEGER DEFAULT 1;
ALTER TABLE snapshots ADD COLUMN tracked_application_id INTEGER DEFAULT 1;
ALTER TABLE events ADD COLUMN profile_id INTEGER DEFAULT 1;
ALTER TABLE events ADD COLUMN tracked_application_id INTEGER DEFAULT 1;
ALTER TABLE http_logs ADD COLUMN profile_id INTEGER DEFAULT 1;
ALTER TABLE evidence_records ADD COLUMN profile_id INTEGER DEFAULT 1;
ALTER TABLE evidence_records ADD COLUMN tracked_application_id INTEGER DEFAULT 1;

UPDATE snapshots SET profile_id = 1, tracked_application_id = 1 WHERE profile_id IS NULL;
UPDATE events SET profile_id = 1, tracked_application_id = 1 WHERE profile_id IS NULL;
UPDATE http_logs SET profile_id = 1 WHERE profile_id IS NULL;
UPDATE evidence_records SET profile_id = 1, tracked_application_id = 1 WHERE profile_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_snapshots_profile ON snapshots(profile_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_tracked_app ON snapshots(tracked_application_id);
CREATE INDEX IF NOT EXISTS idx_events_profile ON events(profile_id);
CREATE INDEX IF NOT EXISTS idx_events_tracked_app ON events(tracked_application_id);
CREATE INDEX IF NOT EXISTS idx_http_logs_profile ON http_logs(profile_id);
CREATE INDEX IF NOT EXISTS idx_evidence_profile ON evidence_records(profile_id);
