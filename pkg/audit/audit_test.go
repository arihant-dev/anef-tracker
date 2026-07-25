package audit_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/audit"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"testing"
)

func TestAuditRecordingAndChain(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	e1, err := audit.RecordAudit(database, audit.ActionConfigChanged, "notification.webhook", 1, "Set webhook URL")
	if err != nil {
		t.Fatalf("RecordAudit 1 failed: %v", err)
	}

	e2, err := audit.RecordAudit(database, audit.ActionExportCreated, "evidence_bundle", 1, "Exported ZIP archive")
	if err != nil {
		t.Fatalf("RecordAudit 2 failed: %v", err)
	}

	if e2.PreviousHash != e1.EntryHash {
		t.Errorf("hash chain broken: entry 2 previous_hash (%s) != entry 1 hash (%s)", e2.PreviousHash, e1.EntryHash)
	}

	valid, err := audit.VerifyAuditChain(database)
	if err != nil || !valid {
		t.Errorf("expected audit chain to verify, got valid=%t err=%v", valid, err)
	}
}
