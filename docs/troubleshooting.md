# System Diagnostics & Troubleshooting Guide

If operational issues occur, execute system diagnostics:

```bash
anef doctor
```

---

## Common Issues & Solutions

### 1. Token Expired (HTTP 401)
- **Symptom**: `anef fetch` returns HTTP 401 Unauthorized.
- **Fix**: Re-authenticate via `./bin/anef login --curl "curl ..."` or `--user <num> --pass <pwd>`.

### 2. Database Locks
- **Symptom**: SQLite database busy error.
- **Fix**: Run database maintenance via `anef db vacuum`.

### 3. Verification Failures
- **Symptom**: `anef evidence verify` reports invalid hashes.
- **Fix**: Check storage permissions and restore from latest valid backup via `anef backup restore <path>`.
