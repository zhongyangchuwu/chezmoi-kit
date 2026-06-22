# Summary: Phase 1 Runtime Stability

## Completed Changes

- Removed production `tea.WithoutSignals()` from `internal/ui/tui.go` so Bubble Tea can own normal signal cleanup.
- Routed `modeExecuting` keypresses through a no-op update path instead of review-mode quit handling.
- Removed executing-mode quit help text from `internal/ui/view.go` and made `executingHelp()` return no advertised key bindings.
- Made source git status parsing strict: malformed non-empty `git status --porcelain=v1` lines now return an error with line context.
- Added CLI command wiring coverage for `cm edit <target>` and completed the CLI fake service's `editService` implementation.
- Added tests for malformed source git status parsing and executing-mode quit/help behavior.
- Ran `go mod tidy`, moving `github.com/rogpeppe/go-internal` into the direct dependency block and adding missing `go.sum` checksums.

## Files Changed

- `internal/ui/tui.go`
- `internal/ui/update.go`
- `internal/ui/keys.go`
- `internal/ui/view.go`
- `internal/ui/sync_test.go`
- `internal/cli/service.go`
- `internal/cli/service_test.go`
- `internal/cli/cli_test.go`
- `go.mod`
- `go.sum`
- `.planning/` initialization artifacts

## Deviations

- Full subprocess cancellation was not implemented. Phase 1 intentionally uses conservative executing-mode semantics for v0.1.0: once execution starts, the UI does not advertise or process quit as cancellation.
- Manual Ctrl+C terminal cleanup could not be proven in the non-interactive harness; the production code path now delegates signals to Bubble Tea, and manual smoke remains in Phase 4.

## Evidence

- `go test ./internal/ui` passed.
- `go test ./internal/cli` passed.
- `go mod tidy -diff` passed with no output.
- Go workspace diagnostics reported no issues.
- `go test ./...` passed.
- `go test -race ./...` passed.
- `go vet ./...` passed with no output.
- `go mod verify` passed.
- `go build -o /tmp/cm-phase1-check ./cmd/cm && /tmp/cm-phase1-check --help && /tmp/cm-phase1-check version` passed.

## Unresolved Risks

- Manual TUI Ctrl+C terminal cleanup still needs an interactive terminal smoke test.
- `cm sync` executing mode still waits synchronously for long-running chezmoi subprocesses; this is accepted for v0.1.0 but should be revisited if cancellation becomes a requirement.
