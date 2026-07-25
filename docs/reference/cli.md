# CLI Command Reference

Complete reference of all `anef-tracker` command line interface commands and subcommands.

---

## Command Matrix

| Command | Usage | Description |
|---|---|---|
| `anef login` | `anef login [--browser \| --curl "..." \| --user X --pass Y]` | Authenticate session via browser redirect, cURL, or Keycloak |
| `anef fetch` | `anef fetch` | Fetch current application status, record HTTP log & save raw snapshot |
| `anef status` | `anef status` | Display current application status code, label, & description |
| `anef timeline` | `anef timeline` | Render chronological progress timeline |
| `anef watch` | `anef watch [--interval 360m \| --once]` | Launch background watcher daemon or run single execution |
| `anef context` | `anef context` | Show active profile and application scope context |
| `anef profile` | `anef profile <list \| create <name> \| switch <id> \| delete --confirm <id>>` | Applicant profile management |
| `anef privacy` | `anef privacy audit` | Run PII scanner and audit policy report |
| `anef audit` | `anef audit <list \| verify>` | Display audit log or verify cryptographic SHA-256 hash chain |
| `anef evidence` | `anef evidence <search <term> \| verify \| export <id> \| bundle [--redact]>` | Forensic evidence commands and ZIP bundle export |
| `anef workflow` | `anef workflow [explain \| audit]` | Empirical state machine diagram and duration evidence |
| `anef analytics` | `anef analytics [explain <state>]` | Status duration medians & percentiles |
| `anef doctor` | `anef doctor` | Run 8-point system diagnostic health check |
| `anef version` | `anef version` | Display version, commit, build date, and platform info |
