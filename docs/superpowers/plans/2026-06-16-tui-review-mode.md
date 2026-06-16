# TUI Review Mode Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace immediate sync behavior with a TUI-first review/confirm workflow and remove the plain prompt path.

**Architecture:** `internal/reconcile` owns action types and review service contracts. `internal/ui` owns TUI presentation state and invokes one explicit execution command after confirmation. `internal/cli` wires `cm sync` directly to the TUI and adapts chezmoi commands to the review service.

**Tech Stack:** Go, Cobra, Bubble Tea v2, Bubbles help/key, Lipgloss, existing internal `chezmoi`, `reconcile`, and `syncdiff` packages.

---

## File Structure

- Modify `internal/reconcile/service.go`: add action model and `ReviewService` contract.
- Modify `internal/cli/service.go`: implement `Execute([]reconcile.Action)` and use `reconcile.ReviewService` for sync command wiring.
- Modify `internal/cli/cli.go`: remove `--plain`, always run `ui.RunSyncTUI` for `cm sync`.
- Delete `internal/ui/sync.go`: remove plain prompt.
- Rewrite `internal/ui/tui.go`: pending action review model, confirm mode, preflight, execute.
- Modify `internal/ui/sync_test.go`: replace plain/immediate tests with small review-model contract tests.
- Modify `internal/cli/cli_test.go`: remove `--plain` test; keep sync wiring behavior.
- Update `README.md`, `docs/usage.md`, `docs/development.md`, `docs/design.md`: TUI-first sync and no plain prompt.

## Task 1: Reconcile Contracts

**Files:**
- Modify: `internal/reconcile/service.go`
- Modify: `internal/cli/service.go`

- [ ] **Step 1: Define action model**

Add to `internal/reconcile/service.go`:

```go
type ActionKind int

const (
	ActionAdd ActionKind = iota
	ActionApply
	ActionMerge
)

type Action struct {
	Target string
	Kind   ActionKind
}

type ReviewService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	Execute(actions []Action) error
}
```

Keep `Service` only if still needed by direct callers during transition; remove it once TUI uses `ReviewService`.

- [ ] **Step 2: Implement execution adapter**

Add `Execute(actions []reconcile.Action) error` to `chezmoiService` in `internal/cli/service.go`. It loops in order and delegates `ActionAdd`, `ActionApply`, and `ActionMerge` to `runTarget`. Unknown kind returns an error with target context.

- [ ] **Step 3: Run focused compile**

Run: `go test ./internal/reconcile ./internal/cli`

Expected: compile failure until UI wiring updates, or pass if no old interface remains.

## Task 2: Remove Plain Sync Path

**Files:**
- Delete: `internal/ui/sync.go`
- Modify: `internal/cli/cli.go`
- Modify: `internal/cli/cli_test.go`

- [ ] **Step 1: Remove plain prompt file**

Delete `internal/ui/sync.go`.

- [ ] **Step 2: Simplify sync command**

In `internal/cli/cli.go`, remove `--plain` flag and all terminal detection branching. `cm sync` should call:

```go
return ui.RunSyncTUI(services.Sync, args, stdin, stdout)
```

- [ ] **Step 3: Update CLI tests**

Remove `TestRunSyncPlainFlagKeepsPromptMode`. Keep one sync test proving `cm sync --tui` or `cm sync` enters the TUI path and mutates only after confirm once TUI refactor is complete.

- [ ] **Step 4: Run focused compile**

Run: `go test ./internal/cli ./internal/ui`

Expected: compile failures in UI tests until TUI refactor is complete.

## Task 3: Review TUI Model

**Files:**
- Rewrite: `internal/ui/tui.go`
- Modify: `internal/ui/sync_test.go`

- [ ] **Step 1: Replace immediate model fields**

Use model fields:

```go
type syncMode int

const (
	modeReview syncMode = iota
	modeConfirm
)

type syncTUIModel struct {
	service    reconcile.ReviewService
	entries    []chezmoi.StatusEntry
	cursor     int
	pending    map[string]reconcile.ActionKind
	diffTarget string
	diff       string
	mode       syncMode
	help       help.Model
	message    string
	err        error
}
```

- [ ] **Step 2: Add pending helpers**

Implement:

```go
func (m syncTUIModel) togglePending(kind reconcile.ActionKind) syncTUIModel
func (m syncTUIModel) clearPending() syncTUIModel
func (m syncTUIModel) pendingActions() []reconcile.Action
```

`pendingActions` returns actions in `entries` order.

- [ ] **Step 3: Add diff loading command**

Keep `d` behavior as a command that calls `DiffOutput(current.Path)`. Cache by `diffTarget`.

- [ ] **Step 4: Add confirm mode**

`enter` switches to confirm mode only when pending action count > 0. In confirm mode:

- `y` runs execute command;
- `esc` returns to review;
- `q` quits without executing.

- [ ] **Step 5: Add preflight execution command**

Execution command:

1. Build `pendingActions()`.
2. Call `Status(targets)` for those targets.
3. Drop actions whose target is absent from fresh status.
4. If none remain, return clean message without calling `Execute`.
5. Call `Execute(actions)`.

- [ ] **Step 6: Keep tests small**

Replace `internal/ui/sync_test.go` with tests for:

- toggling same action clears it;
- replacing action changes pending kind;
- pending actions preserve entry order;
- confirm execution preflight drops clean targets.

Do not assert exact help text, prompt wording, colors, or full layout.

## Task 4: Documentation Cleanup

**Files:**
- Modify: `README.md`
- Modify: `docs/usage.md`
- Modify: `docs/development.md`
- Modify: `docs/design.md`

- [ ] **Step 1: Remove plain references**

Delete `--plain`, plain prompt, and `RunSync` references.

- [ ] **Step 2: Document review workflow**

Document that `cm sync` is TUI-first: mark actions, enter confirm, `y` execute, `esc` review, `q` quit.

## Task 5: Verification and Commit

**Files:** all modified files.

- [ ] **Step 1: Format**

Run: `gofmt -w internal/reconcile/service.go internal/cli/service.go internal/cli/cli.go internal/cli/cli_test.go internal/ui/tui.go internal/ui/sync_test.go`

- [ ] **Step 2: Test**

Run: `go test ./...`

Expected: all packages pass.

- [ ] **Step 3: Build**

Run: `go build ./cmd/cm`

Expected: no output and exit code 0.

- [ ] **Step 4: Remove build artifact**

Run: `rm -f cm`

- [ ] **Step 5: Commit**

Commit message:

```text
refactor(sync): make TUI review mode primary
```
