# Summary: Phase 1 Output Controls

## Completed Changes

- Added CLI render options for `--output plain|ansi|markdown` and `--color auto|always|never`.
- Threaded render options through report-backed commands: root `cm`, `cm status`, `cm diff`, and `cm version`.
- Preserved app service contracts: services still return semantic report documents.
- Updated report color behavior so explicit `ColorAlways` forces ANSI while default `ColorAuto` still respects `NO_COLOR` and TTY detection.
- Added focused CLI and report tests for output selection, color policy, invalid values, and `NO_COLOR` precedence.
- Updated README with optional report controls.

## Files Changed

- `internal/cli/render.go`
- `internal/cli/root.go`
- `internal/cli/status.go`
- `internal/cli/diff_cmd.go`
- `internal/cli/cli_test.go`
- `internal/report/render.go`
- `internal/report/report_test.go`
- `README.md`
- `.planning/ROADMAP.md`
- `.planning/STATE.md`
- `.planning/phases/phase-1-output-controls/CONTEXT.md`
- `.planning/phases/phase-1-output-controls/PLAN.md`

## Deviations

- `--output` and `--color` were implemented as persistent flags. They are meaningful for report-backed commands; completion generation may display inherited flags, but completion output itself is not reformatted.
- `--color always` explicitly overrides `NO_COLOR`. This is documented and covered by report tests.

## Evidence

- `go test ./internal/cli ./internal/report` passed.
- `go run ./cmd/cm version --output markdown` printed Markdown report output.
- `go run ./cmd/cm version --color always` ran successfully.
- `go run ./cmd/cm status --output json` failed with `unsupported output "json"` and exit status 1.
- `go test ./...` passed.

## Unresolved Risks

- No real TTY ANSI smoke was performed; unit coverage proves render policy selection, while terminal display remains covered by existing renderer/TTY boundaries.
