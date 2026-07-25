# Developer Guide & Architecture Manual

Welcome to the `anef-tracker` developer guide. This document provides developer onboarding instructions, coding conventions, testing standards, and build procedures.

---

## 1. Project Layout Overview

```
anef-tracker/
├── cmd/anef/            # Main CLI entry point
├── internal/
│   ├── service/         # TrackerService application orchestration
│   └── version/         # Embedded build & version metadata
├── pkg/
│   ├── audit/           # Tamper-evident audit log & hash chain
│   ├── context/         # Profile & application scope management
│   ├── db/              # Pure-Go SQLite engine & migrations (001-017)
│   ├── domain/          # Domain models & ANEF status mappers
│   ├── evidence/        # SHA-256 payload verification & integrity
│   ├── export/          # Export generators & evidence ZIP bundle
│   ├── knowledge/       # Knowledge graph construction & validation
│   ├── notify/          # Desktop, Webhook, Email, Telegram alerts
│   ├── privacy/         # PII scanner, redactor & privacy policy
│   ├── profile/         # Profile & TrackedApplication CRUD
│   ├── storage/         # ScopedRepository row-level isolation
│   ├── timeline/        # Human timeline reconstruction engine
│   ├── tui/             # 17-tab Bubble Tea TUI with grouped navigation
│   ├── watch/           # Watcher daemon runner & scheduler
│   └── workflow/        # Empirical state machine duration engine
└── tests/acceptance/   # Integration & acceptance test suite
```

---

## 2. Development Setup

### Prerequisites
- Go ≥ 1.24
- SQLite 3.42+
- `make`

### Building from Source

```bash
git clone https://github.com/arihant-dev/anef-tracker.git
cd anef-tracker
make build
```

---

## 3. Running Unit & Integration Tests

Run the full test suite across all 35+ packages:

```bash
make test
```

Run tests with Go race detection:

```bash
make race
```

Run linters:

```bash
make lint
```

---

## 4. Coding Conventions

1. **Evidence Integrity**: Never mutate raw snapshot payloads or evidence records.
2. **Profile Isolation**: Always use `ScopedRepository` and pass `context.Scope` when querying database tables.
3. **No AI/ML**: Keep all logic 100% factual, explainable, and empirical.
