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
cmd/cm/main.go                  Cobra command wiring
internal/chezmoi/
  client.go                     chezmoi CLI wrapper
  client_test.go
  status.go                     status parsing
  status_test.go
internal/ui/
  sync.go                       interactive sync prompt
  sync_test.go
internal/reconcile/
  recommend.go                  status description helpers
  recommend_test.go
internal/build/
  info.go                       version from runtime/debug
  info_test.go
```

### `cmd/cm`

Wires Cobra commands. Owns the `service` interface and `sourceEntry`
type. The `service` interface is the only place that knows about both
chezmoi and source git operations.

### `internal/chezmoi`

Owns chezmoi CLI execution and status parsing. The `Status` method
adds `--path-style=absolute` so callers always get absolute target paths.
The `ParseStatus` function returns raw two-column codes.

### `internal/ui`

Owns the interactive sync prompt. Injects a `SyncService` interface for
tests. Does not know about chezmoi binary paths or git.

### `internal/reconcile`

Thin helpers for describing status entries. Used by `cmd/cm` and
`internal/ui` for display text.

### `internal/build`

Reads `runtime/debug.ReadBuildInfo()` for `cm version` output.

## Dependencies

```text
github.com/spf13/cobra     CLI framework
github.com/fatih/color     terminal colours
```

No Viper. No TUI framework.

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
