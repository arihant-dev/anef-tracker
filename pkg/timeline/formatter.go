package timeline

import (
	"fmt"
	"strings"
)

func (ht *HumanTimeline) FormatASCII() string {
	var sb strings.Builder
	sb.WriteString("=== RECONSTRUCTED APPLICATION HUMAN TIMELINE ===\n\n")

	sb.WriteString(fmt.Sprintf("Application Number: %s\n", ht.ApplicationID))
	sb.WriteString(fmt.Sprintf("Application Started: %s\n", ht.StartedAt.Format("02 January 2006")))
	sb.WriteString(fmt.Sprintf("Current State:       %s\n", ht.CurrentStatus))
	sb.WriteString(fmt.Sprintf("Elapsed Duration:    %d days\n", ht.ElapsedDays))
	sb.WriteString(fmt.Sprintf("Total Milestones:    %d milestones\n\n", len(ht.Milestones)))

	for i, m := range ht.Milestones {
		symbol := "[ ]"
		if m.IsCurrent {
			symbol = "[★]"
		} else if m.Completed {
			symbol = "[✓]"
		}

		sb.WriteString(fmt.Sprintf("%s %s (%s) — %s\n",
			symbol, m.Date.Format("2006-01-02"), m.StatusLabel, m.StatusCode))

		if m.IsCurrent {
			sb.WriteString(fmt.Sprintf("    ├── Current Position (Days in state: %d days)\n", m.DaysInState))
			sb.WriteString(fmt.Sprintf("    └── Evidence: snapshot %s\n", m.SnapshotID))
		} else if m.Completed {
			sb.WriteString(fmt.Sprintf("    └── Evidence: %d snapshot records\n", m.EvidenceCount))
		}

		if i < len(ht.Milestones)-1 {
			sb.WriteString("    │\n")
		}
	}

	return sb.String()
}
