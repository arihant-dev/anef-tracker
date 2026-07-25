# Compatibility Policy

This document outlines the compatibility matrix, operating system support, Go language version policy, component stability expectations, and database migration guarantees for `anef-tracker`.

---

## 1. Component Stability Policy

| Component | Stability Level | Policy & Guarantees |
|---|---|---|
| **CLI Interface** | Stable after `v1.0` | Backward compatible flag & command syntax after `v1.0.0` |
| **Configuration Schema** | Best effort | Documented YAML parameters; deprecations announced in CHANGELOG |
| **Database Schema** | Forward migrations supported | Schema migrations strictly forward-only (`001` through `017`) |
| **Internal Go Packages** | Unstable before `v1.0` | `pkg/` packages subject to internal refactoring before `v1.0.0` |
| **Plugin Interfaces** | Stable after `v1.1` | Plugin contracts (`pkg/plugin`) frozen at `v1.1.0` |

---

## 2. Supported Go Versions

| Go Version | Status | Notes |
|---|---|---|
| Go 1.24.x | **Supported** | Primary release & CI target |
| Go 1.25.x | **Supported** | CI matrix target |

---

## 3. Operating System Support

| OS | Architecture | Status | Support Level |
|---|---|---|---|
| macOS (Darwin) | `arm64` (Apple Silicon), `amd64` (Intel) | **Fully Supported** | Primary TUI & Desktop notification platform |
| Linux | `amd64`, `arm64` | **Fully Supported** | Desktop notifications via `notify-send` |
| Windows | `amd64` | **Supported** | CLI & headless watcher mode |

---

## 4. Database Engine & Migration Policy

- **Database Engine**: Pure Go SQLite (`github.com/glebarez/go-sqlite`) with **CGO disabled**.
- **SQLite Version**: SQLite ≥ 3.42 compatibility.
- **Migration Guarantees**:
  - Migrations are strictly **forward-only** (`001_initial.sql` through `017_audit_log.sql`).
  - Schema migrations run automatically upon application startup via `db.InitDB()`.
  - Database schema changes preserve backward compatibility for historical evidence records.
