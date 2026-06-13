# cm MVP Design

## Goal

Build `cm`, a small Go CLI that makes chezmoi-managed configuration reconciliation safer and faster.

`cm` is not a replacement for chezmoi. It is a thin interactive layer over chezmoi commands. Its default command is read-only.

## Scope

### In scope

- Binary name: `cm`.
- Default command: `cm status`.
- Read-only commands:
  - `cm`
  - `cm status [target...]`
  - `cm diff [target...]`
- Mutating commands:
  - `cm sync [target...]`
  - `cm add [target...]`
  - `cm apply [target...]`
  - `cm merge [target...]`
- Parse `chezmoi status` two-column status output.
- Explain whether a file changed locally, changed in chezmoi target state, or changed on both sides.
- Recommend a safe next command in `cm status`.
- In `cm sync`, present an interactive per-entry menu:
  - `d`: show `chezmoi diff <target>`
  - `a`: run `chezmoi add <target>`
  - `p`: run `chezmoi apply <target>`
  - `m`: run `chezmoi merge <target>`
  - `s`: skip
  - `q`: quit
- Re-read a target's status after a mutating sync action.

### Out of scope

- Reimplementing chezmoi diff, add, apply, or merge.
- Template support or auto-template behavior.
- Git pull, push, commit, or source repository automation inside `cm`.
- Full-screen TUI.
- Daemon/watch mode.
- Extra state database.
- Automatic bulk accept-all behavior.

## Command behavior

### `cm` and `cm status`

`cm` runs the same code path as `cm status`.

`cm status [target...]` calls:

```bash
chezmoi status [target...]
```

It parses each non-empty status line as:

```text
XY path
```

Where:

- `X` is the local destination change relative to the last state written by chezmoi.
- `Y` is the target-state change that `chezmoi apply` would make.

If no changed entries are returned, output:

```text
clean
```

Otherwise output one line per entry:

```text
MM ~/.zshrc  both changed; run cm sync ~/.zshrc
```

`cm status` never modifies home files, chezmoi source files, chezmoi state, or git repositories.

### `cm diff`

`cm diff [target...]` calls:

```bash
chezmoi diff [target...]
```

It is read-only from `cm`'s perspective and does not parse the diff.

### `cm sync`

`cm sync [target...]` calls `chezmoi status [target...]`, filters changed entries, and processes entries sequentially.

Each entry displays:

- raw status code;
- target path;
- human explanation;
- recommended action;
- available keys.

After `add`, `apply`, or `merge`, `cm sync` calls `chezmoi status <target>` again. If the target is clean, sync continues to the next entry. If it still has changes, the same entry is shown again with the fresh status.

### Direct mutating wrappers

These commands directly forward arguments:

```bash
cm add [target...]    -> chezmoi add [target...]
cm apply [target...]  -> chezmoi apply [target...]
cm merge [target...]  -> chezmoi merge [target...]
```

They exist for explicit user intent. They are not used by the read-only default command.

## Status recommendation rules

| Status shape | Meaning | `cm status` text | `cm sync` default |
|---|---|---|---|
| `M ` | Local destination changed | local changed | add |
| ` M` | Target state differs | source changed | apply |
| `MM` | Local and target both changed | both changed | merge |
| `D ` | Local destination deleted | local deleted | inspect |
| ` D` | Target state would delete | source wants delete | apply |
| `A ` | Local destination added | local added | add |
| ` A` | Target state would add | source adds target | apply |
| other | Unclassified state | inspect | diff |

The MVP does not force the recommended action. It only preselects or labels it.

## Architecture

```text
cmd/cm/main.go
internal/chezmoi/
  client.go
  status.go
  status_test.go
internal/reconcile/
  recommend.go
  recommend_test.go
internal/ui/
  sync.go
  sync_test.go
```

### `internal/chezmoi`

Owns process execution and status parsing.

- `Client` stores the chezmoi binary path, stdout, stderr, stdin, and working directory.
- `Run(args ...string) error` executes chezmoi with inherited or supplied streams.
- `Output(args ...string) ([]byte, error)` executes chezmoi and captures stdout.
- `Status(targets []string) ([]StatusEntry, error)` calls `chezmoi status`.
- `ParseStatus(out []byte) ([]StatusEntry, error)` parses status output.

### `internal/reconcile`

Owns pure status interpretation.

- `Describe(entry StatusEntry) string`
- `Recommend(entry StatusEntry) Action`
- `ActionCommand(action Action, target string) []string`

### `internal/ui`

Owns terminal interaction for `cm sync`.

- Accepts injected reader/writer for tests.
- Does not know how to parse chezmoi output.
- Calls a small interface so tests can fake chezmoi without shelling out.

## Error handling

- If `chezmoi` is missing, return a clear error: `chezmoi not found in PATH`.
- If a status line is malformed, return an error naming the line.
- If a forwarded chezmoi command fails, return its exit error after letting stderr pass through.
- In sync, invalid key input re-prompts without running commands.

## Testing

Unit tests cover:

- status parser accepts normal two-column lines, including first-column space statuses;
- parser rejects malformed lines;
- recommendation mapping for local, source, both, delete, add, and unknown states;
- `cm sync` dispatches the expected action for a chosen key;
- `cm sync` re-checks status after a mutating action;
- `cm` and `cm status` do not call mutating commands in command-level tests.

Smoke verification covers:

```bash
go test ./...
go run ./cmd/cm status
```

If the host has no chezmoi state changes, `go run ./cmd/cm status` should print `clean` or the current reconciliation list without modifying files.
