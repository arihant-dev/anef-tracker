package knowledge_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"path/filepath"
	"testing"
)

func TestGraphBuilder(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_graph.db")
	database, err := db.InitDBWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed initializing test db: %v", err)
	}
	defer database.Conn.Close()

	builder := knowledge.NewGraphBuilder(database)
	g, err := builder.BuildFromDB()
	if err != nil {
		t.Fatalf("BuildFromDB failed: %v", err)
	}

	if len(g.Nodes) < 8 {
		t.Errorf("expected at least 8 nodes in knowledge graph, got %d", len(g.Nodes))
	}

	repo := knowledge.NewRepository(database)
	if err := repo.SaveGraph(g); err != nil {
		t.Fatalf("SaveGraph failed: %v", err)
	}
}
