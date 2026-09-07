# cm

A small chezmoi reconciliation helper for personal config management.

`cm` does not replace chezmoi. It adds a simpler status view, forces chezmoi's
builtin diff for authoritative previews, and provides an explicit review TUI
before applying changes between destination and source state.

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
cm                  # same as cm status
cm ui ~/.config     # browse managed state and scoped unmanaged candidates
cm sync             # interactively review and reconcile dirty managed files
```

## Mental model

```text
chezmoi destination = working configuration files (usually $HOME)
chezmoi target      = rendered desired state
chezmoi source      = source files, templates, and encrypted state
chezmoi git         = history for the source repository
```

`cm status` shows three independent facts when present:

- `local:` non-script targets whose destination state differs from the rendered target.
- `automation:` chezmoi scripts that apply would execute; `cm sync` does not execute them.
- `chezmoi:` git working-tree changes inside `chezmoi source-path`.

## Commands

| Command | Mutates? | Meaning |
|---|---|---|
| `cm` | no | Same as `cm status`. |
| `cm status [target...]` | no | Show local mismatch and chezmoi source git status. |
| `cm diff [target...]` | no | Show bounded chezmoi builtin diff from rendered target to destination. |
| `cm sync [target...]` | yes, after confirm | Review typed targets, bind actions to reviewed state, confirm, and verify reconciliation. |
| `cm ui [path...]` | no | Persistent workspace for managed, ignored, and explicitly scoped unmanaged files. |
| `cm add [target...]` | yes | Run `chezmoi add`; local destination content becomes source state. |
| `cm apply [target...]` | yes | Run `chezmoi apply`; may also run chezmoi scripts. |
| `cm merge [target...]` | yes | Run `chezmoi merge` for manual conflict resolution. |
| `cm edit <target>` | yes | Run `chezmoi edit`; relative paths resolve from chezmoi's configured destination. |
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

1. Load dirty targets and report scripts separately.
2. Load the selected target's type, template state, and forced chezmoi builtin diff.
3. Offer only actions valid for that review: regular files may add/apply/merge,
   templates may apply/merge, and symlinks/directories/removes may apply.
4. Bind each pending action to the reviewed fingerprint and show confirmation.
5. Recompute the review immediately before execution; changed reviews are deferred.
6. Execute one target at a time and re-check it afterward.
7. Remove only targets that are verified clean; unresolved targets return to review.

Confirmed `a` actions run `chezmoi re-add`, preserving `encrypted_` source
attributes. Confirmed `p` actions run `chezmoi apply --force` because the exact
reviewed state was confirmed. Successful subprocess warnings remain visible.
Once execution starts, `cm` waits for the current chezmoi command to finish; it
does not advertise cancellation for already-started mutating subprocesses.

Chezmoi scripts are deliberately outside this file-reconciliation flow. Use
`chezmoi diff` and `chezmoi apply` when you intend to inspect and execute them.

## Non-goals

- No template authoring helpers.
- No automatic git commit, push, or pull.
- No daemon/watch mode.
- No persistent state database.
- No replacement for `chezmoi`.

## License

MIT. See `LICENSE`.
