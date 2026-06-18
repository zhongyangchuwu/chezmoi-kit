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

Opens a two-pane TUI with changed files on the left and the selected file's
diff on the right. Diffs load lazily when you switch into the diff pane or
press `d`. `cm sync` is interactive; use `cm diff`, `cm add`, `cm apply`, or
`cm merge` for non-interactive workflows.

### Sync keys

| Key | Action | Mutates? |
|---|---|---|
| `tab` | switch focus between files and diff pane | no |
| `d` | refresh diff from chezmoi target to local file | no |
| `a` | mark local → chezmoi source | not until confirm |
| `p` | mark chezmoi source → local | not until confirm |
| `m` | mark merge | not until confirm |
| `s` | clear pending action for this entry | no |
| `enter` | review pending actions for confirmation | no |
| `y` | execute pending actions in confirm mode | yes |
| `esc` | leave confirm mode | no |
| `q` | quit without executing more actions | no |

In files focus, `j/k` moves between files without loading diffs. In diff
focus, `j/k` scrolls the diff. Selecting the same action twice clears it.
Selecting a different action for the same target replaces the previous pending
action. Before execution, `cm sync` re-checks selected targets and drops any
target that is already clean.
Confirmed `p` actions run `chezmoi apply --force` because the TUI has already
shown the diff and collected confirmation. Confirmed `a` and `p` actions are
batched; confirmed `m` actions run last, one at a time.

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
`cm sync` to review diffs before choosing an action:

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
