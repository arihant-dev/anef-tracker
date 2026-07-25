package workflow_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/workflow"
	"path/filepath"
	"testing"
)

func TestTransitionDiscovery(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_workflow.db")
	database, err := db.InitDBWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed initializing test db: %v", err)
	}

	sm := workflow.NewStateMachine(database)
	transitions, err := sm.DiscoverTransitions()
	if err != nil {
		t.Fatalf("DiscoverTransitions failed: %v", err)
	}

	if len(transitions) == 0 {
		t.Errorf("expected transitions to be discovered")
	}
}

func TestWorkflowGraph(t *testing.T) {
	sm := workflow.NewStateMachine(nil)
	ascii := sm.RenderASCII("TITRE_A_FABRIQUER")

	if len(ascii) == 0 {
		t.Errorf("expected non-empty ASCII representation")
	}
}
