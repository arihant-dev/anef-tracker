package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"
)

type SnapshotRef struct {
	SnapshotID string              `json:"snapshot_id"`
	Directory  string              `json:"directory"`
	Timestamp  time.Time           `json:"timestamp"`
	Metadata   *domain.Application `json:"metadata"`
	RawBytes   []byte              `json:"-"`
}

func GetBaseSnapshotDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cwd, "data", "snapshots")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func GenerateSnapshotID(t time.Time, rawData []byte) string {
	hash := sha256.Sum256(rawData)
	hashShort := hex.EncodeToString(hash[:])[:6]
	return fmt.Sprintf("%s_%s", t.Format("20060102T150405"), hashShort)
}

// SaveHierarchicalSnapshot stores application.json, metadata.json, http.json, and schema.json under data/snapshots/YYYY/MM/DD/SnapshotID/ with immutable read-only permissions (0444)
func SaveHierarchicalSnapshot(app *domain.Application, reqURL string, statusCode int) (*SnapshotRef, error) {
	baseDir, err := GetBaseSnapshotDir()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	snapshotID := GenerateSnapshotID(now, app.RawJSON)

	datePath := filepath.Join(baseDir, now.Format("2006"), now.Format("01"), now.Format("02"))
	snapshotDir := filepath.Join(datePath, snapshotID)

	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return nil, fmt.Errorf("failed creating snapshot directory: %w", err)
	}

	// 1. application.json (Raw JSON Bytes - Read Only 0444)
	if len(app.RawJSON) > 0 {
		if err := os.WriteFile(filepath.Join(snapshotDir, "application.json"), app.RawJSON, 0444); err != nil {
			return nil, err
		}
	}

	// 2. metadata.json (Domain Application + Snapshot ID - Read Only 0444)
	metaMap := map[string]interface{}{
		"snapshot_id":    snapshotID,
		"application_id": app.ID,
		"created_at":     now,
		"domain":         app,
	}
	metaBytes, err := json.MarshalIndent(metaMap, "", "  ")
	if err == nil {
		_ = os.WriteFile(filepath.Join(snapshotDir, "metadata.json"), metaBytes, 0444)
	}

	// 3. http.json (HTTP Metadata - Read Only 0444)
	httpMeta := map[string]interface{}{
		"snapshot_id": snapshotID,
		"url":         reqURL,
		"status_code": statusCode,
		"timestamp":   now,
	}
	httpBytes, err := json.MarshalIndent(httpMeta, "", "  ")
	if err == nil {
		_ = os.WriteFile(filepath.Join(snapshotDir, "http.json"), httpBytes, 0444)
	}

	// 4. schema.json (Extracted Field Schema Map - Read Only 0444)
	if len(app.RawPayload) > 0 {
		schemaMap := extractSchemaGraph("", app.RawPayload)
		schemaBytes, err := json.MarshalIndent(schemaMap, "", "  ")
		if err == nil {
			_ = os.WriteFile(filepath.Join(snapshotDir, "schema.json"), schemaBytes, 0444)
		}
	}

	return &SnapshotRef{
		SnapshotID: snapshotID,
		Directory:  snapshotDir,
		Timestamp:  now,
		Metadata:   app,
		RawBytes:   app.RawJSON,
	}, nil
}

func extractSchemaGraph(prefix string, payload map[string]interface{}) map[string]string {
	schema := make(map[string]string)
	for k, v := range payload {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		if v == nil {
			schema[path] = "null"
		} else if m, ok := v.(map[string]interface{}); ok {
			subSchema := extractSchemaGraph(path, m)
			for subK, subV := range subSchema {
				schema[subK] = subV
			}
		} else {
			schema[path] = reflect.TypeOf(v).String()
		}
	}
	return schema
}

// GetLatestTwoSnapshots finds the two most recent snapshots for diff comparison.
func GetLatestTwoSnapshots() (*SnapshotRef, *SnapshotRef, error) {
	baseDir, err := GetBaseSnapshotDir()
	if err != nil {
		return nil, nil, err
	}

	var snapshotPaths []string
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode().IsRegular() && info.Name() == "application.json" {
			snapshotPaths = append(snapshotPaths, filepath.Dir(path))
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	if len(snapshotPaths) == 0 {
		return nil, nil, fmt.Errorf("no snapshots found")
	}

	sort.Strings(snapshotPaths)

	var latest, previous *SnapshotRef

	latestDir := snapshotPaths[len(snapshotPaths)-1]
	latest, err = loadSnapshotFromDir(latestDir)
	if err != nil {
		return nil, nil, err
	}

	if len(snapshotPaths) > 1 {
		prevDir := snapshotPaths[len(snapshotPaths)-2]
		previous, _ = loadSnapshotFromDir(prevDir)
	}

	return latest, previous, nil
}

func loadSnapshotFromDir(dir string) (*SnapshotRef, error) {
	appPath := filepath.Join(dir, "application.json")
	data, err := os.ReadFile(appPath)
	if err != nil {
		return nil, err
	}

	app, err := domain.MapJSONToApplication(data, "")
	if err != nil {
		return nil, err
	}

	snapshotID := filepath.Base(dir)

	return &SnapshotRef{
		SnapshotID: snapshotID,
		Directory:  dir,
		Timestamp:  time.Now(),
		Metadata:   app,
		RawBytes:   data,
	}, nil
}
