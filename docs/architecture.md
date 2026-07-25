# ANEF Tracker Architecture

The **ANEF Residence Permit Tracker (`anef-tracker`)** is structured as an evidence-backed operational intelligence platform.

---

## Architectural Diagram

```
                 ANEF Web Portal
                       │
                       │ authenticated HTTP
                       ▼
              pkg/client & pkg/recorder
                       │
                       │ request/response traffic
                       ▼
               pkg/raw Payload Bytes
                       │
                       │ domain entity mapping
                       ▼
              pkg/domain Entities
                       │
         ┌─────────────┴─────────────┐
         ▼                           ▼
  Snapshot Store             SQLite Evidence Engine
  (pkg/snapshot)             (001-011 migrations)
         │                           │
         └─────────────┬─────────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
  Workflow Engine   Analytics   Evidence Graph
  (pkg/workflow)  (pkg/analytics)(pkg/knowledge)
         │             │             │
         └─────────────┼─────────────┘
                       │
                       ▼
             Bubble Tea TUI / CLI
```

---

## Component Boundaries

- `internal/service`: Orchestrates business workflows and status checks.
- `pkg/client`: Authenticated HTTP client with Keycloak token handling.
- `pkg/recorder`: Sanitized HTTP recording layer.
- `pkg/raw`: Raw API payload storage and JSON payload hashing.
- `pkg/domain`: Domain entities and status dictionary mapping (`status_mapping.yaml`).
- `pkg/storage` & `pkg/db`: Persistence interfaces and SQLite migration manager.
- `pkg/workflow`: Reconstructed state machine and transition evidence auditor.
- `pkg/analytics`: Duration percentile calculator ($N \ge 5$ threshold).
- `pkg/knowledge`: Multi-evidence knowledge graph and consistency validator.
- `pkg/evidence`: SHA-256 integrity verifier and audit engine.
- `pkg/backup`: Compressed tarball backup and restoration subsystem.
- `pkg/security`: Token security auditor and header credential redaction.
- `pkg/tui`: 12-tab Bubble Tea viewport user interface.
