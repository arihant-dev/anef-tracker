# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| v0.9.x  | :white_check_mark: |
| < 0.9   | :x: |

---

## Security Boundaries & Guarantees

### Guaranteed
- **Immutable Evidence Storage**: Raw prefetoral HTTP responses are stored read-only with cryptographic SHA-256 payload hashes.
- **Row-Level Profile Isolation**: `ScopedRepository` structurally requires `context.Scope` parameters, preventing cross-applicant data leakage.
- **Tamper-Evident Audit Chain**: Operation audit log uses cryptographic `SHA-256` hash chaining to detect any unauthorized log edits.
- **Evidence-Preserving Redaction**: Redacted exports produce derived copies using `pkg/privacy` without mutating original evidence.
- **Zero Secondary PII Storage**: Ingestion-time PII scanner logs classification metadata only (`PIITypes`), never sensitive values.

### Not Currently Guaranteed
- **Database Encryption at Rest**: Default SQLite storage relies on OS filesystem file permissions (`0600`). Full SQLCipher database encryption is planned for optional CGO builds in v1.1.
- **Secure SSD Deletion**: Deleted records do not overwrite flash memory cells on modern SSDs.
- **Hardware-Backed Key Storage**: Keys are read from local environment or `.env` files rather than TPM or Secure Enclave hardware.

---

## Browser Authentication Trust Model

When using `anef login --browser`:
1. **Loopback Binding Only**: The temporary callback listener binds strictly to local loopback interface (`127.0.0.1:8484`) and is inaccessible over external network interfaces.
2. **Automated Shutdown**: The callback listener shuts down immediately upon receiving the session callback or timing out after 3 minutes.
3. **Direct Authentication**: Browser authentication occurs directly between your browser and official ANEF Keycloak servers (`https://administration-etrangers-en-france.interieur.gouv.fr`).
4. **Zero Third-Party Relays**: No credentials or tokens are transmitted to external servers, cloud proxies, or third-party analytics.

---

## Reporting a Vulnerability

If you discover a potential security vulnerability in `anef-tracker`, please report it privately:

- **Email**: `security@anef-tracker.org` (or via GitHub Private Security Advisories)
- **Response SLA**: We acknowledge receipt of vulnerability reports within 48 hours and aim to issue a patch release within 14 days for high-severity issues.
