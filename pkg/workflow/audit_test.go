package workflow_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/workflow"
	"testing"
)

func TestTransitionAudit(t *testing.T) {
	sm := workflow.NewStateMachine(nil)
	rep := sm.AuditTransition("INSTRUCTION_EN_COURS", "DECISION_VALIDEE")

	if rep == nil {
		t.Fatalf("expected non-nil AuditReport")
	}

	if rep.Observations == 0 {
		t.Errorf("expected > 0 observations for default transition")
	}

	summary := workflow.FormatAuditReport(rep)
	if len(summary) == 0 {
		t.Errorf("expected non-empty formatted audit report string")
	}
}
