# ADR-0003: Hash-Chained Audit Log

## Status
Accepted

## Context
Administrative tools must log configuration changes, profile switches, and export operations. The log must be tamper-evident.

## Decision
We implement a **Hash-Chained Audit Log**:
1. Each audit log entry computes a cryptographic SHA-256 hash incorporating the entry data and the preceding entry's hash (`SHA256(action:resource:profileID:metadata:prevHash:timestamp)`).
2. The verification engine `VerifyAuditChain(db)` checks for any broken hash linkages.

## Consequences
- Retroactive modifications or deletions in the audit log are instantly detected.
- Fulfills forensic audit requirements.
