package notify

import (
	"fmt"
)

func FormatStatusChangeTemplate(oldStatus, newStatus, snapshotID string, eventID int64) string {
	return fmt.Sprintf("STATUS_CHANGE\n\nFROM: %s\nTO:   %s\n\nEvidence Provenance:\n  Snapshot ID: %s\n  Event ID:    #%d",
		oldStatus, newStatus, snapshotID, eventID)
}

func FormatDocumentDiscoveredTemplate(docName string, httpLogID int64) string {
	return fmt.Sprintf("DOCUMENT_DISCOVERED\n\nDocument: %s\n\nEvidence Provenance:\n  HTTP Log ID: #%d", docName, httpLogID)
}
