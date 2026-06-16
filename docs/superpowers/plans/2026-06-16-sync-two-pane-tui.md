# Sync Two Pane TUI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the current review-mode sync TUI into a lazygit-style two-pane interface while keeping the existing reconciliation contract.

**Architecture:** Keep Bubble Tea. Split the UI by responsibility: root program/update plumbing, model helpers, keys/styles, diff cache/rendering, and pane layout. The left pane owns file selection and pending action markers; the right pane shows either the selected file diff or confirm summary.

**Tech Stack:** Go, Bubble Tea v2, Bubbles help/key, Lipgloss, existing `internal/reconcile` review service.

---

## File Structure

- Modify `internal/ui/tui.go`: keep `RunSyncTUI`, `Update`, command plumbing.
- Create `internal/ui/model.go`: model fields, constructor, focus/mode/pending helpers.
- Create `internal/ui/keys.go`: key map and styles.
- Create `internal/ui/diff.go`: diff cache, diff loading, diff scrolling, line coloring.
- Create `internal/ui/view.go`: two-pane layout and render helpers.
- Create `internal/ui/confirm.go`: confirm rendering and execution command.
- Modify `internal/ui/sync_test.go`: keep behavior-level model tests; add cache/order/scroll contract only if needed.

## Tasks

### Task 1: Split current TUI file without behavior change

Move types/functions into the new files listed above. Run `go test ./internal/ui` after the split. Expected: pass.

### Task 2: Add focus, dimensions, and diff cache

Add:

```go
type syncFocus int

const (
    focusFiles syncFocus = iota
    focusDiff
)

type diffState struct {
    content string
    loading bool
    err     error
}
```

Extend `syncTUIModel` with `focus`, `width`, `height`, `diffs map[string]diffState`, and `diffScroll`.

Behavior:
- `WindowSizeMsg` updates width/height.
- `tab` toggles focus in review mode.
- `Init` loads current diff.
- moving file cursor resets `diffScroll` and loads diff if uncached.
- `d` refreshes current diff.

### Task 3: Render two panes

Render body as left files pane and right main pane:
- left width: clamp between 28 and 48 columns, around one-third of terminal width;
- right width: remaining width;
- height: terminal height minus header/footer/message;
- review mode right pane: diff;
- confirm mode right pane: pending action summary.

Use simple borders and truncation. Do not introduce a new layout library.

### Task 4: Scroll diff pane

When focus is files: `j/k` move selection.
When focus is diff: `j/k` scroll diff.
Action keys work regardless of focus. `tab` switches focus. `enter`, `y`, `esc`, and `q` keep existing semantics.

### Task 5: Verify and commit

Run:

```bash
go test ./...
go build ./cmd/cm
rm -f cm
```

Commit:

```text
feat(ui): add two pane sync review layout
```
