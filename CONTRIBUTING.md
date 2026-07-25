# Contributing to ANEF Tracker

Thank you for your interest in contributing to `anef-tracker`! We welcome bug reports, documentation updates, and pull requests.

---

## 1. Development Process

1. **Fork the Repository**: Create a personal fork on GitHub.
2. **Clone and Setup**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/anef-tracker.git
   cd anef-tracker
   make build
   ```
3. **Branching**: Create a feature branch named `feature/your-feature` or `fix/your-fix`.
4. **Code Quality**: Ensure all code passes formatting, vetting, linters, and race detection:
   ```bash
   make lint
   make race
   ```
5. **Commit Messages**: Use clean, descriptive commit messages (e.g. `pkg/notify: add retry logic for webhook delivery`).

---

## 2. Invariants & Rules

- **No AI / Inference Models**: All intelligence logic must remain 100% empirical, deterministic, and evidence-backed.
- **Evidence Immutability**: Original raw payloads and evidence records must never be mutated or overwritten.
- **Profile Isolation**: Always use `ScopedRepository` and include `context.Scope` on queries.
