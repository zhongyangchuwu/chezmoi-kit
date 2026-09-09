# Verification: Workspace Source Edit

**Date:** 2026-09-10  
**Result:** Pass — locally verified; awaiting PR CI.

## Behavior Evidence

- `SourceEditCommand` accepts eligible managed files and symlinks with source mappings, including template/encrypted/uninspected entries, and rejects directories, scripts, ignored/unmanaged, unmapped, remove, external, and unknown entries.
- The command is exactly `chezmoi edit --apply=false --watch=false <absolute-target>` through the existing terminal command/process boundary.
- `TestServiceSourceEditCommandDoesNotApplyDestination` configures a disposable editor, changes source, and observes a dirty review while destination remains `before\n`.
- `TestServiceSourceEditCommandEditsTemplateAndEncryptedSourcesWithoutApply` changes real template and age-encrypted sources, confirms rendered source is updated and encrypted source remains ciphertext, and confirms neither destination auto-applies.
- TUI coverage exercises `e` eligibility, contextual footer, input locking during handoff, original-scope inventory refresh, surviving-selection fallback, editor-error refresh, refresh failure preservation, and stale pre-refresh preview rejection using a refresh epoch.

## Actual PTY

A disposable chezmoi wrapper and editor ran `/tmp/cm-source-edit ui` with one managed file. Pressing `e` returned through the Bubble Tea terminal handoff; the TUI reported `editor closed · workspace refreshed · destination unchanged`. Observed files after return:

```text
source/app.conf:      edited through cm ui
destination/app.conf: before source edit
```

## Quality Gates

| Check | Result |
|---|---|
| `go test ./internal/app -run TestServiceSourceEditCommand -count=1` | pass |
| `go test ./internal/cli ./internal/app ./internal/tui` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go vet ./...` | pass |
| `go mod tidy -diff` | pass |
| `git diff --check` | pass |
| Go LSP diagnostics for `internal/app/**/*.go` and `internal/tui/**/*.go` | pass |

## Scope Check

Implemented only source editing and reusable handoff/refresh support. No local destination editor, merge, lazygit, workspace apply, action menu, or auto-apply/watch behavior was added.
