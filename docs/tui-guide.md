# Terminal User Interface (TUI) User Guide

Launch the interactive 12-tab terminal user interface via:

```bash
anef tui
```

---

## Interactive View Tabs

1. `1: Overview`: Current application position & status summary.
2. `2: Timeline`: Milestone progression status.
3. `3: Diff`: Structural snapshot diff view.
4. `4: Events`: Filterable application state transition log.
5. `5: Docs`: Uploaded document attachments.
6. `6: API Explorer`: Observed ANEF REST endpoints.
7. `7: Schema`: Registered field dictionary & search.
8. `8: Replay`: Traffic replay interface.
9. `9: Logs`: Recorded HTTP traffic log.
10. `10: Workflow`: Reconstructed state machine diagram.
11. `11: Analytics`: State duration distribution percentiles.
12. `12: Evidence Graph`: Interactive node & edge graph view.

---

## Keyboard Shortcuts

- `Tab` / `Shift+Tab`: Navigate tabs.
- `1` to `9`: Jump directly to tab index.
- `↑ / ↓ / j / k`: Scroll viewport content.
- `Ctrl+P`: Launch Global Command Palette fuzzy finder.
- `Ctrl+E`: Export active view.
- `Ctrl+D`: Toggle SQLite Database Status modal overlay.
- `Ctrl+I`: Toggle Evidence Integrity Verification modal overlay.
- `v`: Validate graph consistency (Tab 12).
- `/`: Search view content (Tab 7 & 12).
- `?`: Help overlay.
- `q`: Quit.
