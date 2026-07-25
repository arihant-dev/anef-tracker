package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"time"
)

func EventsToCSV(events []domain.Event) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	_ = w.Write([]string{"ID", "Timestamp", "Type", "Severity", "Confidence", "FieldPath", "OldVal", "NewVal"})

	for _, ev := range events {
		_ = w.Write([]string{
			fmt.Sprintf("%d", ev.ID),
			ev.Timestamp.Format(time.RFC3339),
			ev.Type,
			string(ev.Severity),
			fmt.Sprintf("%.2f", ev.Confidence),
			ev.FieldPath,
			ev.OldVal,
			ev.NewVal,
		})
	}

	w.Flush()
	return buf.Bytes(), nil
}

func FieldsToCSV(fields []domain.FieldObservation) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	_ = w.Write([]string{"ID", "Endpoint", "Path", "Type", "Occurrences", "Confidence", "FirstSeen"})

	for _, f := range fields {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			f.Endpoint,
			f.Path,
			f.Type,
			fmt.Sprintf("%d", f.Occurrences),
			fmt.Sprintf("%.2f", f.Confidence),
			f.FirstSeen.Format(time.RFC3339),
		})
	}

	w.Flush()
	return buf.Bytes(), nil
}
