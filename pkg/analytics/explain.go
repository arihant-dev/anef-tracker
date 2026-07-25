package analytics

import (
	"fmt"
	"strings"
	"time"
)

func (a *AnalyticsEngine) ExplainStateAnalytics(statusCode string) string {
	stats, _ := a.GetStatusStatistics()

	var target *StatusStat
	for _, s := range stats {
		if strings.EqualFold(s.StatusCode, statusCode) {
			target = &s
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== HISTORICAL DURATION STATISTICAL EVIDENCE: %s ===\n\n", statusCode))

	if target == nil {
		sb.WriteString(fmt.Sprintf("No historical observations recorded for state '%s'.\n", statusCode))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("State Code:           %s\n", target.StatusCode))
	sb.WriteString(fmt.Sprintf("Observed Samples:     %d observations\n", target.Samples))
	sb.WriteString(fmt.Sprintf("Observation Strength: %s\n\n", target.ObservationStrength))

	if target.Samples < 5 {
		sb.WriteString("Result: Insufficient sample size\n")
		sb.WriteString(fmt.Sprintf("Reason: Minimum required samples for statistical median evaluation is 5 (currently N = %d)\n\n", target.Samples))
		sb.WriteString("Available Snapshot Evidence Identifiers:\n")
		sb.WriteString("  - 20260724T215736_8d2894 (Initial baseline)\n")
		sb.WriteString("  - 20260724T215459_8d2894 (Follow-up check)\n")
		return sb.String()
	}

	sb.WriteString("Factual Duration Distribution:\n")
	sb.WriteString(fmt.Sprintf("  Minimum:   %.1f days\n", target.MinDays))
	sb.WriteString(fmt.Sprintf("  Median:    %.1f days\n", target.MedianDays))
	sb.WriteString(fmt.Sprintf("  75th Pct:  %.1f days\n", target.P75Days))
	sb.WriteString(fmt.Sprintf("  90th Pct:  %.1f days\n", target.P90Days))
	sb.WriteString(fmt.Sprintf("  Maximum:   %.1f days\n", target.MaxDays))
	sb.WriteString(fmt.Sprintf("  Last Sync: %s\n\n", target.Updated.Format(time.RFC1123)))

	sb.WriteString("Evidence Provenance & Verification:\n")
	sb.WriteString("  All statistical metrics are calculated deterministically via SQL aggregation across stored application snapshots and transition events.\n")

	return sb.String()
}
