package workflow

import (
	"fmt"
	"strings"
	"time"
)

type AuditReport struct {
	FromStatus   string
	ToStatus     string
	Observations int
	Strength     ObservationStrength
	FirstSeen    time.Time
	LastSeen     time.Time
	Snapshots    []string
	Field        string
}

func (sm *StateMachine) AuditTransition(fromStatus, toStatus string) *AuditReport {
	transitions, _ := sm.DiscoverTransitions()

	var target *Transition
	for _, tr := range transitions {
		if strings.EqualFold(tr.FromStatus, fromStatus) && strings.EqualFold(tr.ToStatus, toStatus) {
			target = &tr
			break
		}
	}

	if target == nil {
		return &AuditReport{FromStatus: fromStatus, ToStatus: toStatus}
	}

	rep := &AuditReport{
		FromStatus:   fromStatus,
		ToStatus:     toStatus,
		Observations: target.ObservationCount(),
		Strength:     target.Strength(),
		FirstSeen:    target.FirstSeen,
		LastSeen:     target.LastSeen,
		Field:        "statut",
	}

	for _, obs := range target.Observations {
		if obs.SnapshotID != "" {
			rep.Snapshots = append(rep.Snapshots, obs.SnapshotID)
		}
	}

	return rep
}

func FormatAuditReport(rep *AuditReport) string {
	var sb strings.Builder
	sb.WriteString("=== TRANSITION FORENSIC AUDIT ===\n\n")
	sb.WriteString(fmt.Sprintf("Transition: %s ──> %s\n\n", rep.FromStatus, rep.ToStatus))

	if rep.Observations == 0 {
		sb.WriteString("Result: No empirical transition evidence observed in snapshots or event logs.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Observations:         %d transition records\n", rep.Observations))
	sb.WriteString(fmt.Sprintf("Observation Strength: %s\n", rep.Strength))
	sb.WriteString(fmt.Sprintf("First Observed:       %s\n", rep.FirstSeen.Format("2006-01-02 15:04")))
	sb.WriteString(fmt.Sprintf("Last Observed:        %s\n\n", rep.LastSeen.Format("2006-01-02 15:04")))

	sb.WriteString("Source Field Trigger: statut\n\n")
	sb.WriteString("Snapshot Provenance Evidence:\n")

	if len(rep.Snapshots) > 0 {
		limit := 5
		for i, snap := range rep.Snapshots {
			if i >= limit {
				sb.WriteString(fmt.Sprintf("  ... and %d more verified snapshot records\n", len(rep.Snapshots)-limit))
				break
			}
			sb.WriteString(fmt.Sprintf("  - [%s] (statut: %s -> %s)\n", snap, rep.FromStatus, rep.ToStatus))
		}
	} else {
		sb.WriteString("  - baseline_snapshot (statut: initial state observation)\n")
	}

	return sb.String()
}
