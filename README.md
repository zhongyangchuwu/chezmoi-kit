# cm

A small chezmoi reconciliation helper for personal config management.

`cm` does not replace chezmoi. It adds a simpler status view, internal diffs,
and an explicit review TUI before applying changes between local files and the
chezmoi source state.

## Requirements

- Go, for source builds and local install.
- `chezmoi`, required by every command that inspects or mutates managed files.
- `git`, required for source repository status.
- `lazygit`, only required for `cm git`.
- `just`, only required for the helper recipes in this repository.

## Install

Local development install:

```bash
just install
```

Versioned local install:

```bash
VERSION=v0.1.0 just install
```

Release build:

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

## Quick start

```bash
cm          # same as cm status
cm sync     # interactively review and reconcile dirty managed files
```

## Mental model

```text
local config    = working files in $HOME
chezmoi source  = rendered desired config state
chezmoi git     = history for the source repository
```

`cm status` shows two independent facts:

- `local:` managed files whose current local content differs from the rendered chezmoi target.
- `chezmoi:` git working-tree changes inside `chezmoi source-path`.

## Commands

| Command | Mutates? | Meaning |
|---|---|---|
| `cm` | no | Same as `cm status`. |
| `cm status [target...]` | no | Show local mismatch and chezmoi source git status. |
| `cm diff [target...]` | no | Show internal sync diff from rendered target to local file. |
| `cm sync [target...]` | yes, after confirm | TUI review, pending action selection, preflight re-check, confirmed reconciliation. |
| `cm add [target...]` | yes | Run `chezmoi add`; local file content becomes source state. |
| `cm apply [target...]` | yes | Run `chezmoi apply`; source state is applied locally. |
| `cm merge [target...]` | yes | Run `chezmoi merge` for manual conflict resolution. |
| `cm edit <target>` | yes | Run `chezmoi edit` for a managed file. |
| `cm git` | yes | Open `lazygit` in the chezmoi source repository. |
| `cm doctor` | no | Check required and optional environment prerequisites. |
| `cm version` | no | Print build information. |
| `cm completion bash\|zsh\|fish\|powershell` | no | Generate shell completion. |

Optional report controls for `cm`, `cm status`, `cm diff`, `cm doctor`, and `cm version`:

```bash
cm status --output plain       # plain text, default-compatible for scripts
cm status --output ansi        # ANSI-capable text, color auto-detected by default
cm status --output markdown    # Markdown report
cm diff --color never          # disable ANSI color
cm version --color always      # request ANSI color when NO_COLOR is unset
```

`--output` accepts `plain`, `ansi`, or `markdown`. `--color` accepts `auto`,
`always`, or `never`; the default is `auto`, and `NO_COLOR` disables ANSI color unless `--color always` is set.

## Sync safety model

`cm sync` is explicit:

1. Load dirty managed files.
2. Let the user review diffs and mark pending `add`, `apply`, or `merge` actions.
3. Show a confirmation view.
4. Re-check selected targets before executing.
5. Execute only targets that are still dirty.

Confirmed `a` actions run `chezmoi re-add`, which preserves `encrypted_` source
attributes for managed files. Confirmed `apply` actions run `chezmoi apply --force`
because the TUI has already shown the diff and collected confirmation. Once execution
starts, `cm` waits for chezmoi commands to finish; it does not advertise cancellation
for already-started mutating subprocesses.

## Non-goals

- No template authoring helpers.
- No automatic git commit, push, or pull.
- No daemon/watch mode.
- No persistent state database.
- No replacement for `chezmoi`.

## License

MIT. See `LICENSE`.
