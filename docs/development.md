# Development

## Requirements

- Go matching `go.mod`.
- `chezmoi` for commands that inspect or mutate managed files.
- `git` for source repository status checks.
- `lazygit` for `cm git`.
- `just` for repository helper recipes.

## Run locally

```bash
go run ./cmd/cm status
go run ./cmd/cm sync
go run ./cmd/cm diff ~/.zshrc
go run ./cmd/cm edit .zshrc
```

## Test

Project-wide:

```bash
go test ./...
go test -race ./...
go vet ./...
go mod tidy -diff
go mod verify
```

Targeted:

```bash
go test ./internal/cli
go test ./internal/chezmoi
go test ./internal/syncdiff
go test ./internal/ui
go test ./internal/build
```

Vulnerability scan:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Install locally

```bash
just install
```

With explicit version injection:

```bash
VERSION=v0.1.0 just install
```

The install recipe runs `go install` and writes zsh completion to
`${HOME}/.zfunc/_cm`.

## Release build

For local one-platform builds:

```bash
VERSION=v0.1.0 just build-release
./dist/cm version
```

The local recipe uses:

```bash
go build -trimpath -ldflags "-s -w -X github.com/zhongyangchuwu/cm/internal/build.Version=$VERSION" -o dist/cm ./cmd/cm
```

For release artifacts, use GoReleaser:

```bash
go run github.com/goreleaser/goreleaser/v2@latest check
go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean
```

GoReleaser reads `.goreleaser.yaml`, builds Linux/macOS/Windows archives, injects
`internal/build.Version={{ .Version }}`, and writes checksums under `dist/`.

## v0.1.0 local release gate

Run before tagging:

```bash
go test ./...
go test -race ./...
go vet ./...
go mod tidy -diff
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go run github.com/goreleaser/goreleaser/v2@latest check
go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean
```

Manual smoke in a real terminal:

```bash
cm status
cm diff
cm sync
cm edit <known-managed-file>
cm git
cm completion zsh
```

For `cm sync`, verify:

- Ctrl+C restores the terminal.
- Quitting before confirmation does not execute actions.
- Confirmed execution no longer advertises `q` as cancellation.
- Errors from chezmoi commands are returned to the CLI.

## GitHub workflows

`.github/workflows/ci.yml` runs on pushes to `main` and pull requests. It runs
module tidy checks, module verification, tests, race tests, vet, govulncheck,
GoReleaser config validation, and a GoReleaser snapshot build.

`.github/workflows/release.yml` runs on `v*` tags. It uses
`goreleaser/goreleaser-action@v7` with `fetch-depth: 0`, Go from `go.mod`, and
`GITHUB_TOKEN` with `contents: write` to publish GitHub release artifacts.

## Tag flow

After local checks and CI pass:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Pushing the tag triggers the GoReleaser workflow and publishes the GitHub
release artifacts.

## Project layout

```text
cmd/cm/main.go                  thin process entrypoint
internal/cli/
  cli.go                       Cobra command wiring
  service.go                   chezmoi-backed CLI service
  status.go                    status rendering
  diff.go                      diff command rendering
  service_test.go              service adapter tests
  cli_test.go                  command wiring tests
internal/chezmoi/
  client.go                    chezmoi CLI wrapper
  status.go                    status and managed-file parsing
  content.go                   chezmoi/local content loader for sync diffs
internal/syncdiff/
  diff.go                      internal sync diff generation
internal/reconcile/
  service.go                   shared reconciliation action model and service contract
internal/ui/
  tui.go                       terminal sync program entry/update
  model.go                     sync TUI state and pending actions
  view.go                      two-pane layout rendering
  diff.go                      diff cache, scroll, and coloring
  confirm.go                   preflight and confirmed execution command
  keys.go                      key bindings and styles
  update.go                    key handling
internal/process/
  runner.go                    external process execution abstraction
internal/build/
  info.go                      version from runtime/debug and ldflags
```

## Package responsibilities

### `cmd/cm`

Owns only process startup and delegates to `internal/cli`.

### `internal/cli`

Wires Cobra commands, renders command output, and adapts `internal/chezmoi` to
the interfaces consumed by command handlers and sync UIs.

### `internal/chezmoi`

Owns chezmoi CLI execution and status parsing. The `Status` method adds
`--path-style=absolute` so callers always get absolute target paths.
`ManagedFiles` backs `cm edit` completion. `ContentLoader` adapts
chezmoi-rendered target content and local files to `internal/syncdiff`.

### `internal/syncdiff`

Generates sync diffs from rendered chezmoi target content to the current local
file without shelling out to `chezmoi diff` or an external diff renderer.

### `internal/reconcile`

Defines the shared reconciliation action model and review service contract used
by the terminal UI. It keeps sync-domain capabilities out of UI presentation
packages.

### `internal/process`

Owns subprocess execution. Chezmoi, git status, and lazygit all go through this
runner boundary so tests can inject one process fake.

### `internal/ui`

Owns the terminal sync UI. Depends on the `internal/reconcile` service contract
and does not know about chezmoi binary paths or git. The selected file's diff
loads by default on entry and when selection changes; confirmed actions execute
one target at a time so the UI can return to review mode with remaining files.

### `internal/build`

Reads `runtime/debug.ReadBuildInfo()` and the `Version` ldflag variable for
`cm version` output.

## Dependencies

```text
charm.land/bubbletea/v2          terminal UI runtime
charm.land/bubbles/v2            TUI help/key bindings
charm.land/lipgloss/v2           TUI styling
github.com/spf13/cobra           CLI framework
github.com/fatih/color           plain terminal colours
github.com/rogpeppe/go-internal  unified diff implementation
golang.org/x/term                terminal detection
```

No Viper. No external diff renderer. No generated code.
