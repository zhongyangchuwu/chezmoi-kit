# Summary: Phase 2 Doctor Diagnostics

## Completed Changes

- Added app-level `DoctorService` and semantic `DoctorReport` generation.
- Added read-only checks for `chezmoi`, `git`, chezmoi source path, source git repository readability, and optional `lazygit`.
- Added `cm doctor` CLI command using the existing report rendering path and Phase 1 output/color flags.
- Added app and CLI tests for pass, warning, failure, and command rendering behavior.
- Updated README command table and optional report-control command list.

## Files Changed

- `internal/app/services.go`
- `internal/app/services_test.go`
- `internal/cli/doctor.go`
- `internal/cli/root.go`
- `internal/cli/cli_test.go`
- `README.md`
- `.planning/ROADMAP.md`
- `.planning/STATE.md`
- `.planning/phases/phase-2-doctor-diagnostics/CONTEXT.md`
- `.planning/phases/phase-2-doctor-diagnostics/PLAN.md`

## Deviations

- Doctor executable checks use `--version` through the existing runner output seam instead of adding a separate look-path abstraction.
- `cm doctor` returns a rendered report and then returns a non-zero error when required checks fail, matching existing CLI error handling.

## Evidence

- `go test ./internal/app ./internal/cli ./internal/report` passed.
- `go run ./cmd/cm doctor --output markdown` printed a Markdown doctor report.
- `go run ./cmd/cm doctor --output plain` printed a plain doctor report.
- `go run ./cmd/cm doctor --color never` printed a plain-colored doctor report.
- `go test ./...` passed.

## Unresolved Risks

- Host smoke observed all tools available; missing-tool behavior is covered by unit tests, not by altering the host PATH.
