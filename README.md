# ANEF Residence Permit Tracker & Evidence-Based Workflow Intelligence Platform

[![Continuous Integration](https://github.com/arihant-dev/anef-tracker/actions/workflows/test.yml/badge.svg)](https://github.com/arihant-dev/anef-tracker/actions/workflows/test.yml)
[![Security Audit](https://github.com/arihant-dev/anef-tracker/actions/workflows/security.yml/badge.svg)](https://github.com/arihant-dev/anef-tracker/actions/workflows/security.yml)
[![CodeQL Analysis](https://github.com/arihant-dev/anef-tracker/actions/workflows/codeql.yml/badge.svg)](https://github.com/arihant-dev/anef-tracker/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/arihant-dev/anef-tracker)](https://goreportcard.com/report/github.com/arihant-dev/anef-tracker)
[![Release](https://img.shields.io/github/v/release/arihant-dev/anef-tracker)](https://github.com/arihant-dev/anef-tracker/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`anef-tracker` is a **private, multi-application evidence vault** and CLI/TUI operations assistant for tracking French residence permit applications (*Administration Numérique pour les Étrangers en France* — ANEF).

It transforms raw prefetoral API responses into **tamper-evident, explainable operational intelligence**, reconstructing administrative state machines, calculating duration percentiles, verifying snapshot integrity, and managing applicant profiles with row-level data isolation.

---

## Installation Methods

### 1. Via `go install` (Go 1.24+)
```bash
go install github.com/arihant-dev/anef-tracker/cmd/anef@latest
```

### 2. Pre-built Release Binaries
Download standalone cross-platform archives for macOS (`arm64`/`amd64`), Linux, or Windows directly from [GitHub Releases](https://github.com/arihant-dev/anef-tracker/releases).

### 3. Docker Container
Run containerized headless watcher or CLI:
```bash
docker run -d --name anef-watcher -v $(pwd)/data:/app/data anef-tracker:v0.9.0
```

### 4. Build from Source
```bash
git clone https://github.com/arihant-dev/anef-tracker.git
cd anef-tracker
make build
```

### Initial Authentication & Tracking

1. Authenticate via browser-assisted cURL import:
   ```bash
   anef login
   ```
   Or paste a copied DevTools cURL directly:
   ```bash
   anef login --curl "curl 'https://...'"
   ```

2. Validate operational readiness before fetching:
   ```bash
   anef session validate
   ```

3. Fetch current status & record evidence:
   ```bash
   anef fetch
   ```

4. View active context & application status:
   ```bash
   anef context
   anef status
   ```

5. Launch interactive 17-tab TUI:
   ```bash
   anef tui
   ```

---

## Architecture at a Glance

```text
Browser Authentication / DevTools cURL
                  │
                  ▼
          Importer Subsystem
                  │
                  ▼
         Provider-Agnostic Session
                  │
                  ▼
          HTTP Log Recorder
                  │
                  ▼
      Immutable Evidence Store (SQLite / JSON)
                  │
                  ▼
      State Machine & Duration Analytics Engine
                  │
                  ▼
      Privacy Observer & Redaction Layer
                  │
                  ▼
      Cryptographic Hash-Chained Audit Log
                  │
                  ▼
     CLI / 17-Tab TUI / Evidence Bundles
```

---

## Command Reference

| Category | Command | Description |
|---|---|---|
| **Authentication & Session** | `anef login` | Browser-assisted or direct cURL session import |
| | `anef session inspect` | View tokens, claims, cookies & import metadata |
| | `anef session validate` | Check operational readiness to fetch (`Ready for fetch`) |
| | `anef session doctor` | Diagnostic audit of session vault permissions & encryption |
| **Monitoring** | `anef status` | Print current application status & legal classification |
| | `anef timeline` | Display chronological human progress timeline |
| | `anef watch` | Run watcher daemon (`--interval 360m` or `--once`) |
| | `anef context` | Display active profile & application scope context |
| | `anef profile` | `list`, `create`, `switch`, or `delete` applicant profiles |
| | `anef tui` | Launch 17-tab grouped terminal interface |
| **Evidence & Vault** | `anef privacy audit` | Run PII scanner & privacy policy audit |
| | `anef audit list` | View tamper-evident audit log |
| | `anef audit verify` | Verify audit log hash chain integrity |
| | `anef evidence bundle` | Export application-scoped evidence ZIP (`--redact` supported) |
| | `anef evidence search` | Keyword search across evidence graph nodes & events |
| | `anef evidence verify` | Verify snapshot SHA-256 hashes & link integrity |
| **Workflow** | `anef workflow` | Render empirical state machine diagram |
| | `anef analytics` | Display status duration medians & percentiles |
| **System** | `anef doctor` | Run system installation diagnostics & health checks |
| | `anef backup` | `create` compressed state archive or `restore` state |
| | `anef security audit`| Audit local permissions (0600) & credential protection |

---

## Documentation

- [Architecture Overview](docs/architecture.md)
- [Configuration Manual](docs/configuration.md)
- [Compatibility Policy](docs/compatibility.md)
- [Upgrading & Migration Guide](docs/upgrading.md)
- [Developer Guide](docs/development.md)
- [Architecture Decision Records (ADRs)](docs/adr/)

---

## Security Model

| Guarantee | Status |
|---|---|
| Evidence immutability | ✅ |
| Row-level profile isolation | ✅ |
| Tamper-evident audit log | ✅ |
| Evidence-preserving redaction | ✅ |
| Secret & PII scanning in CI | ✅ |
| CodeQL SAST security analysis | ✅ |
| SBOM generation (SPDX/CycloneDX) | ✅ |
| At-rest database encryption | Planned |
| Secure hardware key storage | Planned |

---

## License & Security

- **License**: Released under the [MIT License](LICENSE).
- **Security Policy**: See [SECURITY.md](SECURITY.md) for security design philosophy and disclosure procedures.
