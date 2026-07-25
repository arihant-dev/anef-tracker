# Plugin System Architecture Specification (v1.1+ Target)

This document defines the interface specifications for the community plugin ecosystem planned for `v1.1+`.

---

## Plugin Interfaces

### 1. Notification Provider Plugin (`pkg/plugin/notifier`)

```go
type NotifierPlugin interface {
    Name() string
    Configure(config map[string]string) error
    Send(event domain.Event) error
}
```

### 2. Evidence Exporter Plugin (`pkg/plugin/exporter`)

```go
type ExporterPlugin interface {
    Name() string
    Format() string
    Export(scope appcontext.Scope, data *export.ExportData) ([]byte, error)
}
```

### 3. Storage Backend Plugin (`pkg/plugin/storage`)

```go
type StoragePlugin interface {
    Name() string
    Init(connStr string) error
    SaveSnapshot(scope appcontext.Scope, snapshot *snapshot.SnapshotRef) error
}
```

---

## Design Goals
- Decouples core evidence tracking from third-party notification or cloud export integrations.
- Keeps core `anef-tracker` lightweight and 100% pure-Go.
