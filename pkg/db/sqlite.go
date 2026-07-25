package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/arihant-dev/anef-tracker/pkg/domain"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DB struct {
	Conn *sql.DB
}

func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".anef")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "anef_tracker.db"), nil
}

func InitDB() (*DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}
	return InitDBWithPath(dbPath)
}

func InitDBWithPath(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	database := &DB{Conn: conn}
	if err := database.Migrate(); err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return database, nil
}

func (db *DB) Migrate() error {
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed reading embedded migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for idx, fileName := range files {
		var count int
		_ = db.Conn.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", idx).Scan(&count)
		if count > 0 {
			continue // Already applied
		}

		content, err := migrationFS.ReadFile(path.Join("migrations", fileName))
		if err != nil {
			return fmt.Errorf("failed reading migration file %s: %w", fileName, err)
		}
		if _, err := db.Conn.Exec(string(content)); err != nil {
			return fmt.Errorf("error running migration %s: %w", fileName, err)
		}

		_, _ = db.Conn.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)", idx, time.Now())
	}

	return nil
}

func (db *DB) SaveSnapshotRef(appID, directory string) (int64, error) {
	res, err := db.Conn.Exec(
		"INSERT INTO snapshots (application_id, snapshot_dir, created_at) VALUES (?, ?, ?)",
		appID, directory, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) SaveApplication(app *domain.Application) error {
	query := `INSERT INTO applications (id, user_login, numero_demande, legal_category, status_code, status_label, version, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status_code = excluded.status_code,
		status_label = excluded.status_label,
		version = excluded.version,
		updated_at = excluded.updated_at`

	_, err := db.Conn.Exec(query, app.ID, app.ForeignerID, app.NumeroDemande, app.LegalCategory, app.Status.Code, app.Status.Label, app.Version, time.Now(), time.Now())
	return err
}

func (db *DB) RecordHTTPLog(method, url string, statusCode int, latencyMs int64, reqHeaders, respHeaders, respBody string) error {
	_, err := db.Conn.Exec(
		"INSERT INTO http_logs (method, url, status_code, latency_ms, req_headers, resp_headers, resp_body, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		method, url, statusCode, latencyMs, reqHeaders, respHeaders, respBody, time.Now(),
	)
	return err
}

func (db *DB) SaveEvent(event domain.Event) error {
	_, err := db.Conn.Exec(
		"INSERT INTO events (application_id, event_type, severity, confidence, field_path, old_val, new_val, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		event.ApplicationID, event.Type, string(event.Severity), event.Confidence, event.FieldPath, event.OldVal, event.NewVal, time.Now(),
	)
	return err
}

func (db *DB) GetEvents(limit int) ([]domain.Event, error) {
	rows, err := db.Conn.Query("SELECT id, application_id, event_type, severity, confidence, field_path, old_val, new_val, created_at FROM events ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		var sevStr string
		if err := rows.Scan(&e.ID, &e.ApplicationID, &e.Type, &sevStr, &e.Confidence, &e.FieldPath, &e.OldVal, &e.NewVal, &e.Timestamp); err != nil {
			return nil, err
		}
		e.Severity = domain.Severity(sevStr)
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
