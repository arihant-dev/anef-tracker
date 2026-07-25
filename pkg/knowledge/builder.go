package knowledge

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"time"
)

type GraphBuilder struct {
	DB *db.DB
}

func NewGraphBuilder(database *db.DB) *GraphBuilder {
	return &GraphBuilder{DB: database}
}

func (b *GraphBuilder) BuildFromDB() (*Graph, error) {
	g := NewGraph()

	if b.DB == nil {
		return g, nil
	}

	now := time.Now()

	// 1. Application Node
	appNode := Node{
		ID:    "app:dossier",
		Type:  NodeTypeApplication,
		Label: "ANEF Residence Permit Application",
		Evidence: []Provenance{
			{SourceType: SourceSnapshot, SnapshotID: "latest_snapshot", CreatedAt: now},
		},
	}
	g.AddNode(appNode)

	// 2. Status Nodes & Sequential Workflow Edges
	statuses := []string{
		"DEMANDE_SOUMISE", "DOSSIER_DEPOSE", "VALIDATION_FORMELLE",
		"INSTRUCTION_EN_COURS", "DECISION_VALIDEE", "TITRE_A_FABRIQUER", "TITRE_DISPONIBLE",
	}
	for i, st := range statuses {
		stNode := Node{
			ID:    "status:" + st,
			Type:  NodeTypeStatus,
			Label: st,
			Evidence: []Provenance{
				{SourceType: SourceSnapshot, SnapshotID: "baseline_snapshot", CreatedAt: now},
			},
		}
		g.AddNode(stNode)

		// Link Application to active status
		if st == "TITRE_A_FABRIQUER" {
			g.AddEdge(Edge{
				From: appNode.ID,
				To:   stNode.ID,
				Type: EdgeTypeHasStatus,
				Evidence: []Provenance{
					{SourceType: SourceSnapshot, SnapshotID: "live_snapshot", CreatedAt: now},
				},
			})
		}

		// Link sequential status workflow edges
		if i < len(statuses)-1 {
			nextStNodeID := "status:" + statuses[i+1]
			g.AddEdge(Edge{
				From: stNode.ID,
				To:   nextStNodeID,
				Type: EdgeTypeHasStatus,
				Evidence: []Provenance{
					{SourceType: SourceSnapshot, SnapshotID: "baseline_snapshot", CreatedAt: now},
				},
			})
		}
	}

	// 3. Endpoint Nodes
	epRows, err := b.DB.Conn.Query("SELECT DISTINCT method, url FROM http_logs")
	if err == nil {
		defer epRows.Close()
		idx := 1
		for epRows.Next() {
			var method, url string
			_ = epRows.Scan(&method, &url)
			epID := fmt.Sprintf("ep:%d", idx)
			idx++
			g.AddNode(Node{
				ID:    epID,
				Type:  NodeTypeEndpoint,
				Label: fmt.Sprintf("%s %s", method, url),
				Evidence: []Provenance{
					{SourceType: SourceHTTP, HTTPLogID: int64(idx), CreatedAt: now},
				},
			})
			g.AddEdge(Edge{
				From: epID,
				To:   appNode.ID,
				Type: EdgeTypeReturns,
				Evidence: []Provenance{
					{SourceType: SourceHTTP, HTTPLogID: int64(idx), CreatedAt: now},
				},
			})
		}
		_ = epRows.Err()
	}

	// 4. Events
	events, err := b.DB.GetEvents(50)
	if err == nil {
		for _, ev := range events {
			evID := fmt.Sprintf("event:%d", ev.ID)
			g.AddNode(Node{
				ID:    evID,
				Type:  NodeTypeEvent,
				Label: fmt.Sprintf("%s: %s", ev.Type, ev.FieldPath),
				Evidence: []Provenance{
					{SourceType: SourceEvent, EventID: ev.ID, CreatedAt: ev.Timestamp},
				},
			})
			g.AddEdge(Edge{
				From: appNode.ID,
				To:   evID,
				Type: EdgeTypeTriggeredEvent,
				Evidence: []Provenance{
					{SourceType: SourceEvent, EventID: ev.ID, CreatedAt: ev.Timestamp},
				},
			})
		}
	}

	return g, nil
}
