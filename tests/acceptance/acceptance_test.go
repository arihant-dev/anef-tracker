package acceptance_test

import (
	"github.com/arihant-dev/anef-tracker/internal/version"
	"github.com/arihant-dev/anef-tracker/pkg/backup"
	"github.com/arihant-dev/anef-tracker/pkg/config"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"os"
	"path/filepath"
	"testing"
)

func TestFullSystemAcceptance(t *testing.T) {
	// 1. Version Info Verification
	vInfo := version.GetVersionInfo()
	if len(vInfo) == 0 {
		t.Errorf("expected non-empty version metadata")
	}

	dbPath := filepath.Join(t.TempDir(), "acceptance_test.db")
	database, err := db.InitDBWithPath(dbPath)
	if err != nil {
		t.Fatalf("InitDBWithPath failed: %v", err)
	}
	defer database.Conn.Close()

	// 3. Config Engine Verification
	cfg, err := config.LoadConfig()
	if err != nil || cfg == nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 4. Evidence Verification
	evRep, err := evidence.VerifyIntegrity(database)
	if err != nil || !evRep.HashesValid {
		t.Errorf("expected clean evidence integrity verification")
	}

	// 5. Evidence Graph Validation
	gb := knowledge.NewGraphBuilder(database)
	g, _ := gb.BuildFromDB()
	valRep := g.ValidateGraph()
	if !valRep.Valid {
		t.Errorf("expected graph consistency validation to pass with 0 orphan nodes")
	}

	// 6. Backup Creation & Verification
	backupRes, err := backup.CreateBackup(database)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}
	defer os.Remove(backupRes.BackupPath)

	if backupRes.BackupPath == "" {
		t.Errorf("expected non-empty backup path")
	}
}
