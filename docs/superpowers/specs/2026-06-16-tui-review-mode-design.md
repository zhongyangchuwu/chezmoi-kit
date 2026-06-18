# TUI Review Mode Design

## Goal

Make `cm sync` a TUI-first review workflow: inspect changed managed files, mark pending reconciliation actions, confirm once, re-check status, then execute the selected actions.

## Current behavior being replaced

`cm sync` currently has two paths:

- plain prompt: reads one-letter choices from stdin and executes immediately.
- TUI: renders a file list and diff area, but `a`, `p`, and `m` still execute immediately.

The plain path is removed. `cm sync` always runs the TUI. Non-interactive callers should use explicit commands (`cm add`, `cm apply`, `cm merge`) or `cm diff`.

## User workflow

1. `cm sync [target...]` loads dirty managed targets with `chezmoi status --path-style=absolute`.
2. If no targets are dirty, it prints `clean` and exits.
3. The TUI shows:
   - changed files list;
   - current file diff;
   - pending action marker per file;
   - help/status line.
4. Navigation changes the selected file and loads its diff lazily.
5. `a`, `p`, and `m` mark the current file with an action:
   - `a`: add local file into chezmoi source;
   - `p`: apply chezmoi target locally;
   - `m`: open chezmoi merge for target.
6. Repeating the same action on a file clears that pending action. Choosing a different action replaces the previous action.
7. `s` clears any pending action for the current file.
8. `enter` opens confirm mode when at least one action is pending.
9. In confirm mode:
   - `y` executes pending actions;
   - `esc` returns to review;
   - `q` quits without executing.
10. Before execution, `cm` re-runs status for pending targets. Any target that is no longer dirty is dropped from the execution set. Dirty targets proceed.
11. Actions execute in original status order.
12. After execution, `cm` exits. Direct command wrappers remain available for scripting.

## Domain contract

`internal/reconcile` owns sync/reconcile domain types:

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

The UI owns presentation state only. It may store pending `reconcile.Action` values, but action semantics and execution belong to `internal/reconcile` service implementations.

`internal/cli.chezmoiService` implements `ReviewService` by delegating:

- `ActionAdd` → batched `chezmoi add <targets...>`;
- `ActionApply` → batched `chezmoi apply --force <targets...>` because the TUI has already shown the diff and collected explicit confirmation;
- `ActionMerge` → sequential `chezmoi merge <target>` after add/apply batches, so each merge tool session blocks the current process and can surface failures immediately.

Batch execution stops on the first error and returns context naming the failed action and target set.

## TUI model

The TUI model contains:

- `entries []chezmoi.StatusEntry`: current changed files.
- `cursor int`: selected file.
- `pending map[string]reconcile.ActionKind`: pending actions keyed by absolute target path.
- `diffTarget string` and `diff string`: cached diff for selected target.
- `mode`: review or confirm.
- `message string` and `err error`.

The model should not call `Add`, `Apply`, or `Merge` directly. It can call `DiffOutput` while reviewing. Execution is a single command from confirm mode that performs preflight and `Execute`.

## Preflight behavior

Before executing:

1. Build pending actions in current `entries` order.
2. Run `Status(pendingTargets)`.
3. Keep only actions whose target is still present in the returned status entries.
4. If no actions remain, exit cleanly with a message.
5. Execute remaining actions in order.

This avoids applying stale choices to files that became clean after the TUI loaded.

## CLI behavior

`cm sync` always runs the TUI. Remove:

- `--plain` flag;
- `RunSync` plain prompt;
- plain prompt tests and docs.

Keep:

- `cm diff [target...]` for non-interactive diff review;
- direct `cm add/apply/merge [target...]` wrappers.

## Testing approach

Keep tests light and contract-focused:

- action toggle/replace/clear behavior in the model;
- confirm executes pending actions in file-list order;
- preflight drops clean targets before execute;
- CLI `sync` wires to TUI path and no longer exposes `--plain`.

Do not test prompt wording, exact help text, colors, or full layout snapshots. Manual terminal testing remains primary for interaction quality.

## Non-goals

- No external diff renderer.
- No git commit/push/pull automation.
- No persistent state database.
- No support for the removed plain prompt.
