# Summary: Phase 1 Authoritative Reconciliation

## Completed Changes

- Replaced cm's byte-oriented custom diff engine with bounded chezmoi builtin reverse diff, explicitly disabling color, pagers, configured external diff commands, and configured type exclusions.
- Added app-owned reconciliation status, entry, target type, review, SHA-256 fingerprint, action, and action-result values.
- Derived reconciliation work from chezmoi's second status column and separated pending scripts into an `automation:` status surface outside `cm sync`.
- Added bounded target metadata parsing, template detection, NUL-safe managed paths, and configured destination lookup to the chezmoi adapter.
- Added conservative target-aware action gating and dynamic TUI help.
- Bound pending actions to reviewed state, deferred changed reviews before mutation, and verified target state after successful execution.
- Preserved successful buffered stdout/stderr as visible TUI notices.
- Fixed relative `cm edit` targets for custom chezmoi destinations while preserving absolute paths.
- Added real chezmoi integration coverage for text direction, mode-only drift, symlinks, directories, removes, scripts, templates, encrypted re-add, and custom-destination editing.
- Added pinned chezmoi release-asset installation to CI and release workflows.
- Removed `internal/diff`, `internal/chezmoi/content.go`, and `github.com/rogpeppe/go-internal`.
- Updated public docs, changelog, planning roots, phase artifacts, and active codebase maps.

## Files Changed

- Runtime: `internal/app/*`, `internal/chezmoi/*`, `internal/tui/*`, `internal/cli/*`.
- Removed: `internal/diff/*`, `internal/chezmoi/content.go`.
- Automation: `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `go.mod`, `go.sum`.
- User docs: `README.md`, `docs/usage.md`, `docs/design.md`, `docs/development.md`, `CHANGELOG.md`.
- Planning: root project/requirements/roadmap/state, refreshed `.planning/codebase/*`, and this phase directory.

## Deviations

- Review found an additional status correctness requirement: first-column-only chezmoi history drift must not appear as destination/target reconciliation work. `AUTH-STATUS-01` was added and implemented.
- `chezmoi dump` does not emit a target entry for deletion effects. Remove reviews now derive `TargetRemove` from second-column `D` and use authoritative diff without dump metadata.
- `managed --recursive=false` is not supported. Template membership uses an exact target comparison against NUL-delimited `managed --include=templates` output, and non-file types skip that query.
- Installing chezmoi v2.72.1 with `go install` is rejected because its module contains exclude directives. CI/release now download the official `chezmoi-linux-amd64` release asset with `gh` instead; the downloaded binary and tests were exercised locally.
- A CLI TUI test initially consumed keystrokes before asynchronous review load. The test now gates input on the review callback rather than sleeping.
- An LSP rename corrupted adjacent syntax; the corruption was repaired immediately and reported to harness QA.

## Evidence

- Focused package tests for app, chezmoi, CLI, and TUI passed after each behavioral cutover.
- Real integration tests passed with installed chezmoi v2.72.1, including builtin age encrypted re-add.
- Isolated CLI status displayed independent `local:`, `automation:`, and `chezmoi:` blocks.
- Isolated targetless diff showed symlink targets, directory modes, file modes/content, and excluded the pending script.
- Actual PTY TUI apply completed and wrote the rendered target.
- Actual PTY stale-review scenario deferred without mutation and refreshed the diff.
- Actual PTY successful no-op scenario retained the dirty target and displayed both warning and postflight reason.
- Final `go test ./...`, race tests, vet, tidy diff, module verification, govulncheck, GoReleaser check, and LSP diagnostics passed.

## Unresolved Risks

- Optimistic preflight cannot eliminate the final filesystem race between fingerprint comparison and the chezmoi subprocess reading state.
- Rendered reviews can expose secret-backed content in terminal scrollback.
- Direct `cm apply` may execute scripts by design.
- Windows real-tool integration is not exercised by the POSIX wrapper harness.
- Remote GitHub Actions and tagged release publication have not been observed for this branch.
