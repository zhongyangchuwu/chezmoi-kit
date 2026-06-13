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

Scans all changed entries and presents each one with a prompt:

```text
! /home/me/.zshrc
local differs from chezmoi
[d]iff [a]dd local a[p]ply chezmoi [m]erge [s]kip [q]uit
>
```

### Keys

| Key | Action | Mutates? |
|---|---|---|
| `d` | show chezmoi diff | no |
| `a` | keep local, write to chezmoi source | yes |
| `p` | discard local, apply chezmoi target | yes |
| `m` | open chezmoi merge | yes |
| `s` | skip this entry | no |
| `q` | quit sync immediately | no |

After `a`, `p`, or `m`, `cm sync` re-checks the entry's status.
If it is clean, sync continues to the next entry.

### Sync a single target

```bash
cm sync ~/.zshrc
cm sync ~/.zshrc ~/.gitconfig
```

## Direct commands

Use when you already know what you want:

```bash
cm add ~/.zshrc       # accept local → chezmoi source
cm apply ~/.zshrc     # accept chezmoi source → local
cm merge ~/.zshrc     # open chezmoi merge
cm diff ~/.zshrc      # show chezmoi diff
```

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
