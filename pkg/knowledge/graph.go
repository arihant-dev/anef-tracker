package knowledge

import (
	"fmt"
	"strings"
	"time"
)

type Graph struct {
	Nodes map[string]Node `json:"nodes" yaml:"nodes"`
	Edges []Edge          `json:"edges" yaml:"edges"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]Node),
		Edges: []Edge{},
	}
}

func (g *Graph) AddNode(n Node) {
	if n.ID != "" {
		if len(n.Evidence) == 0 {
			n.Evidence = append(n.Evidence, Provenance{SourceType: SourceSnapshot, CreatedAt: time.Now()})
		}
		g.Nodes[n.ID] = n
	}
}

func (g *Graph) AddEdge(e Edge) {
	if e.From != "" && e.To != "" {
		if len(e.Evidence) == 0 {
			e.Evidence = append(e.Evidence, Provenance{SourceType: SourceSnapshot, CreatedAt: time.Now()})
		}
		g.Edges = append(g.Edges, e)
	}
}

func (g *Graph) FormatASCII() string {
	var sb strings.Builder
	sb.WriteString("=== ANEF EVIDENCE GRAPH ===\n\n")
	sb.WriteString(fmt.Sprintf("Total Observed Nodes: %d | Total Observed Edges: %d\n\n", len(g.Nodes), len(g.Edges)))

	for id, node := range g.Nodes {
		p := node.PrimaryProvenance()
		sb.WriteString(fmt.Sprintf("[%s] %s (%s) — %d evidence records\n", node.Type, node.Label, id, len(node.Evidence)))
		if p.SnapshotID != "" {
			sb.WriteString(fmt.Sprintf("   Primary Snapshot: %s [%s]\n", p.SnapshotID, p.SourceType))
		}
		for _, edge := range g.Edges {
			if edge.From == id {
				toNode, ok := g.Nodes[edge.To]
				toLabel := edge.To
				if ok {
					toLabel = toNode.Label
				}
				sb.WriteString(fmt.Sprintf("   └── %s ──> [%s] %s\n", edge.Type, edge.To, toLabel))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (g *Graph) ExplainNode(nodeID string) string {
	node, ok := g.Nodes[nodeID]
	if !ok {
		return fmt.Sprintf("Node '%s' not found in Evidence Graph.", nodeID)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== EVIDENCE GRAPH NODE EXPLANATION: %s ===\n\n", nodeID))
	sb.WriteString(fmt.Sprintf("Node Label:     %s\n", node.Label))
	sb.WriteString(fmt.Sprintf("Node Type:      %s\n", node.Type))
	sb.WriteString(fmt.Sprintf("Evidence Count: %d observations\n", len(node.Evidence)))

	sb.WriteString("\nEvidence Provenance Records:\n")
	for i, ev := range node.Evidence {
		sb.WriteString(fmt.Sprintf("  %d. [%s] Source: %s", i+1, ev.SourceType, ev.CreatedAt.Format("2006-01-02 15:04")))
		if ev.SnapshotID != "" {
			sb.WriteString(fmt.Sprintf(" | Snapshot: %s", ev.SnapshotID))
		}
		if ev.EventID > 0 {
			sb.WriteString(fmt.Sprintf(" | Event ID: #%d", ev.EventID))
		}
		if ev.HTTPLogID > 0 {
			sb.WriteString(fmt.Sprintf(" | HTTP Log ID: #%d", ev.HTTPLogID))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nConnected Edges:\n")
	outboundCount := 0
	for _, edge := range g.Edges {
		if edge.From == nodeID {
			outboundCount++
			targetNode := g.Nodes[edge.To]
			sb.WriteString(fmt.Sprintf("  └── Outbound: %s ──> [%s] %s\n", edge.Type, edge.To, targetNode.Label))
		} else if edge.To == nodeID {
			outboundCount++
			sourceNode := g.Nodes[edge.From]
			sb.WriteString(fmt.Sprintf("  ├── Inbound:  %s <── [%s] %s\n", edge.Type, edge.From, sourceNode.Label))
		}
	}
	if outboundCount == 0 {
		sb.WriteString("  (No direct connections recorded)\n")
	}

	return sb.String()
}

func (g *Graph) SearchEvidence(term string) string {
	termLower := strings.ToLower(term)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== FORENSIC EVIDENCE SEARCH: '%s' ===\n\n", term))

	matchedNodes := 0
	for id, node := range g.Nodes {
		if strings.Contains(strings.ToLower(id), termLower) ||
			strings.Contains(strings.ToLower(node.Label), termLower) ||
			strings.Contains(strings.ToLower(string(node.Type)), termLower) {
			matchedNodes++
			p := node.PrimaryProvenance()
			sb.WriteString(fmt.Sprintf("%d. [%s] %s (%s)\n", matchedNodes, node.Type, node.Label, id))
			sb.WriteString(fmt.Sprintf("   Evidence Observations: %d records\n", len(node.Evidence)))
			sb.WriteString(fmt.Sprintf("   Primary Provenance:    [%s] %s\n", p.SourceType, p.SnapshotID))
			sb.WriteString("\n")
		}
	}

	if matchedNodes == 0 {
		sb.WriteString(fmt.Sprintf("No evidence nodes or relationships matched query term '%s'.\n", term))
	}

	return sb.String()
}
