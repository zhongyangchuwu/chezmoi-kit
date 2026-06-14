# cm

A small chezmoi reconciliation helper for personal config management.

## What it does

- Shows whether local config differs from chezmoi source.
- Shows git changes inside the chezmoi source repo.
- Lets you interactively choose add / apply / merge.
- Keeps `cm` itself read-only by default.

## Mental model

```
local config  =  working tree
chezmoi source  =  config stage
chezmoi git    =  history
```

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

## Quick start

```bash
cm          # show status
cm sync     # interactively reconcile all
```

## Commands

| Command | Mutates? | Meaning |
|---|---|---|
| `cm` | no | same as `cm status` |
| `cm status` | no | show local mismatch and chezmoi git status |
| `cm diff [target...]` | no | forward to `chezmoi diff` |
| `cm sync [target...]` | yes | interactive reconciliation |
| `cm add [target...]` | yes | local → chezmoi source |
| `cm apply [target...]` | yes | chezmoi source → local |
| `cm merge [target...]` | yes | open chezmoi merge |
| `cm git` | yes | open lazygit in the chezmoi source repo |
| `cm version` | no | build info |
| `cm completion bash\|zsh\|fish\|powershell` | no | shell completion |

## Non-goals

- No templates.
- No automatic git commit.
- No daemon/watch mode.
- No remplacer for chezmoi.
