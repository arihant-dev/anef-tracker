# Evidence Export Bundles Guide

`anef-tracker` can export self-contained `.zip` archives containing verified evidence claims and reports.

---

## Generating Evidence Bundles

Generate an application-scoped evidence export bundle:

```bash
anef evidence bundle --profile 1 --application 1
```

Generate a **redacted** bundle (masks PII values while preserving original database evidence):

```bash
anef evidence bundle --redact
```

---

## Bundle Contents

```
anef-evidence-bundle-p1-a1-{timestamp}.zip
├── manifest.json       ← Metadata, scope, database hash
├── integrity.json      ← SHA-256 hashes of included files
├── timeline.md         ← Reconstructed human timeline
├── report.md           ← Evidence-backed status report
```
