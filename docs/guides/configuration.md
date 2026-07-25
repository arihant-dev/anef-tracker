# Configuration Manual

This guide details all configuration parameters for `anef-tracker`. Configuration is loaded from `config.yaml` or environment variables.

---

## Configuration Parameter Matrix

| Section | Parameter | Env Variable | Default | Description |
|---|---|---|---|---|
| **Database** | `database.path` | `DATABASE_PATH` | `data/anef.db` | Path to SQLite database file |
| **Storage** | `storage.immutable_snapshots` | `STORAGE_IMMUTABLE_SNAPSHOTS` | `true` | Snapshots stored read-only |
| **Notifications** | `notifications.enabled` | `NOTIFICATIONS_ENABLED` | `false` | Master notification toggle |
| | `notifications.desktop` | `NOTIFICATIONS_DESKTOP` | `true` | OS desktop popups |
| | `notifications.webhook_url` | `NOTIFICATIONS_WEBHOOK_URL` | `""` | Slack/Discord webhook URL |
| **Retention** | `retention.http_logs_days` | `RETENTION_HTTP_LOGS_DAYS` | `365` | HTTP traffic retention window |
| **Security** | `security.encryption_enabled` | `SECURITY_ENCRYPTION_ENABLED` | `false` | Encryption flag |
| | `security.key_source` | `SECURITY_KEY_SOURCE` | `env` | Key source (`env`, `keychain`) |
