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
go run ./cmd/cm ui ~/.config
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
go test ./internal/app
go test ./internal/cli
go test ./internal/chezmoi
go test ./internal/tui
```

`internal/app` integration tests use a real `chezmoi` binary with temporary
source, destination, cache, and persistent state. They skip only when chezmoi is
not installed and never touch the user's configured source or destination.

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
go build -trimpath -ldflags "-s -w -X github.com/zhongyangchuwu/cm/internal/app.Version=$VERSION" -o dist/cm ./cmd/cm
```

For release artifacts, use GoReleaser:

```bash
go run github.com/goreleaser/goreleaser/v2@latest check
go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean
```

GoReleaser reads `.goreleaser.yaml`, builds Linux/macOS/Windows archives, injects
`internal/app.Version={{ .Version }}`, and writes checksums under `dist/`.

## Local release gate

Run before tagging a release:

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
cm ui
cm ui <explicit-directory-scope>
cm sync
cm edit <known-managed-file>
cm git
cm completion zsh
```


For `cm sync`, verify:

- Ctrl+C restores the terminal before execution.
- Quitting before confirmation does not execute actions.
- A target changed after review is deferred and refreshed.
- A successful no-op remains in review after postflight.
- Successful warnings are visible.
- Scripts are reported outside ordinary sync actions.
- Execution mode does not advertise `q` as cancellation.
- Errors from chezmoi commands are returned to the CLI.

For `cm ui`, verify both unscoped and explicit-scope inventory, tree/flat
projection, filters, path and preview search, hunk and horizontal navigation,
full-screen/narrow rendering, template/encrypted reveal, semantic state markers,
the `?` help overlay, `NO_COLOR=1` text fallback, and quit-without-mutation behavior.


To collect sync phase timings during manual smoke:

```bash
cm sync --debug
```

Inspect the temporary log path printed to stderr for `sync initial status`,
`sync review`, `sync review preflight`, `sync execute target`, `sync review
postflight`, and `sync execution finish` entries. Diff and subprocess contents
must not appear in the log.

## GitHub workflows

`.github/workflows/ci.yml` runs on pushes to `main` and pull requests. It installs
pinned chezmoi v2.72.1 for isolated integration coverage, then runs module tidy
checks, module verification, tests, race tests, vet, govulncheck, GoReleaser
configuration validation, and a snapshot build.

`.github/workflows/release.yml` runs on `v*` tags. It installs the same pinned
chezmoi release asset before GoReleaser's test hook, then uses
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
internal/app/
  options.go                    process-wide CLI options
  version.go                    version from runtime/debug and ldflags
  reconcile.go                  sync status, review, target type, fingerprint, action result
  workspace.go                  workspace entries, snapshots, previews, service contract
  workspace_service.go          managed/scoped-unmanaged inventory and bounded previews
  services.go                   service graph and application use cases
  workspace_integration_test.go isolated real-chezmoi workspace behavior
internal/report/
  report.go                     semantic document model
  render.go                     plain, ANSI, and Markdown renderers
  diff.go                       shared diff line classification
internal/cli/
  root.go                       Cobra command wiring
  service.go                    alias to app service graph for command wiring
  status.go                     status output stream plumbing
  diff_cmd.go                   diff output stream plumbing
  cli_test.go                   command wiring tests
internal/chezmoi/
  client.go                     generic chezmoi execution and bounded output
  workspace.go                  inventory/path/content adapter
  status.go                     strict status and NUL-path parsing
internal/tui/
  tui.go                        terminal sync and workspace program entry/update
  model.go                      shared typed entries and mode-specific state
  workspace_list.go             workspace tree/flat/filter/search projections
  view.go                       two-pane and full-screen workspace rendering
  diff_state.go                 review/preview cache, scrolling, and splitting
  diff_view.go                  diff line styling
  confirm.go                    sync-only preflight, execution, and postflight
  keys.go                       key bindings and styles
  update.go                     key handling
  workspace_test.go             workspace presentation behavior
internal/process/
  runner.go                     external process execution abstraction
```

## Package responsibilities

### `cmd/cm`

Owns only process startup and delegates to `internal/cli`.

### `internal/app`

Owns process-wide options, version metadata, the application service graph,
semantic status/diff/version reports, app-owned reconciliation entries and
reviews, action gating, reviewed fingerprints, direct wrappers, and postflight
contracts consumed by the TUI. Release builds override `app.Version` with ldflags.

### `internal/cli`

Wires Cobra commands, process streams, shell completion, and exit behavior. It
delegates status, diff, sync, workspace, source git, edit, and mutating command
behavior to `internal/app` services.

### `internal/chezmoi`

Owns chezmoi CLI execution and strict output parsing. Status uses absolute target
paths, and workspace inventory uses managed path mappings, typed membership,
source-ignored entries, and scoped unmanaged candidates. Authoritative diff
forces chezmoi's builtin renderer with reverse direction and no pager; target
and decrypted-source previews are bounded. `target-path` supplies the configured
destination directory.

### `internal/report`

Owns semantic command-output documents, inline roles, diff line classification,
and plain/ANSI/Markdown renderers. ANSI rendering uses semantic palette roles,
auto-detects TTY stdout, and disables color when `NO_COLOR` is non-empty.


### `internal/tui`

Owns the terminal sync and workspace UI. Sync consumes app-owned reviews,
invalidates stale pending actions, executes confirmed targets sequentially, and
removes a target only after app postflight state reports it clean. Workspace
mode is read-only, persists clean entries, and only loads sensitive content
after an explicit reveal.

### `internal/process`

Owns subprocess execution. Chezmoi, git status, and lazygit all go through this
runner boundary so tests can inject one process fake.
During sync execution, add/apply subprocesses use buffered stdout/stderr and nil
stdin to avoid sharing the active TUI terminal. Successful output is returned to
the TUI; merge remains terminal-bound for interactive merge tools.

## Dependencies

```text
charm.land/bubbletea/v2          terminal UI runtime
charm.land/bubbles/v2            TUI help/key bindings
charm.land/lipgloss/v2           TUI styling
github.com/spf13/cobra           CLI framework
golang.org/x/term                terminal detection for CLI reports and TUI
```

No Viper. No generated code. No user-configured external diff process in `cm diff` or `cm sync`; chezmoi's builtin diff is forced.
