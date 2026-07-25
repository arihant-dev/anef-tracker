# ANEF Residence Permit Tracker & Evidence Vault

Welcome to the documentation portal for `anef-tracker` (`v0.9.0-beta`).

`anef-tracker` is a **private, multi-application evidence vault** and CLI/TUI operations assistant for tracking French residence permit applications (*Administration Numérique pour les Étrangers en France* — ANEF).

---

## Core Capabilities

- **Immutable Evidence Storage**: Raw prefetoral HTTP responses backed by cryptographic SHA-256 payload hashes.
- **Row-Level Profile Scoping**: Structural database isolation preventing cross-applicant data leakage (`pkg/storage/scoped.go`).
- **Fact-Based Workflow Reconstruction**: Empirical state machine duration medians and progress percentiles — **zero AI/ML inference**.
- **17-Tab Grouped Bubble Tea TUI**: Terminal user interface with grouped tab navigation (`Monitor`, `Evidence`, `Intelligence`, `Operations`).
- **Tamper-Evident Audit Logging**: Cryptographically chained operation log (`SHA-256`) detecting any unauthorized log edits.

---

## Quick Navigation

- [Quickstart Guide](getting-started/quickstart.md)
- [Installation Manual](getting-started/installation.md)
- [System Architecture](architecture/overview.md)
- [CLI Reference](reference/cli.md)
- [Security Policy](security/policy.md)
