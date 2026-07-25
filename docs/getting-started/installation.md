# Installation Manual

`anef-tracker` supports four primary installation methods across macOS, Linux, and Windows.

---

## 1. Go Install (Recommended for Go Developers)
```bash
go install anef_tracker/cmd/anef@latest
```

---

## 2. Pre-Built Binaries
Download official cross-platform binary archives from [GitHub Releases](https://github.com/anef-tracker/anef-tracker/releases):
- macOS (Apple Silicon `arm64` & Intel `amd64`)
- Linux (`amd64` & `arm64`)
- Windows (`amd64`)

---

## 3. Docker Container
Run containerized headless watcher or CLI:
```bash
docker run -d --name anef-watcher -v $(pwd)/data:/app/data anef-tracker:v0.9.0
```

---

## 4. Build from Source
```bash
git clone https://github.com/anef-tracker/anef-tracker.git
cd anef-tracker
make build
```
