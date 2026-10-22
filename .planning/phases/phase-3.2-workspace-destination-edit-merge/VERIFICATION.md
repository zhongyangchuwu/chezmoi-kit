# Verification: Workspace Destination Edit and Merge

**Date:** 2026-10-23  
**Result:** Local pass; awaiting PR CI.

## Behavior

- `DestinationEditCommand` prefers `VISUAL`, falls back to `EDITOR`, accepts quoted command arguments, invokes the absolute destination path directly, and reports an actionable missing-editor error.
- Isolated chezmoi integration changes a destination through `E`'s app command and proves the mapped source file remains unchanged.
- `MergeCommand` delegates exactly to `chezmoi merge <absolute-target>` for an eligible managed regular file.
- Isolated chezmoi integration configures `merge.command`, records its three arguments, and observes destination, source, and nonempty rendered-target paths without implicit source/destination change.
- TUI tests exercise `E` and `m` key routing, terminal ownership/input lock, original-scope refresh, contextual footer visibility, unsupported symlink/script gates, and refresh after a nonzero merge return.

## Gates

| Check | Result |
|---|---|
| `go test ./internal/app ./internal/tui` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go vet ./...` | pass |
| `go mod tidy -diff` | pass |
| `git diff --check` | pass |
| Go LSP diagnostics after reload | pass |
