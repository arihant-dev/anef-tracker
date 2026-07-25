# ANEF Tracker Configuration Manual

This document details all configuration parameters for `anef-tracker`. Configuration is read from `config.yaml` or environment variables.

---

## 1. Database Configuration (`database`)

| Setting | YAML Key | Env Variable | Default | Description |
|---|---|---|---|---|
| Database Path | `database.path` | `DATABASE_PATH` | `data/anef.db` | Absolute or relative path to SQLite database file |

---

## 2. Storage & Snapshot Immutability (`storage`)

| Setting | YAML Key | Env Variable | Default | Description |
|---|---|---|---|---|
| Immutable Snapshots | `storage.immutable_snapshots` | `STORAGE_IMMUTABLE_SNAPSHOTS` | `true` | When true, raw HTTP payload snapshots are read-only |

---

## 3. Notifications (`notifications`)

| Setting | YAML Key | Env Variable | Default | Description |
|---|---|---|---|---|
| Enabled | `notifications.enabled` | `NOTIFICATIONS_ENABLED` | `false` | Global notification master switch |
| Desktop Alerts | `notifications.desktop` | `NOTIFICATIONS_DESKTOP` | `true` | Local OS desktop popups (macOS osascript, Linux notify-send) |
| Webhook URL | `notifications.webhook_url` | `NOTIFICATIONS_WEBHOOK_URL` | `""` | Slack, Discord, or Telegram JSON webhook endpoint |

---

## 4. Retention Policy (`retention`)

| Setting | YAML Key | Env Variable | Default | Description |
|---|---|---|---|---|
| HTTP Log Days | `retention.http_logs_days` | `RETENTION_HTTP_LOGS_DAYS` | `365` | Retention window for recorded HTTP request/response traffic |

---

## 5. Security & Capability Status (`security`)

| Setting | YAML Key | Env Variable | Default | Description |
|---|---|---|---|---|
| Encryption Enabled | `security.encryption_enabled` | `SECURITY_ENCRYPTION_ENABLED` | `false` | Encryption flag (Pure-Go SQLite uses OS-level 0600 file protection) |
| Key Source | `security.key_source` | `SECURITY_KEY_SOURCE` | `env` | Source of encryption keys (`env`, `keychain`) |
