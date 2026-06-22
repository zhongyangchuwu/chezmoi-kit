# Context: Phase 2 Doctor Diagnostics

## Goal

Add a read-only `cm doctor` command that explains whether the current environment is ready for existing `cm` workflows.

## Constraints

- Do not mutate local files, chezmoi source files, or git state.
- Do not change existing command behavior or error paths outside explicit `cm doctor` invocation.
- Required checks should fail the doctor command when core workflows cannot run.
- Optional checks should report warnings, not universal failure.
- Reuse semantic reports and Phase 1 output/color controls.

## Decisions

- `cm doctor` belongs in app services so CLI remains command wiring and report rendering only.
- External executable checks should use the existing process boundary where possible.
- Doctor output is a report document with pass/warn/fail semantics, not ad hoc text.

## Open Questions

- Existing process runner may need a small read-only lookup method or app may use `exec.LookPath` directly if that keeps the seam smaller.
- Chezmoi source path availability should follow existing client methods discovered during implementation.

## Verification Expectations

- Unit tests cover required tool missing, optional tool missing, all checks passing, and no mutation commands.
- CLI tests cover command wiring and output controls.
- Smoke `go run ./cmd/cm doctor --output markdown` where host tools allow it.
