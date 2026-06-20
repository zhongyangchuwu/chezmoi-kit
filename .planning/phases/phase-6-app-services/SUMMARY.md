# Summary: 6-app-services

## Completed Changes

- Added `internal/app/services.go` as the app service graph and default `app.NewServices(client chezmoi.Client)` entry point.
- Moved status, source git, diff, sync execution, terminal merge command, direct target wrappers, edit, and managed-file completion use cases into `internal/app`.
- Moved sync action domain contracts to app: `ActionKind`, `Action`, `TerminalCommand`, and `SyncService`.
- Reduced `internal/cli` to Cobra command wiring, stream plumbing, completion, version output, and exit behavior; `internal/cli/service.go` now aliases `app.Services` only.
- Updated `internal/tui` to consume `app.SyncService`, `app.Action`, `app.ActionKind`, and `app.TerminalCommand` directly.
- Moved concrete command construction tests from `internal/cli/service_test.go` to `internal/app/services_test.go` and kept CLI/TUI tests at adapter/state boundaries.
- Updated public and planning documentation to reflect app-owned services and Phase 7 remaining semantic-output work.

## Files Changed

- `internal/app/services.go`
- `internal/app/services_test.go`
- `internal/cli/service.go`
- `internal/cli/root.go`
- `internal/cli/status.go`
- `internal/cli/diff_cmd.go`
- `internal/cli/cli_test.go`
- `internal/tui/model.go`
- `internal/tui/confirm.go`
- `internal/tui/tui.go`
- `internal/tui/update.go`
- `internal/tui/sync_test.go`
- Deleted `internal/cli/service_test.go`
- Deleted `internal/tui/service.go`
- `docs/design.md`
- `docs/development.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/STRUCTURE.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/STACK.md`
- `.planning/codebase/CONCERNS.md`
- `.planning/codebase/MAP.md`

## Deviations

- The plan said CLI status/diff renderers would keep formatting until Phase 7. Implementation moved raw status/diff byte rendering into `internal/app` so `internal/cli` owns only stream writes. This is intentionally still a raw-byte interim model; Phase 7 remains responsible for semantic reports and palette renderers.
- `internal/cli/service.go` remains as a tiny type alias file so existing command constructor signatures stay simple while concrete orchestration is app-owned.

## Evidence

- `go test ./internal/app ./internal/cli ./internal/tui` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go mod tidy -diff` produced no diff.
- Go workspace diagnostics reported no issues.
- `go build -o /tmp/cm-app-services ./cmd/cm && /tmp/cm-app-services version` built and printed version metadata.
- Searches found no `chezmoiService`, `tui.Action`, `tui.ReviewService`, `internal/tui/service.go`, CLI `internal/diff`/`internal/process` imports, or CLI source-git orchestration remnants.

## Unresolved Risks

- Status and diff outputs are still raw bytes owned by app as an interim state; Phase 7 must replace them with semantic report documents and renderers.
- Real-terminal smoke for `cm sync` Ctrl+C cleanup and `cm git` lazygit startup remains an external release gate from earlier phases.
