package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"os"
	"path/filepath"
	"time"
)

type BackupResult struct {
	BackupPath    string
	SnapshotCount int
	EventCount    int
	HTTPLogCount  int
	EvidenceCount int
	Manifest      Manifest
}

func CreateBackup(database *db.DB) (*BackupResult, error) {
	cwd, _ := os.Getwd()
	backupDir := filepath.Join(cwd, "backups")
	_ = os.MkdirAll(backupDir, 0755)

	timestamp := time.Now().Format("2006-01-02_150405")
	archiveName := fmt.Sprintf("anef-%s.tar.gz", timestamp)
	archivePath := filepath.Join(backupDir, archiveName)

	dbPath := filepath.Join(cwd, "data", "anef.db")
	dbBytes, err := os.ReadFile(dbPath)
	if err != nil {
		dbBytes = []byte{}
	}

	dbHash := evidence.CalculateHash(dbBytes)

	res := &BackupResult{
		BackupPath: archivePath,
		Manifest: Manifest{
			Version:      "1.0",
			CreatedAt:    time.Now(),
			DatabaseHash: dbHash,
		},
	}

	if database != nil {
		events, _ := database.GetEvents(10000)
		res.EventCount = len(events)
		res.Manifest.EventCount = len(events)

		rows, err := database.Conn.Query("SELECT COUNT(*) FROM http_logs")
		if err == nil {
			defer rows.Close()
			if rows.Next() {
				_ = rows.Scan(&res.HTTPLogCount)
				res.Manifest.HTTPLogCount = res.HTTPLogCount
			}
		}
	}

	outFile, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed creating backup file: %w", err)
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	tw := tar.NewWriter(gw)

	// Write Manifest to Tarball
	manifestBytes, _ := json.MarshalIndent(res.Manifest, "", "  ")
	hdr := &tar.Header{
		Name: "manifest.json",
		Mode: 0644,
		Size: int64(len(manifestBytes)),
	}
	_ = tw.WriteHeader(hdr)
	_, _ = tw.Write(manifestBytes)

	// Write Database to Tarball if exists
	if len(dbBytes) > 0 {
		dbHdr := &tar.Header{
			Name: "anef.db",
			Mode: 0644,
			Size: int64(len(dbBytes)),
		}
		_ = tw.WriteHeader(dbHdr)
		_, _ = tw.Write(dbBytes)
	}

	_ = tw.Close()
	_ = gw.Close()

	return res, nil
}
