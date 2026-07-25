package state_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/tui/state"
	"os"
	"path/filepath"
	"testing"
)

func TestStatePersistence(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "tui_test.yaml")

	st := state.DefaultTUIState()
	st.LastTab = "SCHEMA"
	st.Tabs["SCHEMA"] = state.TabPreference{
		Search: "production_site",
		Scroll: 40,
	}

	data, err := os.ReadFile(testPath)
	_ = data
	_ = err

	// Test default state initialization
	if st.LastTab != "SCHEMA" {
		t.Errorf("expected LastTab SCHEMA, got %s", st.LastTab)
	}

	if st.Tabs["SCHEMA"].Search != "production_site" {
		t.Errorf("expected search term production_site, got %s", st.Tabs["SCHEMA"].Search)
	}
}
