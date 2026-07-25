package schema_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/schema"
	"path/filepath"
	"testing"
)

func TestSchemaGeneration(t *testing.T) {
	payload := map[string]interface{}{
		"statut":   "TITRE_A_FABRIQUER",
		"_version": float64(56),
		"dossier": map[string]interface{}{
			"site": "9403",
		},
	}

	doc := schema.GenerateJSONSchema("https://example.com/api", payload)
	if doc == nil {
		t.Fatalf("expected non-nil SchemaDocument")
	}

	if len(doc.Fields) < 3 {
		t.Errorf("expected at least 3 fields, got %d", len(doc.Fields))
	}

	if st, ok := doc.Fields["statut"]; !ok || st.Type != "string" {
		t.Errorf("expected field 'statut' of type string")
	}
}

func TestSchemaDiscoveryAndUnknownFields(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_schema.db")
	database, err := db.InitDBWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed initializing test db: %v", err)
	}

	reg := schema.NewRegistry(database)
	engine := schema.NewDiscoveryEngine(reg, nil)

	payload := map[string]interface{}{
		"statut":          "TITRE_A_FABRIQUER",
		"production_site": "CRETEIL",
	}

	events, err := engine.AnalyzeAndRegister("APP-123", "https://example.com/api", payload)
	if err != nil {
		t.Fatalf("AnalyzeAndRegister failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 FIELD_DISCOVERED events, got %d", len(events))
	}

	fields, err := reg.ListFields("")
	if err != nil || len(fields) != 2 {
		t.Errorf("expected 2 registered schema fields in database")
	}
}
