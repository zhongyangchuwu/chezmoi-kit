# Design

## Purpose

`cm` is a thin CLI around chezmoi for personal config reconciliation.

It keeps chezmoi as the authority and adds:

- a simplified read-only status view;
- internal rendered-target vs local-file diffs;
- explicit direct wrappers for known actions;
- an interactive review TUI for choosing reconciliation actions.

## Design principles

### Default read-only

`cm`, `cm status`, `cm diff`, `cm version`, and `cm completion` do not mutate
managed files. Mutations happen only through `cm sync` after confirmation or
through direct explicit wrappers: `add`, `apply`, `merge`, `edit`, and `git`.

### One local mismatch state

Chezmoi exposes a richer three-point comparison. `cm` intentionally shows a
smaller user model: a managed local file either differs from the rendered
chezmoi target or it does not.

The raw chezmoi status code, such as `MM`, is parsed for internal routing but is
not shown to users. The status view displays `!` for any local mismatch.

### Explicit reconciliation

`cm sync` opens a review TUI. The user marks per-entry actions, reviews the
pending set, and confirms a batch. Confirmed actions execute one target at a
time; each completed or newly clean target leaves the list. If files remain, the
TUI returns to review mode so the user can handle the next batch. If none remain,
it exits with a completion message.

Before executing each action, `cm` re-checks that target and drops it if it is
already clean. Confirmed execution then runs to completion for the current target
or returns an error; `cm` does not advertise quit as cancellation for
already-started mutating subprocesses.

Available actions:

- add local content to chezmoi source;
- apply chezmoi target content locally;
- merge with the configured chezmoi merge tool;
- skip by leaving the entry unmarked.

There is no automatic recommendation. The user reviews diffs first.

### Chezmoi source is separate from git history

`cm status` shows two independent facts:

1. `local:` — local managed files that differ from rendered chezmoi source.
2. `chezmoi:` — git status inside the chezmoi source repository.

`cm git` opens `lazygit` in the chezmoi source directory, but `cm` never commits,
pulls, or pushes automatically.

### Chezmoi remains the authority

`cm` delegates `add`, `apply`, `merge`, `edit`, `source-path`, `status`,
`managed`, and `cat` behavior to the chezmoi executable.

Diffs are generated internally from rendered chezmoi target content to the
current local file. `cm` does not manipulate chezmoi source files directly.

## Package architecture

```text
cmd/cm                  process entrypoint
internal/cli            Cobra command tree, renderers, concrete service adapter
internal/chezmoi        chezmoi executable wrapper and output parsers
internal/process        external process runner abstraction
internal/reconcile      sync action types and TUI-facing service contract
internal/syncdiff       internal diff generation from content sources
internal/ui             Bubble Tea sync review TUI
internal/build          version/build metadata formatting
```

### Entry point

`cmd/cm/main.go` passes process arguments and standard streams into
`internal/cli.Main`. It owns no command behavior.

### CLI composition

`internal/cli` builds the Cobra command tree and groups narrow service
interfaces in `commandServices`. This keeps command tests independent from real
chezmoi, git, and lazygit processes.

### External process boundary

`internal/process.Runner` is the only package-level abstraction over
`os/exec`. `internal/chezmoi.Client` uses that runner for chezmoi commands, and
`internal/cli.chezmoiService` uses it for git and lazygit.

### Reconciliation boundary

`internal/reconcile.ReviewService` is the contract consumed by the TUI:

```go
type ReviewService interface {
    Status(targets []string) ([]chezmoi.StatusEntry, error)
    DiffOutput(target string) ([]byte, error)
    ExecuteOne(action Action) error
}
```

The UI owns presentation state. The concrete CLI service owns command execution.

### Diff boundary

`internal/syncdiff.Differ` depends on a `ContentSource` interface. The chezmoi
content loader implements that interface by reading rendered target content via
`chezmoi cat` and local content from the filesystem.

## Status model

`cm status` runs:

```bash
chezmoi status --path-style=absolute
```

It parses non-empty lines as:

```text
XY path
```

The parsed code is internal. The public `local:` block displays:

```text
local:
  local config differs from chezmoi source
! /home/me/.zshrc  differs from chezmoi
  run cm sync
```

The `chezmoi:` block comes from:

```bash
chezmoi source-path
git -C <source-dir> status --porcelain=v1
```

Malformed non-empty git porcelain lines are treated as errors rather than being
silently skipped.

## Non-goals

- Chezmoi template authoring helpers.
- Automatic git commit, push, or pull.
- Daemon, watch, or auto-sync.
- Persistent state database beyond chezmoi's own state.
- Replacement for `chezmoi`; `cm` always delegates to it.
