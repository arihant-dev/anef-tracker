CREATE INDEX IF NOT EXISTS idx_knowledge_nodes_type ON knowledge_nodes(type);
CREATE INDEX IF NOT EXISTS idx_knowledge_edges_from ON knowledge_edges(from_node);
CREATE INDEX IF NOT EXISTS idx_knowledge_edges_to ON knowledge_edges(to_node);
CREATE INDEX IF NOT EXISTS idx_knowledge_edges_type ON knowledge_edges(type);
CREATE INDEX IF NOT EXISTS idx_workflow_from_to ON workflow_transitions(from_status, to_status);
