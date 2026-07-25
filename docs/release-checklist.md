# Release Checklist

Use this checklist before cutting any release tag (e.g. `v0.9.0`, `v1.0.0`).

---

## 1. Quality & Hardening Verification

- [ ] Clean build passes: `make build`
- [ ] Linters pass cleanly: `make lint`
- [ ] All unit tests pass with zero data races: `make race`
- [ ] Technical documentation complete: `make docs-check`
- [ ] Database diagnostics pass: `anef doctor`
- [ ] Audit log hash chain intact: `anef audit verify`
- [ ] Secret audit complete (`gitleaks` pass / no tracked `.db` or `.env` files)

---

## 2. Release & Tagging Execution

- [ ] Update [CHANGELOG.md](../CHANGELOG.md) with release highlights.
- [ ] Verify GoReleaser dry-run:
  ```bash
  goreleaser release --snapshot --clean
  ```
- [ ] Tag the git commit:
  ```bash
  git tag -a v0.9.0 -m "Release v0.9.0"
  git push origin v0.9.0
  ```
- [ ] Verify GitHub Actions release workflow succeeds and publishes binaries.
