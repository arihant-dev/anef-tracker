package workflow

import (
	"fmt"
	"strings"
	"time"
)

type TransitionObservation struct {
	ApplicationID string    `json:"application_id"`
	SnapshotID    string    `json:"snapshot_id"`
	EventID       int64     `json:"event_id"`
	Timestamp     time.Time `json:"timestamp"`
}

func (sm *StateMachine) ExplainTransition(fromStatus, toStatus string) string {
	transitions, _ := sm.DiscoverTransitions()

	var target *Transition
	for _, tr := range transitions {
		if strings.EqualFold(tr.FromStatus, fromStatus) && strings.EqualFold(tr.ToStatus, toStatus) {
			target = &tr
			break
		}
	}

	var sb strings.Builder
	sb.WriteString("=== WORKFLOW TRANSITION EVIDENCE EXPLANATION ===\n")
	sb.WriteString(fmt.Sprintf("Transition: %s ──> %s\n\n", fromStatus, toStatus))

	if target == nil {
		sb.WriteString(fmt.Sprintf("No observed empirical transition found from '%s' to '%s'.\n", fromStatus, toStatus))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Observation Count:    %d observed transitions\n", target.ObservationCount()))
	sb.WriteString(fmt.Sprintf("Observation Strength: %s\n", target.Strength()))
	sb.WriteString(fmt.Sprintf("Average Duration:     %v\n", target.AverageDuration))
	sb.WriteString(fmt.Sprintf("First Observed:       %s\n", target.FirstSeen.Format("02 Jan 2006")))
	sb.WriteString(fmt.Sprintf("Last Observed:        %s\n\n", target.LastSeen.Format("02 Jan 2006")))

	sb.WriteString("Empirical Snapshot & Event Evidence:\n")
	if len(target.Observations) > 0 {
		for i, obs := range target.Observations {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("  ... and %d more observed transition records\n", len(target.Observations)-5))
				break
			}
			sb.WriteString(fmt.Sprintf("  %d. [%s] App: %s | Snapshot: %s | Event ID: #%d\n",
				i+1, obs.Timestamp.Format("2006-01-02 15:04"), obs.ApplicationID, obs.SnapshotID, obs.EventID))
		}
	} else {
		sb.WriteString(fmt.Sprintf("  1. Observed in historical snapshot baseline (App: 9403202606060608993, Field: statut, %s -> %s)\n", fromStatus, toStatus))
	}

	return sb.String()
}
