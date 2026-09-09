# Usage

## Requirements

- `chezmoi` on `PATH`.
- `git` on `PATH` for source repository status.
- `lazygit` on `PATH` only for `cm git`.
- Go and `just` only when installing from this repository.

## Install

```bash
just install
```

Install with an explicit version string:

```bash
VERSION=v0.1.0 just install
```

Build a release binary under `dist/`:

```bash
VERSION=v0.1.0 just build-release
./dist/cm version
```

Zsh completion installed by `just install`:

```zsh
fpath=(~/.zfunc $fpath)
autoload -Uz compinit
compinit
```

## Status

```bash
cm          # same as cm status
cm status
cm status ~/.zshrc
```

Example output:

```text
local:
  local config differs from chezmoi source
! /home/me/.zshrc  differs from chezmoi
  run cm sync

automation:
  chezmoi scripts are pending; cm sync does not execute scripts
R /home/me/install-packages.sh  apply would run this script
  use chezmoi diff and chezmoi apply to review and run scripts

chezmoi:
  source repository has git changes
 M dot_zshrc
```

### The three blocks

`local:` — non-script targets whose second chezmoi status column says apply would
change destination state. This is what `cm sync` processes.

`automation:` — scripts that chezmoi apply would run. Scripts can execute
arbitrary commands and do not obey the same clean-after-reconciliation model as
files, so `cm sync` never executes them.

`chezmoi:` — git working-tree changes inside `chezmoi source-path`. This is
shown for awareness only; `cm` does not commit, pull, or push.

### Clean state

```text
clean
```

## Diff

```bash
cm diff
cm diff ~/.zshrc
```

With no targets, `cm diff` diffs every dirty non-script target from `cm status`.
With explicit targets, it diffs only those paths.

`cm` invokes chezmoi with its pager and configured external diff command
disabled:

```text
chezmoi --color=false --no-pager --use-builtin-diff diff --include=all --exclude=none --reverse <target>
```

`--reverse` preserves cm's established direction: removed lines are rendered
target state, and added lines are destination content that a local-to-source
operation would accept. Because chezmoi produces the diff, mode-only changes,
symlinks, directories, templates, and other target-state semantics are included.
Diff capture is bounded; an oversized preview is rejected instead of allowing an
unreviewed action.

If only scripts are pending, targetless `cm diff` explains that there are no
file diffs and directs you to `chezmoi diff` for script review.

## Workspace

```bash
cm ui
cm ui ~/.config
cm ui ~/.config ~/.local/bin
```

`cm ui` is an inspection-first workbench. Browsing and previews do not mutate
state; an explicit source-edit handoff can modify chezmoi source without
automatically changing destination. It opens the selected item's labeled preview:
an authoritative diff for files or a local summary for directories. It retains
clean managed entries instead of exiting when no drift exists. With no paths, it
inventories managed entries and source-ignored entries only; it never scans all
of `$HOME` for unmanaged files. Pass one or more destination paths to discover
unmanaged candidates recursively within those explicit scopes. The workspace
never adds, applies, or merges files.

The Current pane labels `C` clean, `D` dirty, `U` unmanaged, `I` ignored, `R`
script, and `?` uninspected. `[T]` marks templates and `[E]` encrypted source
state. Directories end with `/`; there are no expandable-row glyphs, so state,
type, and name columns align with file rows. A wide terminal adds a read-only
Parent pane and Preview pane; medium widths omit Parent and narrow widths switch
between Current and Preview with `tab`.

The workspace uses semantic color for these markers, but each retains its
letter/badge meaning when color is unavailable. `?` means cm intentionally
skipped authoritative inspection; it is not a clean result.

| Key | Action |
|---|---|
| `tab` | switch Current and Preview focus |
| `j` / `k`, `up` / `down` | move Current selection or scroll the preview |
| `h` / `left` | return to the parent directory in Current focus; scroll preview left in Preview focus |
| `l` / `right` / `enter` | enter the selected directory in Current focus; scroll preview right in Preview focus |
| `f` | cycle all, managed, dirty, unmanaged, ignored, and script filters |
| `/` | search every workspace path in Current focus, or preview text in Preview focus; `enter` locates a path result |
| `n` / `N` | move to next or previous preview match |
| `1` / `2` / `3` / `4` | choose diff, destination, rendered target, or source preview |
| `[` / `]` | move to the previous or next diff hunk |
| `e` | open the selected managed file or symlink through `chezmoi edit`; never auto-apply or watch destination changes |
| `z` | toggle full-screen preview |
| `R` | explicitly reveal a withheld diff, rendered target, or encrypted source |
| `?` | toggle quick-start, complete key reference, and state/type legend |
| `q` / `ctrl+c` | quit without mutation; when help is open, close help instead |
Filtering keeps a directory visible when it contains a matching descendant. A
directory preview is a local matching-child summary and never requests content
from chezmoi; selecting a file still uses the existing authoritative preview.

`e` is available for managed regular files, symlinks, templates, and encrypted
files with a chezmoi source mapping. It opens `chezmoi edit --apply=false
--watch=false` and returns to a full refresh of the original workspace scopes.
The refresh preserves viable directory, filter, focus, preview, and selection
state, clears stale previews and sensitive reveal state, and loads a new
authoritative preview. Direct `cm edit <target>` remains available for shell use.

