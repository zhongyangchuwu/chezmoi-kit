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

chezmoi:
  source repository has git changes
 M dot_zshrc
```

### The two blocks

`local:` — managed files whose current content differs from the rendered
chezmoi source target. This is what `cm sync` processes.

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

With no targets, `cm diff` diffs every dirty managed file from `cm status`.
With explicit targets, it diffs only those paths.

Diff headers use this direction:

```text
--- chezmoi:<path>
+++ local:<path>
```

Added lines show local content that `cm add` would accept into chezmoi source.
Removed lines show target content that `cm apply` would write locally.

## Sync

```bash
cm sync
cm sync ~/.zshrc
cm sync ~/.zshrc ~/.gitconfig
```

`cm sync` opens a two-pane TUI with changed files on the left and the selected
file's diff on the right. The current file's diff loads automatically on entry
and when you move between files; press `d` to reload it.

Use `cm diff`, `cm add`, `cm apply`, or `cm merge` for non-interactive workflows.

### Sync keys

| Key | Action | Mutates? |
|---|---|---|
| `tab` | switch focus between files and diff pane | no |
| `j` / `down` | move down or scroll diff | no |
| `k` / `up` | move up or scroll diff | no |
| `d` | refresh diff from chezmoi target to local file | no |
| `a` | mark local → chezmoi source | not until confirm |
| `p` | mark chezmoi source → local | not until confirm |
| `m` | mark merge | not until confirm |
| `s` | clear pending action for this entry | no |
| `enter` | review pending actions for confirmation | no |
| `esc` | leave confirm mode | no |
| `q` / `ctrl+c` | quit before execution | no |
| `y` | execute pending actions in confirm mode | yes |

Selecting the same action twice clears it. Selecting a different action for the
same target replaces the previous pending action.

Before executing each action, `cm sync` re-checks that target and drops it if it
is already clean. Confirmed `p` actions run `chezmoi apply --force` because the
TUI has already shown the diff and collected confirmation. Actions execute one
target at a time, so completed files leave the list and remaining files can be
handled in later confirm batches.

Once execution starts, `cm` waits for the current chezmoi command to finish. It
does not advertise `q` as cancellation for already-started mutating subprocesses.
If all entries are resolved, `cm sync` exits with a completion message; if files
remain, it returns to review mode with a remaining-file count.

## Direct commands

Use these when you already know the reconciliation action. Prefer `cm diff` or
`cm sync` when you need to review first.

```bash
cm add ~/.zshrc       # accept local file content into chezmoi source
cm apply ~/.zshrc     # apply chezmoi target content locally
cm merge ~/.zshrc     # open chezmoi merge
cm edit .zshrc        # edit a managed file through chezmoi edit
```

`cm edit <target>` completes from `chezmoi managed` output. Pass managed file
names such as `.zshrc`, not arbitrary shell globs.

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
