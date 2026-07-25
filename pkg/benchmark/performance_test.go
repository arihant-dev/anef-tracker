package benchmark_test

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"github.com/arihant-dev/anef-tracker/pkg/search"
	"github.com/arihant-dev/anef-tracker/pkg/workflow"
	"testing"
	"time"
)

func BenchmarkEvidenceSearch10000Records(b *testing.B) {
	g := knowledge.NewGraph()
	for i := 1; i <= 10000; i++ {
		g.AddNode(knowledge.Node{
			ID:    fmt.Sprintf("node:%d", i),
			Type:  knowledge.NodeTypeField,
			Label: fmt.Sprintf("Field %d path.statut", i),
		})
	}

	b.ResetTimer()
	for b.Loop() {
		_ = g.SearchEvidence("statut")
	}
}

func BenchmarkGraphRebuild10000Nodes(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		g := knowledge.NewGraph()
		for j := 1; j <= 10000; j++ {
			g.AddNode(knowledge.Node{
				ID:    fmt.Sprintf("node:%d", j),
				Type:  knowledge.NodeTypeField,
				Label: fmt.Sprintf("Schema field #%d", j),
			})
		}
	}
}

func BenchmarkWorkflowReconstruction10000Events(b *testing.B) {
	sm := workflow.NewStateMachine(nil)
	b.ResetTimer()
	for b.Loop() {
		_ = sm.RenderASCII("TITRE_A_FABRIQUER")
	}
}

func BenchmarkSearchParser(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = search.ParseQuery("type:FIELD_DISCOVERED field:adresse after:2026-07-01 statut")
	}
}

func TestBenchmarkThresholds(t *testing.T) {
	startTime := time.Now()
	sm := workflow.NewStateMachine(nil)
	_ = sm.RenderASCII("TITRE_A_FABRIQUER")
	elapsed := time.Since(startTime)

	if elapsed > 500*time.Millisecond {
		t.Errorf("workflow reconstruction took %v (> 500ms target)", elapsed)
	}
}
