package report_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/report"
	"testing"
)

func TestMarkdownReport(t *testing.T) {
	rep, err := report.GenerateReport(nil)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	mdStr := rep.RenderMarkdown()
	if len(mdStr) == 0 {
		t.Errorf("expected non-empty markdown string")
	}
}

func TestJSONReport(t *testing.T) {
	rep, _ := report.GenerateReport(nil)
	jsonBytes, err := rep.RenderJSON()
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Errorf("expected non-empty JSON bytes")
	}
}
