package crawler

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"strings"
)

type EndpointNode struct {
	Method   string   `json:"method"`
	URL      string   `json:"url"`
	Returns  []string `json:"returns"`
	Requires []string `json:"requires"`
}

type EndpointGraph struct {
	Nodes []EndpointNode `json:"nodes"`
}

// BuildEndpointGraph constructs hypermedia relationships between observed endpoints.
func BuildEndpointGraph(observations []domain.EndpointObservation) *EndpointGraph {
	graph := &EndpointGraph{}

	for _, obs := range observations {
		node := EndpointNode{
			Method: obs.Method,
			URL:    obs.URL,
		}

		if strings.Contains(obs.URL, "demande_sejour") {
			node.Returns = []string{"Application", "Status", "AttestationDepot"}
		} else if strings.Contains(obs.URL, "documents") {
			node.Returns = []string{"Justificatifs", "PDF"}
			node.Requires = []string{"numero_demande"}
		} else if strings.Contains(obs.URL, "profil") || strings.Contains(obs.URL, "moi") {
			node.Returns = []string{"UserProfile", "ForeignerID"}
		}

		graph.Nodes = append(graph.Nodes, node)
	}

	return graph
}

func (g *EndpointGraph) FormatGraphASCII() string {
	var sb strings.Builder
	sb.WriteString("=== ENDPOINT RELATIONSHIP GRAPH ===\n\n")

	for i, node := range g.Nodes {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, node.Method, node.URL))
		if len(node.Returns) > 0 {
			sb.WriteString(fmt.Sprintf("   Returns:  %s\n", strings.Join(node.Returns, ", ")))
		}
		if len(node.Requires) > 0 {
			sb.WriteString(fmt.Sprintf("   Requires: %s\n", strings.Join(node.Requires, ", ")))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
