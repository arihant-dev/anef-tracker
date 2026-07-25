package knowledge_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"testing"
)

func TestGraphValidation(t *testing.T) {
	g := knowledge.NewGraph()
	n1 := knowledge.Node{ID: "app:dossier", Type: knowledge.NodeTypeApplication, Label: "App"}
	n2 := knowledge.Node{ID: "status:TITRE", Type: knowledge.NodeTypeStatus, Label: "Status"}
	g.AddNode(n1)
	g.AddNode(n2)

	g.AddEdge(knowledge.Edge{From: n1.ID, To: n2.ID, Type: knowledge.EdgeTypeHasStatus})

	rep := g.ValidateGraph()
	if !rep.Valid {
		t.Errorf("expected clean graph to be valid")
	}
}

func TestOrphanDetection(t *testing.T) {
	g := knowledge.NewGraph()
	g.AddNode(knowledge.Node{ID: "orphan:1", Type: knowledge.NodeTypeField, Label: "Orphan"})

	rep := g.ValidateGraph()
	if rep.Valid {
		t.Errorf("expected graph with orphan node to fail validation")
	}

	if len(rep.OrphanNodes) != 1 || rep.OrphanNodes[0] != "orphan:1" {
		t.Errorf("expected orphan:1 in orphan nodes report")
	}
}
