# Development

## Requirements

- Go
- chezmoi
- just

## Test

```bash
go test ./...
```

Targeted:

```bash
go test ./internal/chezmoi
go test ./internal/ui
go test ./cmd/cm
```

## Run locally

```bash
go run ./cmd/cm status
go run ./cmd/cm sync
go run ./cmd/cm diff ~/.zshrc
```

## Install locally

```bash
just install
```

## Project layout

```text
cmd/cm/main.go                  thin process entrypoint
internal/cli/
  cli.go                       Cobra command wiring
  service.go                   chezmoi-backed CLI service
  status.go                    status rendering
  diff.go                      diff command rendering
  service_test.go
  cli_test.go
internal/chezmoi/
  client.go                    chezmoi CLI wrapper
  status.go                    status parsing
  content.go                   chezmoi/local content loader for sync diffs
internal/syncdiff/
  diff.go                      internal sync diff generation
internal/reconcile/
  service.go                   shared reconciliation service contract
internal/ui/
  tui.go                       terminal sync program entry/update
  model.go                     sync TUI state and pending actions
  view.go                      two-pane layout rendering
  diff.go                      diff cache, scroll, and coloring
  confirm.go                   confirm execution command
  keys.go                      key bindings and styles
  update.go                    key handling
internal/process/
  runner.go                    external process execution abstraction
internal/build/
  info.go                      version from runtime/debug
```

### `cmd/cm`

Owns only process startup and delegates to `internal/cli`.

### `internal/cli`

Wires Cobra commands, renders command output, and adapts `internal/chezmoi` to
the interfaces consumed by command handlers and sync UIs.

### `internal/chezmoi`

Owns chezmoi CLI execution and status parsing. The `Status` method
adds `--path-style=absolute` so callers always get absolute target paths.
The `ParseStatus` function returns raw two-column codes.
The `ContentLoader` type adapts chezmoi-rendered target content and local files
to `internal/syncdiff`.

### `internal/syncdiff`

Generates sync diffs from rendered chezmoi target content to the current
local file without shelling out to `chezmoi diff`.

### `internal/reconcile`

Defines the shared reconciliation action model and review service contract used
by the terminal UI. It keeps sync-domain capabilities out of UI presentation packages.

### `internal/process`

Owns subprocess execution. Chezmoi, git status, and lazygit all go through this
runner boundary so tests can inject one process fake.

### `internal/ui`

Owns the terminal sync UI. Depends on the `internal/reconcile` service contract
and does not know about chezmoi binary paths or git.

### `internal/build`

Reads `runtime/debug.ReadBuildInfo()` for `cm version` output.

## Dependencies

```text
charm.land/bubbletea/v2          terminal UI runtime
charm.land/bubbles/v2            TUI help/key bindings
charm.land/lipgloss/v2           TUI styling
github.com/spf13/cobra           CLI framework
github.com/fatih/color           plain terminal colours
github.com/rogpeppe/go-internal  anchored unified diff
```

No Viper. No external diff renderer.

## Build info

Version information comes from Go's build info. When installed via
`go install`, VCS fields are populated automatically:

```bash
cm version
```

For release builds, set version explicitly:

```bash
go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/build.Version=v0.1.0"
```
