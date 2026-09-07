# Plan: Phase 1 Authoritative Reconciliation

## Objective

Replace regular-file-only preview assumptions with a chezmoi-authoritative, target-aware review and execution model that refuses stale actions and verifies postconditions before declaring reconciliation complete.

## Scope

- App-owned sync status, entry, review, target type, action fingerprint, and action result contracts.
- Chezmoi client support for reconciliation status partitioning, forced builtin reverse diff, bounded target metadata, template detection, NUL-delimited managed paths, and destination lookup.
- TUI review-state loading, valid-action gating, script exclusion messaging, stale preflight, postflight retention, and successful-output visibility.
- Custom destination edit correction.
- Removal of obsolete custom content diff implementation and dependency.
- Behavior-focused unit and real chezmoi integration coverage.
- Public docs, changelog, planning summary, review, verification, and capture updates.

## Tasks

1. Add chezmoi parsers and client methods for NUL paths, target metadata, template detection, target directory, and forced builtin reverse diff.
2. Introduce app reconciliation domain values and partition raw status into ordinary entries and pending scripts.
3. Replace `internal/diff` use with bounded authoritative diff output and remove obsolete custom diff/content code.
4. Migrate TUI state from `chezmoi.StatusEntry` and cached raw diffs to app-owned entries and reviewed snapshots.
5. Gate actions by type/template state and make pending actions carry the reviewed fingerprint.
6. Recompute review before execution; defer stale actions without mutation.
7. Capture successful subprocess output, run postflight status, and retain unresolved targets.
8. Resolve `cm edit` relative targets against `chezmoi target-path`; keep absolute targets unchanged.
9. Add focused regressions and real chezmoi integration scenarios for target semantics and safety transitions.
10. Update README, usage, design, development, changelog, and phase records.
11. Run review and full verification gates; revise until all committed requirements pass.

## Acceptance Criteria

- `cm diff` shows mode-only and symlink changes that the old byte diff missed or misrepresented.
- Directory changes no longer fail because cm attempts to read a directory as a file.
- `cm status` distinguishes pending scripts from ordinary local reconciliation entries.
- A status row whose second column is blank does not enter `local:` or `cm sync`.
- `cm sync` never offers or executes ordinary add/apply/merge actions for scripts.
- A pending action becomes invalid if its reviewed status, type, template state, or diff changes before execution.
- A command that exits successfully but leaves the target dirty does not remove the target or produce a false completion message.
- Successful buffered warnings/output are shown in the TUI completion or remaining-work message.
- Template targets cannot select local-to-source re-add; regular file, template, symlink, directory, remove, and unknown action matrices match the phase context.
- Relative `cm edit` works when chezmoi destination differs from `$HOME`; absolute targets are not rewritten.
- No production caller references `internal/diff` or `chezmoi.ContentLoader` after cutover.
- Required focused, integration, race, vet, module, and smoke checks pass.

## Verification

- Focused package tests for `internal/chezmoi`, `internal/app`, `internal/tui`, and `internal/cli`.
- Real-tool integration tests using temporary source/destination/cache/state and the installed or CI-provided chezmoi binary.
- CLI smoke for status script separation and authoritative diff of mode/symlink/directory changes.
- Isolated encrypted re-add workflow confirming the encrypted source attribute remains intact.
- `gofmt` on changed Go files.
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go mod tidy -diff`
- `go mod verify`

## Risks

- Target metadata commands can render secret-backed state; keep calls per selected target, bounded, and out of logs.
- Builtin diff output headers differ from the old custom headers; docs and tests must treat the authoritative output as the new `v0.2.0` contract.
- Real chezmoi output can evolve; pin CI integration coverage to a known version while keeping parsers strict enough to fail visibly.
- Conservative action gating may expose fewer shortcuts than chezmoi itself; explicit safety is preferred until behavior is proven.
