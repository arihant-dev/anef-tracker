# ADR-0001: Immutable Evidence Model

## Status
Accepted

## Context
When tracking administrative residence permit applications on ANEF, users require undeniable proof of status changes, submitted documents, and prefetoral decisions. Data storage must be tamper-proof and verifiable.

## Decision
We establish an **Immutable Evidence Model**:
1. Every raw HTTP response payload is saved to disk accompanied by an immutable SHA-256 payload hash.
2. Evidence records link snapshots, events, and HTTP traffic to their underlying raw payload bytes.
3. Original evidence records are strictly read-only and never overwritten or mutated.

## Consequences
- High evidentiary integrity for legal and administrative disputes.
- Storage requirement grows linearly with snapshot count.
