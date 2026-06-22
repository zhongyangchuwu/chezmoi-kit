# Plan: 6-app-services

## Objective

Promote application use cases and sync action ownership into `internal/app`, leaving CLI as command wiring and TUI as terminal interaction.

## Scope

In scope:

- Add `app.Services` and `app.NewServices(client chezmoi.Client)` as the service graph entry point.
- Move service interfaces currently private to CLI into app where they represent application capabilities.
- Move `ActionKind`, `Action`, `TerminalCommand`, and sync service contract from TUI to app.
- Move concrete logic from the old CLI service adapter into app service files.
- Simplify CLI tests to verify command wiring and output streams only.
- Move terminal command construction tests into app.

Out of scope:

- Semantic report document model and palette rendering; Phase 7 owns this.
- New user-facing commands.
- Subprocess cancellation API.
- Community mocking frameworks.

## Tasks

1. Create `internal/app/services.go` with `Services`, service interfaces, and `NewServices`.
2. Move status/source logic into app:
   - `Status(targets)`
   - `SourceStatus()`
   - `sourceDir()`
   - `OpenSourceGit()`
3. Move diff command/use-case logic into app:
   - target expansion when no explicit targets are passed
   - per-target `diff.Differ` use
   - clean output behavior retained until Phase 7 replaces rendering.
4. Move sync action domain and execution into app:
   - `ActionKind`, `Action`, `TerminalCommand`, `SyncService`
   - add/apply buffered execution
   - merge terminal command construction
5. Move direct target wrappers and edit behavior into app.
6. Update CLI to depend on `app.Services` and remove concrete chezmoi/git orchestration from CLI.
7. Update TUI to depend on `app.SyncService` and `app.Action` types.
8. Rewrite tests by ownership:
   - app tests cover command construction and process IO.
   - cli tests cover Cobra wiring only.
   - tui tests cover state transitions and service calls only.
9. Update docs and planning/codebase references.
10. Run targeted and full verification.

## Acceptance Criteria

- `internal/cli` contains no concrete chezmoi/git/lazygit/diff engine orchestration.
- `internal/tui` imports app action/service contracts rather than owning them.
- `internal/app` owns command use cases and concrete service implementation.
- Terminal merge command construction is tested in app only.
- CLI/TUI tests no longer need shared terminal command fakes.
- Existing command behavior remains unchanged.

## Verification

Run after implementation:

```bash
go test ./internal/app ./internal/cli ./internal/tui
go test ./...
go vet ./...
go mod tidy -diff
```

Behavior smoke if practical:

```bash
go build -o /tmp/cm-app-services ./cmd/cm
/tmp/cm-app-services status
/tmp/cm-app-services diff
/tmp/cm-app-services version
```

## Risks

- Moving orchestration may accidentally shift error messages or wrapping; app tests should pin important user-facing errors.
- CLI tests may become too shallow if they only assert exit code; keep assertions around the app interface calls and stdout/stderr streams.
- App owning current raw output is an interim state; Phase 7 will replace it with semantic reports.
