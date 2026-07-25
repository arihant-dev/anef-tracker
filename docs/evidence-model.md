# Evidence Model & Provenance Tracking

`anef-tracker` relies on a fully deterministic evidence model.

---

## Source Attribution (`SourceType`)

Every node, edge, and transition observation carries source provenance:

- `SNAPSHOT`: Sourced from an immutable JSON snapshot file.
- `EVENT`: Sourced from a domain transition event.
- `HTTP_LOG`: Sourced from an authenticated HTTP request/response log.
- `SCHEMA`: Sourced from API schema discovery.

---

## Multi-Evidence Slice Architecture

Nodes and edges carry `Evidence []Provenance` slices, allowing elements to accumulate historical evidence records across multiple fetch executions over time.

---

## Hash Verification

Stored payloads are hashed via SHA-256 (`pkg/evidence/hash.go`). Execute integrity verification via:

```bash
anef evidence verify
```
