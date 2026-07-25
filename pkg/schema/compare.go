package schema

import (
	"fmt"
	"strings"
)

type SchemaDiffResult struct {
	Endpoint      string   `json:"endpoint"`
	AddedFields   []string `json:"added_fields"`
	RemovedFields []string `json:"removed_fields"`
	TypeChanges   []string `json:"type_changes"`
	HasDiff       bool     `json:"has_diff"`
	Summary       string   `json:"summary"`
}

// CompareSchemaDocuments compares two SchemaDocuments to identify field structure shifts.
func CompareSchemaDocuments(oldSchema, newSchema *SchemaDocument) *SchemaDiffResult {
	res := &SchemaDiffResult{
		Endpoint: newSchema.Endpoint,
	}

	if oldSchema == nil || oldSchema.Fields == nil {
		for f := range newSchema.Fields {
			res.AddedFields = append(res.AddedFields, f)
		}
		res.HasDiff = len(res.AddedFields) > 0
		res.Summary = fmt.Sprintf("Initial schema registered with %d fields.", len(res.AddedFields))
		return res
	}

	// Check added fields or type changes
	for f, newSpec := range newSchema.Fields {
		oldSpec, exists := oldSchema.Fields[f]
		if !exists {
			res.AddedFields = append(res.AddedFields, fmt.Sprintf("%s (%s)", f, newSpec.Type))
		} else if oldSpec.Type != newSpec.Type {
			res.TypeChanges = append(res.TypeChanges, fmt.Sprintf("%s: %s → %s", f, oldSpec.Type, newSpec.Type))
		}
	}

	// Check removed fields
	for f := range oldSchema.Fields {
		if _, exists := newSchema.Fields[f]; !exists {
			res.RemovedFields = append(res.RemovedFields, f)
		}
	}

	res.HasDiff = len(res.AddedFields) > 0 || len(res.RemovedFields) > 0 || len(res.TypeChanges) > 0

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== SCHEMA DIFF FOR %s ===\n", newSchema.Endpoint))
	if !res.HasDiff {
		sb.WriteString("No schema structural changes detected.\n")
	} else {
		for _, a := range res.AddedFields {
			sb.WriteString(fmt.Sprintf("+ ADDED FIELD: %s\n", a))
		}
		for _, r := range res.RemovedFields {
			sb.WriteString(fmt.Sprintf("- REMOVED FIELD: %s\n", r))
		}
		for _, tc := range res.TypeChanges {
			sb.WriteString(fmt.Sprintf("~ TYPE MODIFIED: %s\n", tc))
		}
	}
	res.Summary = sb.String()

	return res
}
