package crawler

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"strings"
)

// FormatEndpointSummary formats EndpointObservation for CLI display.
func FormatEndpointSummary(obs []domain.EndpointObservation) string {
	if len(obs) == 0 {
		return "No HTTP endpoints discovered yet."
	}

	var sb strings.Builder
	sb.WriteString("=== DISCOVERED ANEF ENDPOINTS ===\n\n")

	for i, o := range obs {
		sb.WriteString(fmt.Sprintf("%d. %s %s\n   Calls: %d | Last Status: HTTP %d | Latency: %dms\n   First Seen: %s | Last Seen: %s\n\n",
			i+1, o.Method, o.URL, o.Occurrences, o.LastStatusCode, o.LastLatencyMs,
			o.FirstSeen.Format("2006-01-02 15:04"), o.LastSeen.Format("2006-01-02 15:04")))
	}

	return sb.String()
}