Preview output is bounded. Rendered template/encrypted targets, encrypted
source, and sensitive uninspected diffs are withheld until `R`; this makes the
reveal an explicit per-target, per-view decision. Plain template source remains
unrendered source text. Chezmoi remains authoritative for all classification,
rendering, decryption, target type, and diff behavior.

The Current pane keeps a compact state legend when space permits. Its footer adapts
to narrow terminals but always retains a help (`?`) and quit (`q`) route. Set
`NO_COLOR=1` to keep the same textual markers and legend without semantic colors.

## Sync

```bash
cm sync
cm sync ~/.zshrc
cm sync ~/.zshrc ~/.gitconfig
```

`cm sync` opens a two-pane TUI with changed non-script targets on the left and
the selected target's type plus authoritative diff on the right. The current
review loads automatically on entry and when you move between targets; press
`d` to reload it.

Use `cm diff`, `cm add`, `cm apply`, or `cm merge` for direct workflows where
you already understand the operation. Direct `cm apply` keeps chezmoi semantics
and may execute scripts.

### Sync keys

| Key | Action | Mutates? |
|---|---|---|
| `tab` | switch focus between files and diff pane | no |
| `j` / `down` | move down or scroll diff | no |
| `k` / `up` | move up or scroll diff | no |
| `d` | refresh type, template state, and diff | no |
| `a` | mark local → chezmoi source | not until confirm |
| `p` | mark chezmoi target → destination | not until confirm |
| `m` | mark merge | not until confirm |
| `s` | clear pending action for this entry | no |
| `enter` | review pending actions for confirmation | no |
| `esc` | leave confirm mode | no |
| `q` / `ctrl+c` | quit before execution | no |
| `y` | execute pending actions in confirm mode | yes |

The available actions depend on the reviewed target:

| Target | Available actions |
|---|---|
| regular file | add, apply, merge |
| template file | apply, merge |
| symlink | apply |
| directory | apply |
| remove entry | apply |
| unknown type | none; use chezmoi directly |

Selecting the same action twice clears it. Selecting another valid action for
the same target replaces it. An action cannot be selected until its review has
loaded successfully.

Each pending action stores the fingerprint of the exact status, target type,
template state, and diff that was reviewed. Immediately before mutation, cm
recomputes that review. If anything changed, the action is deferred, the new
review replaces the old one, and confirmation is required again.

Confirmed `a` actions run `chezmoi re-add`, preserving `encrypted_` source
attributes. Confirmed `p` actions run `chezmoi apply --force` only after the
review fingerprint matches. Add/apply run as non-interactive subprocesses with
captured output; successful output and warnings remain visible. Merge keeps
terminal control because merge tools can be interactive.

After every successful command, cm re-checks the target. Only verified-clean
targets leave the list. A target that still differs returns to review instead of
producing a false completion message. Actions execute one target at a time, so
remaining files can be handled in later confirmation batches.

Once execution starts, `cm` waits for the current chezmoi command to finish. It
does not advertise `q` as cancellation for already-started mutating subprocesses.

### Debug log

Use the global `--debug` flag to write sync timing diagnostics to a temporary
log file:

```bash
cm sync --debug
# or
cm --debug sync
```

`cm` prints the log path to stderr before the TUI starts and again after it
exits. The log records initial status, review loading, review preflight,
execution, postflight review, and batch completion durations. It does not log
rendered diff contents or subprocess output.

## Direct commands

Use these when you already know the reconciliation action. Prefer `cm diff` or
`cm sync` when you need to review first.

```bash
cm add ~/.zshrc       # accept local file content into chezmoi source
cm apply ~/.zshrc     # apply chezmoi target content locally
cm merge ~/.zshrc     # open chezmoi merge
cm edit .zshrc        # edit a managed file through chezmoi edit
```

`cm edit <target>` completes from NUL-delimited `chezmoi managed` output. Relative
names such as `.zshrc` resolve against `chezmoi target-path`, so custom destination
directories work; absolute managed target paths are passed through unchanged.

## Git source repository

```bash
cm git
```

Opens `lazygit` in `chezmoi source-path`. Use it to review, commit, pull, push,
or otherwise manage the chezmoi source repository without leaving the `cm`
workflow. `cm` does not run git operations automatically.

## Typical workflows

### A local tool changed a managed config

```bash
cm
# shows ! for modified files
cm sync
# d → review diff
# a → mark local content for source
# enter → review pending action
# y → execute
```

### You edited chezmoi source manually

```bash
cm
# shows ! because source now differs from local
cm sync
# d → review diff
# p → mark source content for local apply
# enter → review pending action
# y → execute
```

### Both sides need manual resolution

```bash
cm sync
# d → review diff
# m → mark merge
# enter → review pending action
# y → open merge tool
```

## Shell completion

Generate completion scripts:

```bash
cm completion bash > ~/.bash_completions/cm
cm completion zsh > ~/.zfunc/_cm
cm completion fish > ~/.config/fish/completions/cm.fish
cm completion powershell > ~/.powershell/cm.ps1
```

`just install` generates zsh completion automatically.

## Version

```bash
cm version
```

Example output:

```text
cm: v0.1.0
commit: d87c12b78adb45d1fb7a0dbcabc181edd6f5956b
built: 2026-06-18T11:00:31Z
dirty: false
go: go1.26.4
```

Release builds can inject `cm:` with `VERSION=v0.1.0 just build-release`.
When installed directly with `go install`, Go build metadata may provide VCS
fields from the current checkout.
