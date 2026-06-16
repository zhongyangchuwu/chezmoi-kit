# Design

## Purpose

`cm` is a thin CLI around chezmoi for personal config reconciliation.

It does not replace chezmoi. It wraps chezmoi commands with a simpler
status model and an interactive reconciliation prompt.

## Design principles

### Default read-only

`cm` and `cm status` never mutate files. Only `cm sync` and the direct
mutating wrappers (`add`, `apply`, `merge`) change local config or
chezmoi source state.

### One local state

`cm` hides chezmoi's three-point comparison (last-written → actual → target).

For this project, local only means: the local file content differs from
the chezmoi source target, or it does not. This matches a personal
config workflow without templates.

The raw chezmoi status code (e.g. `MM`) is not shown to the user.
Instead, status displays `!` for any mismatch.

### Explicit reconciliation

`cm sync` asks the user to choose per-entry:

- add local to chezmoi source
- apply chezmoi target to local
- merge with configured merge tool
- skip

There is no automatic recommendation. The user reviews diffs first.

### Chezmoi source is separate from git history

`cm status` shows two independent facts:

1. `local:` — local managed files that differ from chezmoi source.
2. `chezmoi:` — git status inside the chezmoi source repository.

### Chezmoi remains the authority

`cm` delegates `add`, `apply`, and `merge` to chezmoi.
Diffs are generated internally from rendered chezmoi target content to the current local file.
It does not manipulate chezmoi source files directly.

## Status model

`cm status` runs:

```bash
chezmoi status --path-style=absolute
```

It parses the two-column output for internal use only.

The `local:` block displays:

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

It displays git changes in the chezmoi source repo:

```text
chezmoi:
  source repository has git changes
 M dot_zshrc
```

Both blocks are present in `cm status` output. Neither is shown when
empty.

## Non-goals

- Chezmoi templates are not supported or needed.
- `cm` does not run git commit, push, or pull automatically.
- No daemon, watch, or auto-sync.
- No additional state database beyond chezmoi's own.
- No replacement for `chezmoi`; `cm` always delegates to it.
