# Context: Phase 5 Package Architecture Cleanup

## Goal

Make package names and file boundaries express one architecture language before moving more behavior: CLI/TUI adapters, app use cases, and infrastructure capabilities.

## Constraints

- Preserve existing user behavior for `cm status`, `cm diff`, `cm sync`, direct wrappers, completion, and version.
- Do not introduce a generic `util` package.
- Do not keep packages whose contents do not earn a package boundary (`build`, `testutil`).
- Avoid moving business behavior and renderer abstractions in the same step; this phase is mostly naming, ownership, and file layout.
- Keep imports acyclic: `cli` may call `tui` and `app`; `tui` may call `app`; `app` must not call `cli` or `tui`.

## Decisions

- Rename `internal/ui` to `internal/tui` because it is Bubble Tea terminal UI, not a generic UI layer.
- Merge `internal/build` into `internal/app` as version metadata because it is application runtime/version information, not build orchestration.
- Delete `internal/testutil`; duplicated test fakes should disappear through ownership changes or remain local to the owning package test.
- Keep `internal/diff`, `internal/process`, and `internal/chezmoi` as infrastructure/capability packages.
- Split mixed-purpose TUI files by role, especially diff state/loading versus diff presentation.

## Open Questions

- Whether `internal/app` should expose a single `Services` struct immediately in this phase or wait for Phase 6.
- Whether version API names should remain `Info`/`Current` or become `VersionInfo`/`CurrentVersion` during the move.

## Verification Expectations

- `go test ./internal/cli ./internal/tui ./internal/app ./internal/chezmoi ./internal/diff` passes.
- `go test ./...`, `go vet ./...`, and `go mod tidy -diff` pass.
- No imports reference `internal/ui`, `internal/build`, or `internal/testutil` after the phase.
