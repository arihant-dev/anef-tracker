package knowledge

import (
	"fmt"
	"strings"
)

type ValidationReport struct {
	TotalNodes      int
	TotalEdges      int
	OrphanNodes     []string
	BrokenEdges     []string
	MissingEvidence []string
	Valid           bool
}

func (g *Graph) ValidateGraph() *ValidationReport {
	rep := &ValidationReport{
		TotalNodes: len(g.Nodes),
		TotalEdges: len(g.Edges),
		Valid:      true,
	}

	// 1. Check for Orphan Nodes and Missing Evidence
	for id, node := range g.Nodes {
		if len(node.Evidence) == 0 {
			rep.Valid = false
			rep.MissingEvidence = append(rep.MissingEvidence, fmt.Sprintf("Node '%s' has no source evidence attached", id))
		}

		connected := false
		for _, edge := range g.Edges {
			if edge.From == id || edge.To == id {
				connected = true
				break
			}
		}
		if !connected && node.Type != NodeTypeApplication {
			rep.Valid = false
			rep.OrphanNodes = append(rep.OrphanNodes, id)
		}
	}

	// 2. Check for Broken Edges
	for _, edge := range g.Edges {
		if _, ok := g.Nodes[edge.From]; !ok {
			rep.Valid = false
			rep.BrokenEdges = append(rep.BrokenEdges, fmt.Sprintf("Edge from non-existent node '%s'", edge.From))
		}
		if _, ok := g.Nodes[edge.To]; !ok {
			rep.Valid = false
			rep.BrokenEdges = append(rep.BrokenEdges, fmt.Sprintf("Edge to non-existent node '%s'", edge.To))
		}
	}

	return rep
}

func FormatValidationReport(rep *ValidationReport) string {
	var sb strings.Builder
	sb.WriteString("=== EVIDENCE GRAPH CONSISTENCY VALIDATION ===\n\n")

	sb.WriteString(fmt.Sprintf("Nodes checked: %d\n", rep.TotalNodes))
	sb.WriteString(fmt.Sprintf("Edges checked: %d\n\n", rep.TotalEdges))

	symbol := func(ok bool) string {
		if ok {
			return "✓"
		}
		return "✗"
	}

	sb.WriteString(fmt.Sprintf("%s No orphan nodes\n", symbol(len(rep.OrphanNodes) == 0)))
	sb.WriteString(fmt.Sprintf("%s All edges resolve cleanly\n", symbol(len(rep.BrokenEdges) == 0)))
	sb.WriteString(fmt.Sprintf("%s Provenance evidence attached to all nodes\n\n", symbol(len(rep.MissingEvidence) == 0)))

	if !rep.Valid {
		sb.WriteString("Discrepancies Identified:\n")
		for _, o := range rep.OrphanNodes {
			sb.WriteString(fmt.Sprintf("  - Orphan Node: %s\n", o))
		}
		for _, b := range rep.BrokenEdges {
			sb.WriteString(fmt.Sprintf("  - Broken Edge: %s\n", b))
		}
		for _, m := range rep.MissingEvidence {
			sb.WriteString(fmt.Sprintf("  - Missing Evidence: %s\n", m))
		}
		sb.WriteString("\nGraph Status: DISCREPANCIES DETECTED\n")
	} else {
		sb.WriteString("Graph Status: CONSISTENT & FULLY VALIDATED\n")
	}

	return sb.String()
}
