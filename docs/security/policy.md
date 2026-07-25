# Security Policy

See root [SECURITY.md](../../SECURITY.md) for full vulnerability reporting procedures and security design principles.

---

## Security Highlights

- **Local Storage File Permissions**: SQLite database files created with `0600` permissions.
- **Browser Auth Trust Model**: Local callback listener binds strictly to `127.0.0.1:8484` and shuts down automatically after 3 minutes.
- **Tamper-Evident Audit Chains**: `SHA-256` hash chaining across operation audit log entries.
