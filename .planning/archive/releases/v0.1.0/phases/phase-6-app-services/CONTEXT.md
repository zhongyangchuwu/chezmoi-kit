# Context: Phase 6 App Services and Test Ownership

## Goal

Move application behavior out of CLI/TUI packages into `internal/app` so command code is a thin adapter and tests align with ownership.

## Constraints

- Preserve observable command behavior and current sync execution semantics.
- Do not create one adapter per method. Prefer cohesive app services and a single service graph constructor.
- `internal/cli` must not own chezmoi/git/lazygit orchestration after this phase.
- `internal/tui` must not own sync action domain types after this phase.
- Avoid a generic `util` or `testutil` package.

## Decisions

- `internal/app` owns major services: status, diff, sync, targets, edit, source git, options, and version metadata.
- Sync action types (`ActionKind`, `Action`, `TerminalCommand`, `SyncService`) belong to app, not TUI.
- CLI tests should verify Cobra wiring and stream behavior only; app tests verify concrete command construction and terminal command behavior.
- TUI tests should verify state transitions and calls to the app sync interface, not chezmoi command arguments.

## Open Questions

- Whether app services should expose reports directly in this phase or wait for Phase 7. Default: wait; keep current byte/text behavior until report model exists.
- Whether direct target wrappers and sync non-interactive execution should share one app helper or stay separate for clearer command semantics.

## Verification Expectations

- App service tests cover source git status, lazygit source dir, diff target expansion, sync add/apply buffering, terminal merge command construction, edit target expansion, and direct target wrappers.
- CLI tests no longer fake terminal command construction.
- TUI tests use local minimal fakes only for state transitions.
- Full Go test/vet/tidy checks pass.
