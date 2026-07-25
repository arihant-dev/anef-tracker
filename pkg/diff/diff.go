package diff

import (
	"encoding/json"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"reflect"

	"strings"
)

type FieldChange struct {
	Path     string          `json:"path"`
	Change   string          `json:"change"` // ADDED, REMOVED, MODIFIED
	OldVal   interface{}     `json:"old_val"`
	NewVal   interface{}     `json:"new_val"`
	Severity domain.Severity `json:"severity"`
}

type DiffResult struct {
	HasChanges bool            `json:"has_changes"`
	Severity   domain.Severity `json:"severity"`
	Changes    []FieldChange   `json:"changes"`
	Summary    string          `json:"summary"`
}

// CompareSnapshots performs a deep recursive JSON structural diff between old and new JSON responses.
func CompareSnapshots(oldJSON, newJSON []byte) (*DiffResult, error) {
	var oldMap, newMap map[string]interface{}

	if len(oldJSON) > 0 {
		_ = json.Unmarshal(oldJSON, &oldMap)
	}
	if len(newJSON) > 0 {
		_ = json.Unmarshal(newJSON, &newMap)
	}

	result := &DiffResult{
		Severity: domain.SeverityLow,
	}

	compareMaps("", oldMap, newMap, result)

	if len(result.Changes) > 0 {
		result.HasChanges = true
		result.Summary = formatHumanSummary(result.Changes, result.Severity)
	} else {
		result.Summary = "No changes detected between snapshots."
	}

	return result, nil
}

func compareMaps(prefix string, oldMap, newMap map[string]interface{}, res *DiffResult) {
	// Check removed or modified
	for k, oldV := range oldMap {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		newV, exists := newMap[k]
		if !exists {
			ch := FieldChange{
				Path:     path,
				Change:   "REMOVED",
				OldVal:   oldV,
				NewVal:   nil,
				Severity: classifySeverity(path, "REMOVED", oldV, nil),
			}
			res.Changes = append(res.Changes, ch)
			updateMaxSeverity(res, ch.Severity)
			continue
		}

		// Both exist, compare values
		if reflect.TypeOf(oldV) != reflect.TypeOf(newV) {
			ch := FieldChange{
				Path:     path,
				Change:   "MODIFIED",
				OldVal:   oldV,
				NewVal:   newV,
				Severity: classifySeverity(path, "MODIFIED", oldV, newV),
			}
			res.Changes = append(res.Changes, ch)
			updateMaxSeverity(res, ch.Severity)
		} else if om, ok := oldV.(map[string]interface{}); ok {
			nm, _ := newV.(map[string]interface{})
			compareMaps(path, om, nm, res)
		} else if !reflect.DeepEqual(oldV, newV) {
			ch := FieldChange{
				Path:     path,
				Change:   "MODIFIED",
				OldVal:   oldV,
				NewVal:   newV,
				Severity: classifySeverity(path, "MODIFIED", oldV, newV),
			}
			res.Changes = append(res.Changes, ch)
			updateMaxSeverity(res, ch.Severity)
		}
	}

	// Check added fields
	for k, newV := range newMap {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		if _, exists := oldMap[k]; !exists {
			ch := FieldChange{
				Path:     path,
				Change:   "ADDED",
				OldVal:   nil,
				NewVal:   newV,
				Severity: classifySeverity(path, "ADDED", nil, newV),
			}
			res.Changes = append(res.Changes, ch)
			updateMaxSeverity(res, ch.Severity)
		}
	}
}

func classifySeverity(path, _ string, _, newV interface{}) domain.Severity {
	lowerPath := strings.ToLower(path)

	if strings.Contains(lowerPath, "statut") || strings.Contains(lowerPath, "code_statut") {
		if fmt.Sprintf("%v", newV) == "TITRE_DISPONIBLE" || fmt.Sprintf("%v", newV) == "CONVOCATION_GENEREE" {
			return domain.SeverityCritical
		}
		return domain.SeverityHigh
	}

	if strings.Contains(lowerPath, "attestation") || strings.Contains(lowerPath, "site_retrait") {
		return domain.SeverityCritical
	}

	if strings.Contains(lowerPath, "version") || strings.Contains(lowerPath, "_version") {
		return domain.SeverityMedium
	}

	if strings.Contains(lowerPath, "updated") || strings.Contains(lowerPath, "timestamp") || strings.Contains(lowerPath, "date") {
		return domain.SeverityLow
	}

	return domain.SeverityMedium
}

func updateMaxSeverity(res *DiffResult, sev domain.Severity) {
	order := map[domain.Severity]int{
		domain.SeverityLow:      1,
		domain.SeverityMedium:   2,
		domain.SeverityHigh:     3,
		domain.SeverityCritical: 4,
	}
	if order[sev] > order[res.Severity] {
		res.Severity = sev
	}
}

func formatHumanSummary(changes []FieldChange, overallSev domain.Severity) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== SNAPSHOT DIFF SUMMARY [%s] ===\n", overallSev))

	for _, ch := range changes {
		switch ch.Change {
		case "MODIFIED":
			sb.WriteString(fmt.Sprintf("\n● %s changed:\n   %v\n   ↓\n   %v\n", ch.Path, ch.OldVal, ch.NewVal))
		case "ADDED":
			sb.WriteString(fmt.Sprintf("\n+ NEW FIELD %s = %v\n", ch.Path, ch.NewVal))
		case "REMOVED":
			sb.WriteString(fmt.Sprintf("\n- REMOVED FIELD %s (was %v)\n", ch.Path, ch.OldVal))
		}
	}

	return sb.String()
}
