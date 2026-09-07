# Review: Phase 1 Authoritative Reconciliation

## Scope Reviewed

- App reconciliation domain, status partitioning, authoritative diff, action execution, edit resolution, and semantic reports.
- Chezmoi client status, diff, target metadata, template detection, managed completion, and destination lookup.
- TUI review loading, target-aware help/action gating, fingerprint preflight, postflight retention, notices, and clean transitions.
- CLI wiring, unit regressions, real chezmoi integration harness, CI dependency installation, public docs, and planning artifacts.

## Findings

1. **Critical — custom byte diff misrepresented non-file target state.** Metadata-only drift appeared as `no diff`, symlinks were dereferenced, directories failed, and scripts looked like file deletions. Fixed by forcing bounded chezmoi builtin reverse diff with all entry types included.
2. **Critical — confirmation was not bound to reviewed content.** Dirty-only preflight could execute against changed content. Fixed by fingerprinting status, type, template state, and diff, then recomputing immediately before mutation.
3. **Critical — command success was treated as reconciliation success.** Targets left the list without a clean postcondition. Fixed by postflight review; dirty targets remain and require another decision.
4. **High — scripts entered ordinary file reconciliation.** Apply could execute arbitrary scripts through a file-oriented action label. Fixed by partitioning `R` status rows into `automation:` and excluding them from targetless diff/sync.
5. **High — target/action validity was not represented.** Every row advertised add/apply/merge even when chezmoi would ignore or reject the operation. Fixed with app-owned target types, template detection, conservative action gating, and dynamic TUI help.
6. **High — `cm edit` assumed destination equals `$HOME`.** Fixed with `chezmoi target-path` and absolute-target passthrough.
7. **Medium — successful buffered output was discarded.** Fixed with `ActionResult` notices shown after execution.
8. **Medium — configured chezmoi include/exclude filters could hide review state.** Fixed by explicitly requesting all status/diff/dump entry types and excluding none.
9. **Medium — first-column-only status was treated as destination/target mismatch.** Fixed by deriving reconciliation work from the second status column.
10. **Medium — a target becoming clean while its review loaded could leave the TUI waiting or fail to load the next entry.** Fixed with explicit clean transition and next-review scheduling.
11. **Low — template detection was needlessly run for directories.** Fixed by querying template membership only for file and symlink targets.

## Fixes Applied

- Added `internal/app/reconcile.go` and `internal/chezmoi/target.go`.
- Removed `internal/diff` implementation and `internal/chezmoi/content.go`.
- Migrated the TUI from infrastructure status rows to app-owned entries and reviews.
- Added focused unit/state-machine tests and isolated real-chezmoi integration coverage.
- Added pinned chezmoi installation to CI and removed the old diff dependency.
- Updated README, usage, design, development, changelog, requirements, roadmap, state, context, research, and plan records.

## Waivers

- Context-aware cancellation remains deferred; execution mode intentionally waits for active mutating subprocesses.
- The preflight check is optimistic rather than an atomic filesystem lock. It closes the broad stale-review window but cannot eliminate a change occurring between final review and subprocess read.
- Direct `cm apply` remains a thin chezmoi wrapper and may execute scripts; docs state this explicitly.
- Buffered action output is retained in memory for one target and displayed after execution; a separate bounded-output contract is deferred unless real commands show problematic volume.
- Public repository/module naming and package-manager distribution remain outside this safety phase.

## Remaining Risks

- New or changed chezmoi target types map to `unknown` and remain read-only until integration evidence supports them.
- Template/source metadata calls can evaluate secret-backed state; output is bounded and not logged, but the normal review surface can still display rendered secrets.
- CI installation depends on availability of the pinned chezmoi module through the Go proxy.

## Result

Ready for focused runtime verification and full quality gates.
