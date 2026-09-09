# Design

## Purpose

`cm` is a thin CLI around chezmoi for personal configuration reconciliation.
It owns review and decision UX; chezmoi remains authoritative for target-state
calculation, templates, encryption, source naming, diff semantics, and mutation.

`cm` adds:

- a simplified read-only status view;
- bounded, forced chezmoi builtin diffs;
- explicit direct wrappers for known actions;
- a target-aware review TUI with stale-state and postflight checks.

## Design principles

### Default read-only

`cm`, `cm status`, `cm diff`, `cm version`, and `cm completion` do not mutate
managed state. Mutations happen only through confirmed `cm sync` actions or the
explicit direct wrappers `add`, `apply`, `merge`, `edit`, and `git`.

### Chezmoi owns target semantics

`cm` must not reconstruct target state by reading destination bytes. A chezmoi
target can be a regular file, template, symlink, directory, remove entry, or
script, and changes can include content, permissions, type, or execution.

Preview output is generated with:

```text
chezmoi --color=false --no-pager --use-builtin-diff diff --include=all --exclude=none --reverse <target>
```

The forced builtin diff prevents configured external diff programs or pagers
from taking over the TUI. `--reverse` preserves cm's public direction: rendered
target is the old side and destination is the new side. Output is bounded before
it enters reports or TUI state.

### Scripts are not file reconciliation

The second `chezmoi status` column uses `R` when apply would execute a script.
Scripts may be stateful, non-idempotent, or always pending, so they do not fit a
clean-after-file-reconciliation invariant.

`cm status` reports scripts in a separate `automation:` block. Targetless
`cm diff` and `cm sync` exclude them and direct users to `chezmoi diff` and
`chezmoi apply`. Direct `cm apply` remains a thin wrapper and therefore retains
chezmoi's script behavior.

### Inspection-first workspace inventory

`cm ui [path...]` is a persistent inspection surface, separate from focused
`cm sync` reconciliation. Its entries are constructed from chezmoi managed,
status, typed managed, and ignored queries. Scoped unmanaged discovery delegates
each candidate directory back to chezmoi; without explicit paths it does not
traverse the destination directory.

The workspace derives a Yazi-inspired directory browser from one absolute-path
inventory. Its Current pane contains direct children only; missing ancestors are
virtual navigation directories synthesized from inventory paths. Directories use
a trailing `/`, not a disclosure glyph. Wide terminals show Parent, Current, and
Preview panes; medium terminals omit Parent; narrow terminals switch Current and
Preview by focus. In Current focus, `h`/left leaves a directory and
`l`/right/enter enters one. Preview focus retains horizontal `h`/`l` scrolling.
Global path search temporarily lists matching inventory nodes and locates the
selected result on Enter; filters retain ancestors with matching descendants.

File previews load only for selected files and are bounded. A virtual or actual
directory preview is a local child summary, never a chezmoi content call. Diff,
destination, rendered target, and source remain distinct file views.
Template/encrypted rendered targets, encrypted source, and uninspected diffs
require explicit reveal. The workspace has no add, apply, or merge actions.

`e` is the one explicit mutation handoff: an eligible managed file or symlink
opens through `chezmoi edit --apply=false --watch=false`. Bubble Tea yields the
terminal to the configured editor, then re-inventories original scopes regardless
of editor exit status. Refresh preserves viable navigation state but discards old
preview, reveal, search-match, and scroll caches before loading a fresh preview.
This prevents inherited chezmoi edit settings from applying or watching destination
changes outside cm review.

Workspace presentation uses a semantic palette for state, target type, template,
encryption, selection, loading, error, and withheld feedback. Letters and badges
remain the authoritative visual meaning so no-color terminals remain usable. `?`
opens a keyboard-isolated help view containing quick start, contextual key groups,
and the complete state/type legend; `Esc`, `?`, or `q` returns to the workspace.


### Reviews bind confirmation to exact state

The app builds a `Review` from:

- target path and raw status code;
- chezmoi target type;
- template membership;
- authoritative diff text.

A SHA-256 fingerprint covers those fields. A pending `Action` contains the
fingerprint from the review that the user selected.

Immediately before execution, cm recomputes the review. If the target is now
clean, the action is skipped. If the fingerprint changed or the action is no
longer valid for the target type, the action is deferred without mutation and
the refreshed review replaces the stale one.

### Command success is not reconciliation success

After a successful add, apply, or merge command, cm recomputes the target review.
Only a verified-clean target leaves the TUI. A target that still differs remains
visible and requires another explicit decision. This prevents false completion
for template re-add refusals, no-op operations, races, and other successful
commands that do not establish the requested state.

Successful non-interactive stdout and stderr are preserved as action notices.
Rendered contents and action output are never written to the timing log.

### Conservative action matrix

