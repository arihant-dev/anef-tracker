# ADR-0005: Evidence-Preserving Redaction

## Status
Accepted

## Context
Users need to share evidence reports with third parties (lawyers, prefectures, advisors) without disclosing sensitive PII (phone numbers, full names, addresses).

## Decision
We enforce **Evidence-Preserving Redaction**:
1. Original raw payloads and database evidence records are **never mutated or overwritten**.
2. Redacted export bundles produce derived copies using `pkg/privacy`.
3. Ingestion-time privacy scanning logs classification metadata only (`PIITypes`), never sensitive values.

## Consequences
- Preserves the evidentiary hash of raw data while enabling safe third-party export sharing.
- Prevents secondary copies of sensitive PII.
