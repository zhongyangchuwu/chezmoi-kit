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
  cli_test.go
internal/chezmoi/
  client.go                    chezmoi CLI wrapper
  status.go                    status parsing
internal/syncdiff/
  diff.go                      internal sync diff generation
internal/ui/
  sync.go                      plain sync prompt
  tui.go                       terminal sync UI
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

### `internal/syncdiff`

Generates the TUI diff from rendered chezmoi target content to the current
local file without shelling out to `chezmoi diff`.

### `internal/ui`

Owns the plain sync prompt and terminal sync UI. Injects a `SyncService`
interface for tests. Does not know about chezmoi binary paths or git.

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
