CREATE TABLE IF NOT EXISTS knowledge_nodes (
    id TEXT PRIMARY KEY,
    type TEXT,
    label TEXT,
    metadata TEXT
);

CREATE TABLE IF NOT EXISTS knowledge_edges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_node TEXT,
    to_node TEXT,
    type TEXT
);
