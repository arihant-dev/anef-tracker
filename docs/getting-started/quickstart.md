# Quickstart Guide

This guide gets you up and running with `anef-tracker` in less than 5 minutes.

---

## Step 1: Install Binary

```bash
go install github.com/arihant-dev/anef-tracker/cmd/anef@latest
```

Verify installation:

```bash
anef version
```

---

## Step 2: Authenticate Session

Copy a cURL request from Chrome DevTools while logged into the ANEF portal:

```bash
anef login --curl "curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/...'"
```

---

## Step 3: Fetch Status & Record Evidence

```bash
anef fetch
```

---

## Step 4: Inspect Active Context & Launch TUI

```bash
anef context
anef tui
```
