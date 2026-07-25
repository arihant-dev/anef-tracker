# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.9.0] - 2026-07-25

### Added
- **Multi-Application Profile Vault**: Row-level profile isolation (`014`–`016` SQL migrations) with `context.Scope` structural enforcement.
- **Tamper-Evident Audit Log**: Cryptographic `SHA-256` hash chaining (`017_audit_log.sql`) with verification helper `anef audit verify`.
- **Privacy & PII Protection**: Deterministic PII scanner, evidence-preserving redactor, and `PrivacyObserver` ingestion hooks.
- **Scoped Evidence Bundles**: Application-scoped `.zip` export bundle containing manifest, integrity, timeline, and report.
- **17-Tab Grouped TUI**: Grouped navigation tabs (`Monitor`, `Evidence`, `Intelligence`, `Operations`) with `g` key selector.
- **CLI Commands**: `anef context`, `anef profile`, `anef privacy audit`, `anef audit`, `anef evidence bundle`.
- **Release Infrastructure**: Full CI/CD matrix, GoReleaser integration, Docker & Dev Container support.
