# Codebase Structure

**Mapped:** 2026-09-07

## Top Level

- `cmd/cm/` — thin process entrypoint.
- `internal/` — private runtime packages.
- `docs/` — usage, design, and development documentation.
- `.planning/` — active workflow state, codebase maps, phase artifacts, and release archives.
- `.github/workflows/` — CI and tagged release automation.
- `go.mod`, `go.sum` — Go module and reproducible dependencies.
- `.goreleaser.yaml`, `justfile` — release and local build/install helpers.

## Runtime Packages

### `cmd/cm`

- `main.go` passes process arguments and streams to `cli.Main` and exits with its status.

### `internal/cli`

- `root.go` builds the Cobra command tree and wires app services/TUI startup.
- `render.go` selects plain, ANSI, or Markdown semantic report renderers.
- `status.go`, `diff_cmd.go`, `doctor.go` are stream adapters.
- `service.go` aliases the app service graph for command wiring.
- `cli_test.go` covers command behavior and exit contracts.

### `internal/app`

- `services.go` builds services and owns status, diff, doctor, source git, edit, direct target commands, sync execution, and workspace composition.
- `reconcile.go` defines `SyncStatus`, `ReconcileEntry`, `Review`, `TargetType`, `Action`, fingerprints, action gating, and `ActionResult`.
- `workspace.go` defines workspace entries, snapshots, previews, and the app boundary.
- `workspace_service.go` merges bounded managed/ignored/scoped-unmanaged inventory and builds lazily loaded previews.
- `options.go` owns process-wide CLI options.
- `version.go` owns build/runtime metadata and semantic version reports.
- `services_integration_test.go` and `workspace_integration_test.go` exercise isolated real chezmoi behavior.

### `internal/chezmoi`

- `client.go` owns generic command execution, bounded output, buffered execution, status, and managed completion.
- `workspace.go` owns strict managed path mappings, typed membership, ignored/unmanaged lists, target/source content adapters.
- `status.go` strictly parses status lines and NUL-delimited paths.
- `target.go` owns forced builtin reverse diff, bounded target metadata, template membership, and destination lookup.
- Tests live beside each parser/client surface.

### `internal/process`

- `runner.go` is the only runtime `os/exec` abstraction.

### `internal/report`

- `report.go` defines semantic documents and inline roles.
- `render.go` renders plain, ANSI, and Markdown output.
- `diff.go` classifies unified diff lines for CLI and TUI styling.

### `internal/tui`

- `tui.go` starts Bubble Tea sync/workspace programs and routes messages.
- `model.go` owns shared entries, sync pending actions, focus/mode state, workspace preview state, and workspace help visibility.
- `workspace_list.go` projects workspace tree/flat/filter/search views with selection preservation.
- `workspace_help.go` renders responsive quick-start, key, and legend guidance; `styles.go` centralizes semantic color and no-color visual policy.
- `diff_state.go` lazily loads and caches sync reviews/workspace previews with scrolling and stale-result isolation.
- `confirm.go` owns sync-only preflight, sequential execution, terminal handoff, postflight, and completion transitions.
- `view.go`, `diff_view.go`, `keys.go`, `update.go`, `path.go`, `timing.go` own rendering, input, path display, and bounded diagnostics.
- `sync_test.go` and `workspace_test.go` cover target-aware sync safety and workspace interaction behavior.

## Removed Boundaries

- `internal/diff` — removed; authoritative diff now comes from forced chezmoi builtin output.
- `internal/chezmoi/content.go` — removed; cm no longer reconstructs target state from destination byte reads.

## Entrypoints

- Process: `cmd/cm/main.go`
- CLI: `internal/cli.Main`
- Services: `internal/app.NewServices`
- Sync UI: `internal/tui.RunSyncTUI`
- Workspace UI: `internal/tui.RunWorkspaceTUI`
- Reports: `internal/report.Plain`, `ANSI`, `Markdown`
