# Plan: Phase 2 Doctor Diagnostics

## Objective

Implement `cm doctor` as a read-only diagnostic command for required and optional environment prerequisites.

## Scope

In scope:

- New `cm doctor` command.
- Required checks for `chezmoi`, `git`, and chezmoi source path/repository readiness.
- Optional warning for `lazygit` because only `cm git` needs it.
- Semantic report output using the Phase 1 render controls.
- Tests for pass, warning, failure, and command wiring.
- README documentation after behavior is verified.

Out of scope:

- Installing tools or fixing environment problems automatically.
- Changing existing command error handling.
- Richer status interpretation.
- Mutating git, chezmoi, or local files.

## Tasks

1. Inspect existing `chezmoi.Client`, app service, and process runner capabilities for read-only checks.
2. Add doctor service/report types at the app boundary.
3. Wire `cm doctor` in CLI with output/color rendering.
4. Add tests for required failures, optional warnings, pass output, and no mutation behavior.
5. Update README command table and doctor usage.
6. Run targeted and repository tests; record verification.

## Acceptance Criteria

- `cm doctor` performs only read-only checks.
- Missing `chezmoi` or `git` produces a failed diagnostic and non-zero exit.
- Missing `lazygit` produces a warning but does not fail the whole command by itself.
- Output clearly distinguishes pass, warning, and failure.
- `--output` and `--color` work for doctor reports.
- Existing commands retain current behavior.

## Verification

- Run app/cli/process tests touched by implementation.
- Run `go test ./...`.
- Smoke `go run ./cmd/cm doctor --output markdown` when possible.

## Risks

- Checking source git readiness can accidentally run mutating commands if not limited to read-only process calls.
- Host environment may lack optional tools, so smoke expectations must distinguish warnings from failures.
- Adding process seam methods can broaden abstractions; keep the smallest interface needed.
