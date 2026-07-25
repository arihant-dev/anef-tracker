# Security & Token Hardening Architecture

`anef-tracker` enforces strict security practices.

---

## Security Policies

1. **Session Storage**: Session tokens (`data/session.json`) are stored outside the Git repository with `0600` permissions.
2. **Database Permissions**: `data/anef.db` file permissions are restricted (non world-readable).
3. **Header Redaction**: Sensitive HTTP headers (`Authorization`, `Cookie`, `Set-Cookie`) are automatically redacted as `[REDACTED]` in recorded logs.

---

## Security Audit Command

Run automated security audit:

```bash
anef security audit
```
