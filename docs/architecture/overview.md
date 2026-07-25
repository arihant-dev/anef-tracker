# System Architecture Overview

`anef-tracker` is designed as a **private, multi-application evidence vault** and workflow intelligence platform.

---

## High-Level Architecture

```
                       Bubble Tea TUI
                             │
                      Profile Context
                             │
            ┌────────────────┼────────────────┐
            │                │                │
       Privacy Layer    Audit Layer    Security Layer
            │                │                │
            └────────────────┼────────────────┘
                             │
                  Scoped Repository Layer
                  (enforces profile_id)
                             │
            ┌────────────────┼────────────────┐
            │                                 │
      Application Repository          Evidence Repository
      (profile_id + app_id)           (profile_id + app_id)
            │                                 │
            └────────────────┼────────────────┘
                             │
                  SQLite Evidence Vault
```

---

## Architectural Principles

1. **Evidence Immutability**: Raw payloads and HTTP records are saved read-only with cryptographic SHA-256 hashes.
2. **Row-Level Profile Scoping**: `ScopedRepository` structurally requires `context.Scope` parameters on all query methods.
3. **Fact-Based Workflow Reconstruction**: empirical state machine modeling without AI/ML inference.
4. **Tamper-Evident Audit Chains**: `SHA-256` hash chaining across audit log entries.
