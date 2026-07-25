package tui_test

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/tui"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	m := tui.InitialModel(nil)
	rendered := m.View()

	if !strings.Contains(rendered, "ANEF WORKFLOW INTELLIGENCE PLATFORM") {
		t.Errorf("expected header title in TUI view")
	}

	if !strings.Contains(rendered, "1: Overview") {
		t.Errorf("expected tabs header in TUI view")
	}
}

func TestQuitKey(t *testing.T) {
	m := tui.InitialModel(nil)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	updatedModel, cmd := m.Update(keyMsg)
	_ = updatedModel

	if cmd == nil {
		t.Fatalf("expected tea.Quit command on 'q' key press")
	}
}

func TestOverviewDisplaysStatus(t *testing.T) {
	app := &domain.Application{
		NumeroDemande: "9999999999999999999",
		ForeignerID:   "9999999999",
		LegalCategory: "E1",
		Status: domain.ApplicationStatus{
			Code:        "TITRE_A_FABRIQUER",
			Label:       "Residence Permit in Production",
			Description: "The decision has been approved.",
			Severity:    domain.SeverityHigh,
		},
		Version: 56,
	}

	m := tui.ModelWithApplication(app)
	rendered := m.View()

	if !strings.Contains(rendered, "TITRE_A_FABRIQUER") {
		t.Errorf("expected status code TITRE_A_FABRIQUER in TUI output")
	}

	if !strings.Contains(rendered, "Residence Permit in Production") {
		t.Errorf("expected status label in TUI output")
	}
}

func TestAPITabRendering(t *testing.T) {
	m := tui.InitialModel(nil)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}}

	updatedModel, _ := m.Update(keyMsg)
	rendered := updatedModel.View()

	if !strings.Contains(rendered, "6: API Explorer") {
		t.Errorf("expected Tab 6 header to be rendered")
	}
}

func TestSchemaTabRendering(t *testing.T) {
	m := tui.InitialModel(nil)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'7'}}

	updatedModel, _ := m.Update(keyMsg)
	rendered := updatedModel.View()

	if !strings.Contains(rendered, "7: Schema") {
		t.Errorf("expected Tab 7 header to be rendered")
	}
}

func TestTabSwitchResetsViewport(t *testing.T) {
	m := tui.InitialModel(nil)

	// Switch to Tab 7 (Schema, index 6)
	m.SwitchTab(6)

	schemaTS := m.GetActiveTabState()
	var testLines []string
	for i := 1; i <= 100; i++ {
		testLines = append(testLines, fmt.Sprintf("Field %d", i))
	}
	schemaTS.Viewport.SetContent(testLines)
	schemaTS.Viewport.ScrollDown(40)

	if schemaTS.Viewport.Offset != 40 {
		t.Fatalf("expected schema viewport offset 40, got %d", schemaTS.Viewport.Offset)
	}

	// Switch to Tab 4 (Events, index 3)
	m.SwitchTab(3)
	eventsTS := m.GetActiveTabState()

	if eventsTS.Viewport.Offset != 0 {
		t.Fatalf("expected events tab viewport offset to be 0 upon tab switch, got %d", eventsTS.Viewport.Offset)
	}
}
