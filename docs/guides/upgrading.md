# Upgrading & Migration Guide

This guide describes how to upgrade `anef-tracker` releases and migrate existing SQLite databases safely.

---

## Upgrade Steps

1. **Backup State**:
   ```bash
   anef backup create
   ```
2. **Replace Binary**: Download new release binary or run `go install anef_tracker/cmd/anef@latest`.
3. **Run System Diagnostics**:
   ```bash
   anef doctor
   ```
4. **Verify Audit Chain**:
   ```bash
   anef audit verify
   ```
