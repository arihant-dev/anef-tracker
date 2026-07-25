# ADR-0002: Row-Level Profile Scoping

## Status
Accepted

## Context
ANEF Tracker evolved from single-user monitoring to multi-applicant profile management. Without structural enforcement, queries could accidentally leak data across profiles.

## Decision
We enforce **Row-Level Profile Scoping**:
1. All database tables (`snapshots`, `events`, `http_logs`, `evidence_records`) contain `profile_id` and `tracked_application_id` foreign keys.
2. The `ScopedRepository` interface in `pkg/storage/scoped.go` requires a `context.Scope` parameter on all query methods.
3. Unscoped queries are prohibited by the repository API signature.

## Consequences
- Total structural isolation between applicant profiles.
- Eliminates cross-profile data leakage risks at compile time.