| Reviewed target | Add | Apply | Merge |
|---|---:|---:|---:|
| regular file | yes | yes | yes |
| template file | no | yes | yes |
| symlink | no | yes | no |
| directory | no | yes | no |
| remove entry | no | yes | no |
| script | no | no | no |
| unknown type | no | no | no |

`re-add` preserves encrypted attributes but ignores non-files and refuses to
overwrite templates. Unsupported operations stay unavailable until real chezmoi
behavior justifies broadening the matrix.

### Chezmoi source and git history remain separate

`cm status` can show three independent blocks:

1. `local:` — non-script destination/target mismatch.
2. `automation:` — scripts that chezmoi apply would run.
3. `chezmoi:` — git status inside the source repository.

`cm git` opens lazygit in the source directory. `cm` does not automatically
commit, pull, push, or combine source git changes with reconciliation actions.

## Package architecture

```text
cmd/cm                  process entrypoint
internal/cli            Cobra command tree, stream plumbing, exit behavior
internal/app            use cases, semantic reports, reconciliation domain and services
internal/chezmoi        chezmoi executable adapter and strict output parsers
internal/process        external process runner abstraction
internal/report         semantic report model and renderers
internal/tui            Bubble Tea sync and inspection-first workspace TUI
```

### Entry point and CLI

`cmd/cm/main.go` passes arguments and standard streams to `internal/cli.Main`.
`internal/cli` constructs Cobra commands, validates flags, selects report
renderers, and maps command errors to process exit status. It does not calculate
chezmoi state or own sync execution.

### App services and reconciliation domain

`internal/app` owns application use cases and the values consumed by the TUI:

```go
type SyncStatus struct {
    Entries []ReconcileEntry
    Scripts []ReconcileEntry
}

type Review struct {
    Entry       ReconcileEntry
    Type        TargetType
    Template    bool
    Diff        string
    Fingerprint string
    Dirty       bool
}

type Action struct {
    Target      string
    Kind        ActionKind
    Fingerprint string
}

type SyncService interface {
    Status(targets []string) (SyncStatus, error)
    Review(target string) (Review, error)
    ExecuteNonInteractive(action Action) (ActionResult, error)
    TerminalCommand(action Action) (TerminalCommand, error)
}
```

The TUI consumes app-owned reconciliation and workspace values; it never
consumes infrastructure-owned `chezmoi.StatusEntry` values directly.

### Chezmoi adapter

`internal/chezmoi.Client` is the only app-facing chezmoi command adapter. It
provides:

- absolute-path status parsing and `--skip-secrets` inventory status;
- bounded authoritative diff and selected target/source preview output;
- typed managed inventory, NUL-delimited ignored/unmanaged lists, and source mappings;
- configured destination lookup through `target-path`;
- buffered and terminal-bound execution paths.

Target metadata is loaded only for the selected target. It is never written to
debug logs.

### External process boundary

`internal/process.Runner` is the package-level abstraction over `os/exec`.
Chezmoi, source git status, and lazygit all use it. Add and apply capture output
without inheriting the active TUI terminal. Merge receives terminal control
through Bubble Tea's command handoff because the configured merge tool can be
interactive.

### TUI boundary

`internal/tui` owns shared presentation state for the focused sync review and
the inspection-first workspace. Workspace mode keeps clean entries, derives
direct-child directory projections and virtual ancestors from one inventory, and
caches bounded previews by target/view/reveal state. Explicit source edit yields
the terminal to chezmoi, then invalidates workspace caches and re-inventories the
original scopes. Sync-only confirmation, execution, and completion behavior
remain isolated behind explicit mode branches.

### Report boundary

`internal/report` defines semantic blocks, inline roles, diff-line
classification, and plain/ANSI/Markdown renderers. App services construct
semantic status, diff, doctor, and version documents. CLI and TUI presentation
reuse the same diff-line classifier.

## Status and diff boundaries

Raw chezmoi status lines are parsed strictly as:

```text
XY path
```

The raw code remains internal. A space in the second column means destination
already matches target and creates no reconciliation item. `R` in the second
column enters `automation:`; other second-column effects become `ReconcileEntry`
values.

Source repository status comes from:

```text
chezmoi source-path
git status --porcelain=v1
```

Malformed non-empty git porcelain lines are errors rather than silently skipped.

## Edit target resolution

`chezmoi managed` completion returns paths relative to the configured destination.
For a relative `cm edit <target>`, app services call `chezmoi target-path` and
join against that directory. Absolute targets are passed unchanged. This avoids
assuming that chezmoi's destination is `$HOME`.

## Non-goals

- Template authoring or template-data debugging in the authoritative reconciliation phase.
- Automatic git commit, push, pull, or update.
- Daemon, watch, or automatic synchronization.
- Persistent state beyond chezmoi.
- Replacing chezmoi diff, template, encryption, source naming, or merge logic.
- Cancellation of already-started mutating subprocesses without an explicit context-aware runner contract.
