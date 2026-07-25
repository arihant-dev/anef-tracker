package evidence

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{DB: database}
}

func (r *Repository) SaveRecord(rec EvidenceRecord) error {
	if r.DB == nil {
		return nil
	}

	_, err := r.DB.Conn.Exec(
		"INSERT INTO evidence_records (source_type, snapshot_id, event_id, http_log_id, payload_hash, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		string(rec.SourceType), rec.SnapshotID, rec.EventID, rec.HTTPLogID, rec.PayloadHash, rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed saving evidence record: %w", err)
	}
	return nil
}
