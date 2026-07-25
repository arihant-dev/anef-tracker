package export_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/export"
	"strings"
	"testing"
	"time"
)

func TestCSVExport(t *testing.T) {
	events := []domain.Event{
		{
			ID:        1,
			Type:      "FIELD_DISCOVERED",
			Severity:  domain.SeverityLow,
			FieldPath: "statut",
			Timestamp: time.Now(),
		},
	}

	csvBytes, err := export.EventsToCSV(events)
	if err != nil {
		t.Fatalf("EventsToCSV failed: %v", err)
	}

	if !strings.Contains(string(csvBytes), "FIELD_DISCOVERED") || !strings.Contains(string(csvBytes), "statut") {
		t.Errorf("expected CSV to contain field path statut and type FIELD_DISCOVERED")
	}
}

func TestJSONExport(t *testing.T) {
	events := []domain.Event{{ID: 1, Type: "STATUS_CHANGE"}}
	jsonBytes, err := export.ToJSON(events)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !strings.Contains(string(jsonBytes), "STATUS_CHANGE") {
		t.Errorf("expected JSON output to contain STATUS_CHANGE")
	}
}
