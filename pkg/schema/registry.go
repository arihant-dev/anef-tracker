package schema

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"time"
)

type Registry struct {
	DB *db.DB
}

func NewRegistry(database *db.DB) *Registry {
	return &Registry{DB: database}
}

func (r *Registry) RegisterField(endpoint, fieldPath, fieldType string) (bool, error) {
	if r.DB == nil {
		return false, nil
	}

	var id int64
	var occurrences int
	err := r.DB.Conn.QueryRow("SELECT id, occurrences FROM schema_fields WHERE endpoint = ? AND field_path = ?", endpoint, fieldPath).Scan(&id, &occurrences)

	now := time.Now()
	if err != nil {
		// New field discovered!
		_, err := r.DB.Conn.Exec(
			"INSERT INTO schema_fields (endpoint, field_path, field_type, first_seen, last_seen, occurrences, confidence) VALUES (?, ?, ?, ?, ?, 1, 1.0)",
			endpoint, fieldPath, fieldType, now, now,
		)
		if err != nil {
			return false, fmt.Errorf("failed inserting new schema field: %w", err)
		}
		return true, nil
	}

	// Update existing field
	_, err = r.DB.Conn.Exec(
		"UPDATE schema_fields SET occurrences = occurrences + 1, last_seen = ?, field_type = ? WHERE id = ?",
		now, fieldType, id,
	)
	return false, err
}

func (r *Registry) ListFields(endpoint string) ([]domain.FieldObservation, error) {
	if r.DB == nil {
		return nil, nil
	}

	query := "SELECT id, endpoint, field_path, field_type, first_seen, last_seen, occurrences, confidence FROM schema_fields ORDER BY occurrences DESC, field_path ASC"

	sqlRows, err := r.DB.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	var fields []domain.FieldObservation
	for sqlRows.Next() {
		var f domain.FieldObservation
		if err := sqlRows.Scan(&f.ID, &f.Endpoint, &f.Path, &f.Type, &f.FirstSeen, &f.LastSeen, &f.Occurrences, &f.Confidence); err != nil {
			return nil, err
		}
		if endpoint == "" || f.Endpoint == endpoint {
			fields = append(fields, f)
		}
	}
	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}
