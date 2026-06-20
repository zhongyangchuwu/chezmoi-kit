# Capture: 6-app-services

## Durable Docs Updated

- `docs/design.md` now describes `internal/app` as the owner of the service graph, app use cases, and sync action contracts.
- `docs/development.md` now lists `internal/app/services.go`, `internal/app/services_test.go`, and the reduced CLI/TUI package responsibilities.
- `.planning/codebase/ARCHITECTURE.md`, `STRUCTURE.md`, `CONVENTIONS.md`, `STACK.md`, `CONCERNS.md`, and `MAP.md` now reflect app-owned services and Phase 7 semantic-output risks.

## Planning Records Updated

- `.planning/phases/phase-6-app-services/SUMMARY.md` records implementation changes, files changed, deviations, evidence, and risks.
- `.planning/phases/phase-6-app-services/VERIFICATION.md` records success-criteria evidence and requirement coverage.
- Root planning artifacts were updated to mark Phase 6 complete and advance current focus to Phase 7.

## Learnings

- Moving raw status/diff byte generation into app was the smallest way to satisfy ARCH-03 because CLI then only writes app output to streams.
- Keeping `internal/cli/service.go` as a type alias avoided duplicating service graph names in CLI while eliminating CLI-owned orchestration.
- App service tests are the right home for command construction assertions; CLI tests should assert command routing and stream behavior only.
- TUI remains easier to test when it depends on `app.SyncService` and records only state transitions plus service calls.

## Ship Inputs

- Phase 7 should replace app-owned raw output bytes with semantic report documents and renderers.
- Existing external release gates remain: observed remote CI, explicit tag approval, and real-terminal smoke where available.
