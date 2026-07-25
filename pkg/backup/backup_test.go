package backup_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/backup"
	"os"
	"testing"
)

func TestBackupCreation(t *testing.T) {
	res, err := backup.CreateBackup(nil)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if res.BackupPath == "" {
		t.Errorf("expected non-empty backup path")
	}

	if _, err := os.Stat(res.BackupPath); os.IsNotExist(err) {
		t.Errorf("expected backup file to exist at %s", res.BackupPath)
	}

	manifest, err := backup.VerifyBackup(res.BackupPath)
	if err != nil {
		t.Fatalf("VerifyBackup failed: %v", err)
	}

	if manifest.Version != "1.0" {
		t.Errorf("expected manifest version 1.0, got %s", manifest.Version)
	}

	_ = os.Remove(res.BackupPath)
}

func TestManifestVerification(t *testing.T) {
	res, err := backup.CreateBackup(nil)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}
	defer os.Remove(res.BackupPath)

	mf, err := backup.VerifyBackup(res.BackupPath)
	if err != nil {
		t.Fatalf("VerifyBackup failed: %v", err)
	}

	if mf.CreatedAt.IsZero() {
		t.Errorf("expected non-zero CreatedAt timestamp in manifest")
	}
}

func TestBackupRestore(t *testing.T) {
	res, err := backup.CreateBackup(nil)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}
	defer os.Remove(res.BackupPath)

	restoreRes, err := backup.RestoreBackup(res.BackupPath)
	if err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	if restoreRes.Manifest.Version != "1.0" {
		t.Errorf("expected restored manifest version 1.0, got %s", restoreRes.Manifest.Version)
	}
}
