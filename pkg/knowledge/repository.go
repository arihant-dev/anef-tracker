package knowledge

import (
	"encoding/json"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{DB: database}
}

func (r *Repository) SaveGraph(g *Graph) error {
	if r.DB == nil || g == nil {
		return nil
	}

	for id, node := range g.Nodes {
		metaBytes, _ := json.Marshal(node.Metadata)
		_, err := r.DB.Conn.Exec(
			"INSERT OR REPLACE INTO knowledge_nodes (id, type, label, metadata) VALUES (?, ?, ?, ?)",
			id, string(node.Type), node.Label, string(metaBytes),
		)
		if err != nil {
			return fmt.Errorf("failed saving node %s: %w", id, err)
		}
	}

	for _, edge := range g.Edges {
		_, err := r.DB.Conn.Exec(
			"INSERT INTO knowledge_edges (from_node, to_node, type) VALUES (?, ?, ?)",
			edge.From, edge.To, string(edge.Type),
		)
		if err != nil {
			return fmt.Errorf("failed saving edge %s->%s: %w", edge.From, edge.To, err)
		}
	}

	return nil
}
