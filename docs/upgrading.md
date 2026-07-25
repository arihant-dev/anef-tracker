# Upgrading & Migration Guide

This guide describes how to upgrade `anef-tracker` across version releases and migrate existing SQLite databases safely.

---

## 1. Upgrade Path Summary

| From Version | To Version | Schema Migration | Action Required |
|---|---|---|---|
| Phase 4 (`v0.4.0`) | Phase 5 (`v0.5.0`) | `012_notifications`, `013_watch` | Automatic on start |
| Phase 5 (`v0.5.0`) | Phase 6 (`v0.6.0`) | `014`–`017` (Profiles, Scoping, Audit) | Automatic on start |
| Phase 6 (`v0.6.0`) | `v0.9.0` / `v1.0.0` | None (Final release schema) | Binary replacement |

---

## 2. Standard Upgrade Procedure

### Step 1: Create Backup
Before upgrading binary versions, generate a compressed state backup:

```bash
anef backup create
```

### Step 2: Replace Binary
Replace the `anef` binary with the new version release:

```bash
# Example for macOS ARM64
cp bin/dist/anef-darwin-arm64 ~/.gopath/bin/anef
```

### Step 3: Run System Diagnostics
Execute system diagnostics to verify database migration status and integrity:

```bash
anef doctor
```

Output check:
```
✓ Database Connection : SQLite database healthy and migrated (017_audit_log)
✓ Evidence Integrity  : Verified
```

---

## 3. Rollback Procedure

If a rollback is ever required, restore from the backup archive:

```bash
anef backup restore backups/anef-2026-07-25_060000.tar.gz
```
