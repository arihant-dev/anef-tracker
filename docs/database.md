# SQLite Database Architecture & Schema

`anef-tracker` uses an embedded SQLite database (`data/anef.db`) managed via 11 incremental SQL migrations.

---

## Migration History

1. `001_initial.sql`: Core `applications` projection table.
2. `002_snapshots.sql`: Immutable snapshot metadata index.
3. `003_events.sql`: Domain transition event log.
4. `004_http_logs.sql`: Recorded HTTP traffic log.
5. `005_raw_payloads.sql`: Raw API response payload store.
6. `006_schema_registry.sql`: Observed API field schema registry.
7. `007_workflow_transitions.sql`: Historical state transition counts.
8. `008_knowledge_graph.sql`: Nodes and edges knowledge graph tables.
9. `009_provenance_indexes.sql`: Indexing for snapshot IDs and event references.
10. `010_evidence_integrity.sql`: Immutable evidence records and SHA-256 hashes.
11. `011_retention.sql`: Evidence retention policies (`retention_policy` table).

---

## Exporting Database Schema

To export the complete active schema DDL to `docs/database-schema.sql`:

```bash
anef db schema export
```
