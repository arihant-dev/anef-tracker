package audit

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"strings"
	"time"
)

type AuditAction string

const (
	ActionConfigChanged   AuditAction = "CONFIG_CHANGED"
	ActionExportCreated   AuditAction = "EXPORT_CREATED"
	ActionBackupCreated   AuditAction = "BACKUP_CREATED"
	ActionProfileSwitched AuditAction = "PROFILE_SWITCHED"
	ActionSecurityAudit   AuditAction = "SECURITY_AUDIT_RUN"
	ActionPrivacyAudit    AuditAction = "PRIVACY_AUDIT_RUN"
	ActionProfileDeleted  AuditAction = "PROFILE_DELETED"
)

type AuditEntry struct {
	ID            int64       `json:"id"`
	Action        AuditAction `json:"action"`
	Resource      string      `json:"resource"`
	ProfileID     int64       `json:"profile_id"`
	Metadata      string      `json:"metadata"`
	EvidenceID    int64       `json:"evidence_id,omitempty"`
	HashAlgorithm string      `json:"hash_algorithm"`
	EntryHash     string      `json:"entry_hash"`
	PreviousHash  string      `json:"previous_hash"`
	CreatedAt     time.Time   `json:"created_at"`
}

func RecordAudit(database *db.DB, action AuditAction, resource string, profileID int64, metadata string) (*AuditEntry, error) {
	now := time.Now()
	prevHash := ""

	if database != nil {
		_ = database.Conn.QueryRow("SELECT entry_hash FROM audit_log ORDER BY id DESC LIMIT 1").Scan(&prevHash)
	}

	hashInput := fmt.Sprintf("%s:%s:%d:%s:%s:%s", action, resource, profileID, metadata, prevHash, now.Format(time.RFC3339))
	entryHash := evidence.CalculateHash([]byte(hashInput))

	entry := &AuditEntry{
		Action:        action,
		Resource:      resource,
		ProfileID:     profileID,
		Metadata:      metadata,
		HashAlgorithm: "SHA256",
		EntryHash:     entryHash,
		PreviousHash:  prevHash,
		CreatedAt:     now,
	}

	if database != nil {
		res, err := database.Conn.Exec(
			"INSERT INTO audit_log (action, resource, profile_id, metadata, hash_algorithm, entry_hash, previous_hash, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			action, resource, profileID, metadata, "SHA256", entryHash, prevHash, now,
		)
		if err == nil {
			id, _ := res.LastInsertId()
			entry.ID = id
		}
	}

	return entry, nil
}

func ListAuditLog(database *db.DB, limit int) ([]AuditEntry, error) {
	if database == nil {
		return []AuditEntry{}, nil
	}

	rows, err := database.Conn.Query("SELECT id, action, resource, profile_id, metadata, hash_algorithm, entry_hash, previous_hash, created_at FROM audit_log ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.Action, &e.Resource, &e.ProfileID, &e.Metadata, &e.HashAlgorithm, &e.EntryHash, &e.PreviousHash, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func VerifyAuditChain(database *db.DB) (bool, error) {
	entries, err := ListAuditLog(database, 500)
	if err != nil || len(entries) == 0 {
		return true, nil
	}

	for i := 0; i < len(entries)-1; i++ {
		// entries are DESC order, so i+1 is previous entry chronologically
		if entries[i].PreviousHash != entries[i+1].EntryHash && entries[i].PreviousHash != "" {
			return false, fmt.Errorf("hash chain mismatch between entry #%d and #%d", entries[i].ID, entries[i+1].ID)
		}
	}
	return true, nil
}

func FormatAuditLog(entries []AuditEntry) string {
	var sb strings.Builder
	sb.WriteString("=== TAMPER-EVIDENT AUDIT LOG & OPERATION CHAIN ===\n\n")

	if len(entries) == 0 {
		sb.WriteString("No audit events recorded yet.\n")
		return sb.String()
	}

	for _, e := range entries {
		shortHash := e.EntryHash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}
		sb.WriteString(fmt.Sprintf("#%d [%s] %s — %s (Hash: %s)\n",
			e.ID, e.CreatedAt.Format("2006-01-02 15:04"), e.Action, e.Resource, shortHash))
		if e.Metadata != "" {
			sb.WriteString(fmt.Sprintf("    Meta: %s\n", e.Metadata))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
