# Verification: 6-app-services

## Claims Checked

1. `internal/app` exposes the service graph and owns status, diff, sync, source git, targets, edit, options, and version metadata.
2. `internal/cli` contains no concrete chezmoi/git/lazygit/diff engine orchestration.
3. `internal/tui` depends on app sync action/service contracts instead of owning them.
4. Terminal merge command construction is tested once in app, not duplicated across CLI/TUI tests.
5. Existing command behavior is preserved under targeted and full checks.

## Evidence Observed

- `internal/app/services.go` defines `Services`, `NewServices`, app service interfaces, `ActionKind`, `Action`, `TerminalCommand`, and `SyncService`.
- `internal/cli/service.go` is only `type commandServices = app.Services`.
- `internal/tui/model.go`, `internal/tui/confirm.go`, `internal/tui/tui.go`, `internal/tui/update.go`, and `internal/tui/sync_test.go` use `internal/app` contracts directly.
- `internal/app/services_test.go` covers source git status, lazygit source directory, diff target expansion, add/apply buffered execution, terminal merge command construction, edit target expansion, and direct target wrappers.
- `internal/cli/service_test.go` and `internal/tui/service.go` were removed.
- Search checks found no `chezmoiService`, no `tui.Action`/`tui.ReviewService`, no `internal/tui/service.go`, and no CLI `internal/diff`/`internal/process` orchestration imports or source-git orchestration remnants.
- `go test ./internal/app ./internal/cli ./internal/tui` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go mod tidy -diff` produced no diff.
- Go workspace diagnostics reported no issues.
- `go build -o /tmp/cm-app-services ./cmd/cm && /tmp/cm-app-services version` built and printed version metadata.

## Coverage

- ARCH-02 covered by `internal/app/services.go`, `internal/app/services_test.go`, and app docs updates.
- ARCH-03 covered by `internal/cli` simplification, search checks, and CLI tests focused on command wiring and stream behavior.
- TEST-01 covered by app-owned command construction tests and removal of CLI/TUI ownership of concrete terminal command construction.
- Phase 6 success criteria 1-5 all have observed evidence.

## Gaps

- Raw status/diff bytes remain the interim app-owned output representation; semantic report documents and palette renderers are intentionally deferred to Phase 7.
- Real-terminal external release gates from earlier phases remain: `cm sync` Ctrl+C cleanup and `cm git` lazygit startup in an actual terminal.

## Result

Pass. Phase 6 App Services and Test Ownership is verified against its success criteria and mapped requirements.
