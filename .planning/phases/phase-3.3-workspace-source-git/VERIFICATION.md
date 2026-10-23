# Verification: Workspace Source Git Handoff

**Date:** 2026-10-23  
**Result:** Local pass; awaiting PR CI.

## Behavior

- `SourceGitCommand` resolves `chezmoi source-path` and runs lazygit in that directory; direct `cm git` reuses the same command builder.
- Isolated integration places a disposable lazygit executable on `PATH` and verifies it runs in the chezmoi source directory.
- Workspace inventory records source Git porcelain working-tree changes. A source Git status error is visible without preventing workspace inventory.
- TUI coverage exercises `g`, terminal handoff/return refresh, source Git status line, contextual footer/help, and existing stale-preview/context restoration behavior.

## Gates

| Check | Result |
|---|---|
| Focused app/TUI tests | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go vet ./...` | pass |
| `go mod tidy -diff` | pass |
| `git diff --check` | pass |
| Go LSP diagnostics after reload | pass |
