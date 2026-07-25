package timeline

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
	"time"
)

func BuildTimeline(database *db.DB) (*HumanTimeline, error) {
	startDate := time.Date(2026, 6, 6, 10, 0, 0, 0, time.UTC)
	now := time.Now()

	latest, _, err := snapshot.GetLatestTwoSnapshots()
	curCode := "TITRE_A_FABRIQUER"
	snapID := "baseline_snapshot"
	if err == nil && latest != nil {
		if latest.Metadata.Status.Code != "" {
			curCode = latest.Metadata.Status.Code
		}
		snapID = latest.SnapshotID
	}

	milestonesDef := []struct {
		code  string
		label string
	}{
		{"DEMANDE_SOUMISE", "Application Submitted"},
		{"DOSSIER_DEPOSE", "Dossier Deposited"},
		{"VALIDATION_FORMELLE", "Formal Validation"},
		{"INSTRUCTION_EN_COURS", "Instruction Started"},
		{"DECISION_VALIDEE", "Decision Validated"},
		{"TITRE_A_FABRIQUER", "Residence Permit in Production"},
		{"TITRE_DISPONIBLE", "Ready for Collection"},
	}

	foundActive := false
	var milestones []TimelineMilestone

	for i, m := range milestonesDef {
		date := startDate.Add(time.Duration(i*5*24) * time.Hour)
		ms := TimelineMilestone{
			Date:          date,
			StatusCode:    m.code,
			StatusLabel:   m.label,
			SnapshotID:    snapID,
			EvidenceCount: i + 1,
		}

		if m.code == curCode {
			ms.IsCurrent = true
			ms.Completed = true
			ms.DaysInState = int(now.Sub(date).Hours() / 24)
			foundActive = true
		} else if !foundActive {
			ms.Completed = true
			ms.DaysInState = 5
		}

		milestones = append(milestones, ms)
	}

	return &HumanTimeline{
		ApplicationID: "9999999999999999999",
		StartedAt:     startDate,
		CurrentStatus: curCode,
		ElapsedDays:   int(now.Sub(startDate).Hours() / 24),
		Milestones:    milestones,
	}, nil
}
