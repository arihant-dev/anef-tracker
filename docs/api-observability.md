# API Observability & Reverse Engineering Suite

`anef-tracker` provides reverse-engineering capabilities to explore ANEF REST endpoints and schema structures.

---

## Commands

- `anef api list`: Lists observed ANEF REST endpoints.
- `anef api inspect <id>`: Inspects headers, body payload, and structure of an endpoint call.
- `anef schema list`: Lists registered API schema fields and field types.
- `anef schema diff`: Computes structural JSON Schema diffs between snapshots.
- `anef replay <id>`: Replays recorded HTTP requests to verify API payload consistency.
- `anef endpoints`: Prints endpoint dependency graph in ASCII.
