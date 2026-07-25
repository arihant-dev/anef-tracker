package export

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"strings"
	"time"
)

func EventsToMarkdown(events []domain.Event) string {
	var sb strings.Builder
	sb.WriteString("# Recorded Application Events\n\n")
	sb.WriteString("| ID | Timestamp | Type | Severity | Field | Old Value | New Value |\n")
	sb.WriteString("|---|---|---|---|---|---|---|\n")

	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s | `%s` | `%s` | `%s` |\n",
			ev.ID, ev.Timestamp.Format(time.RFC3339), ev.Type, ev.Severity, ev.FieldPath, ev.OldVal, ev.NewVal))
	}

	return sb.String()
}
