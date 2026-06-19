# Plan: 1-runtime-fixes

## Objective

Fix release-blocking runtime and module issues before documentation cleanup and CI setup.

## Scope

In scope:

- TUI signal handling in `internal/ui/tui.go`.
- Executing-mode key handling/help semantics in `internal/ui`.
- Source git status parser behavior in `internal/cli/service.go` and tests.
- `cm edit` command wiring test coverage in `internal/cli/cli_test.go`.
- Go module tidy state in `go.mod` and `go.sum`.

Out of scope:

- Public documentation rewrites.
- `docs/superpowers/` deletion.
- GitHub Actions workflows.
- New commands such as `cm doctor`.
- Full subprocess cancellation API.

## Tasks

1. Inspect current TUI key/update/help code and tests around `modeExecuting`.
2. Remove `tea.WithoutSignals()` from production TUI program options.
3. Route `modeExecuting` keypresses to a no-op update path and adjust help text so quit is not advertised while executing.
4. Add/adjust UI tests for executing mode quit behavior and existing execute error behavior.
5. Change `parseSourceStatus` to return an error for malformed non-empty lines and update `SourceStatus`/tests.
6. Extend CLI fake service with `EditTarget`/`ManagedFiles` and add `cm edit` command wiring test.
7. Run `go mod tidy` to fix dependency metadata.
8. Run targeted and full verification commands.
9. Write `SUMMARY.md` with changed files, evidence, and unresolved risks.

## Acceptance Criteria

- `cm sync` production TUI no longer disables Bubble Tea signal handling.
- Executing mode does not process or advertise quit as a cancel path.
- Malformed source git status output returns an error with line context.
- `cm edit <target>` command wiring is covered by tests.
- `go mod tidy -diff` exits cleanly.
- Relevant Go tests, race tests, vet, module verify, and build checks pass.

## Verification

Run after implementation:

```bash
go test ./internal/ui
go test ./internal/cli
go test ./...
go test -race ./...
go vet ./...
go mod tidy -diff
go mod verify
go build -o /tmp/cm-phase1-check ./cmd/cm
/tmp/cm-phase1-check --help
/tmp/cm-phase1-check version
```

Manual or environment-limited:

- `cm sync` Ctrl+C cleanup in a real terminal.

## Risks

- Removing `tea.WithoutSignals()` may affect deterministic test behavior if any test runs the full program with synthetic input.
- Executing-mode help text may be shared with review/confirm modes; keep change narrow.
- Strict git porcelain parsing may reject valid but unusual porcelain lines if parser assumptions are too tight; match `git status --porcelain=v1` two-character status plus space shape.
