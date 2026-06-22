# Verification: Phase 2 Doctor Diagnostics

## Claims Checked

- `cm doctor` exists and is read-only.
- Required checks for `chezmoi`, `git`, source path, and source git repository can fail the command.
- Optional `lazygit` absence is reported as a warning without failing the command by itself.
- Doctor output distinguishes pass, warning, and failure.
- Doctor uses the same `--output` and `--color` rendering controls as other report-backed commands.
- Existing commands retain current behavior under repository test coverage.

## Evidence Observed

- `go test ./internal/app ./internal/cli ./internal/report` passed.
- `go test ./...` passed.
- `go run ./cmd/cm doctor --output markdown` printed:
  - `# **doctor:**`
  - pass rows for `chezmoi`, `git`, `source-path`, `source-git`, and `lazygit` on this host.
- `go run ./cmd/cm doctor --output plain` printed plain pass rows.
- `go run ./cmd/cm doctor --color never` completed successfully without visible ANSI escapes.
- App tests cover required tool failure, source git failure, and optional lazygit warning.
- CLI tests cover doctor report rendering and non-zero result for failed required checks.

## Coverage

- App diagnostic logic: covered by `internal/app` tests.
- CLI command wiring and report rendering: covered by `internal/cli` tests.
- Shared report rendering: covered by `internal/report` tests and Phase 1 tests.
- End-to-end command smoke: covered on the current host for all-pass doctor output.

## Gaps

- Host smoke did not observe a real missing required executable because current environment has required tools.
- Doctor does not attempt automatic remediation; this is intentionally out of scope.

## Result

Passed with missing-tool behavior covered by tests rather than host mutation.
