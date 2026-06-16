# Usage

## Install

```bash
just install
```

Zsh completion:

```zsh
fpath=(~/.zfunc $fpath)
autoload -Uz compinit
compinit
```

## Status

```bash
cm          # same as cm status
cm status
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

`local:` — managed files whose current content differs from the chezmoi
source target. This is what `cm sync` will process.

`chezmoi:` — git working-tree changes inside the chezmoi source
repository. Shown for awareness only.

### Clean state

```text
clean
```

## Sync

```bash
cm sync
```

On a terminal, opens a TUI with the changed files, current diff pane,
and key help. In non-terminal input/output, `cm sync` falls back to the
plain prompt:

```text
! /home/me/.zshrc
local differs from chezmoi
[d]iff [a]dd local a[p]ply chezmoi [m]erge [s]kip [q]uit
>
```

### Sync keys

| Key | Action | Mutates? |
|---|---|---|
| `d` | show diff from chezmoi target to local file | no |
| `a` | keep local, write to chezmoi source | yes |
| `p` | discard local, apply chezmoi target | yes |
| `m` | open chezmoi merge | yes |
| `s` | skip this entry | no |
| `q` | quit sync immediately | no |

Use `cm sync --tui` to force the TUI, or `cm sync --plain` to force the
plain prompt. Both modes use the same internal diff for `d`.

After `a`, `p`, or `m`, `cm sync` re-checks the entry's status.
If it is clean, sync continues to the next entry.

In the TUI diff pane, `--- chezmoi:<path>` is the rendered chezmoi target and
`+++ local:<path>` is the current local file. Added lines therefore show local
content that `a` would accept into chezmoi source; removed lines show target
content that `p` would apply locally.

### Sync a single target

```bash
cm sync ~/.zshrc
cm sync ~/.zshrc ~/.gitconfig
```

## Direct commands

Use when you already know the reconciliation action. Use `cm diff` or
`cm sync --tui` to review diffs before choosing an action:

```bash
cm diff ~/.zshrc      # show internal sync diff
cm add ~/.zshrc       # accept local → chezmoi source
cm apply ~/.zshrc     # accept chezmoi source → local
cm merge ~/.zshrc     # open chezmoi merge
```

## Git source repository

```bash
cm git
```

Opens `lazygit` in `chezmoi source-path`. Use it to review, commit,
pull, push, or otherwise manage the chezmoi source repository without
leaving the `cm` workflow. `cm` does not run git operations automatically.

## Typical workflows

### A local tool changed a managed config

```bash
cm
# shows ! for modified files
cm sync
# d → review diff
# a → accept local
```

### You edited chezmoi source manually

```bash
cm
# shows ! because source now differs from local
cm sync
# d → review diff
# p → apply to local
```

### Both sides have changes

```bash
cm sync
# d → review diff
# m → merge manually
```

## Shell completion

Generate completion script:

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

Output:

```text
cm: dev
commit: unknown
built: unknown
dirty: unknown
go: go1.26.4
```

When installed via `go install`, VCS fields are populated from the
current git commit.
