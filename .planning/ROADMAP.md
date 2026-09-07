# Roadmap: cm

## Overview

The verified `v0.2.0` reconciliation baseline and Phase 2 unified workspace are intact. The remaining development steps make cm a complete chezmoi management entrypoint by integrating mature external tools, then adding contextual management and diagnostics. These are ordered development steps, not promises of releases.

## Version Policy

- Patch releases (`v0.1.Z`) may add optional diagnostics or fix defects without changing the user mental model.
- Minor releases (`v0.Y.0`) carry safety-model, status-model, output-contract, or command-semantics changes.
- `v1.0.0` remains deferred until public command contracts and the safety model are stable.

## Phases

- [x] Phase 1: Authoritative Reconciliation
- [x] Phase 2: Unified File Workspace
- [x] Phase 2.1: Workspace Usability

- [x] Phase 2.2: Directory Browser

- [ ] Phase 3: Mature Tool Handoffs — next step 2
- [ ] Phase 4: Management and Diagnostics — next step 3

## Phase Details

### Phase 1: Authoritative Reconciliation

- Goal: Ensure every previewed and confirmed sync operation reflects chezmoi's actual target type and the exact state the user reviewed.
- Depends on: `v0.1.1` sync TUI, semantic reports, doctor diagnostics, and encrypted `re-add` behavior.
- Requirements: AUTH-DIFF-01, AUTH-STATUS-01, AUTH-TYPE-01, AUTH-SCRIPT-01, AUTH-REVIEW-01, AUTH-VERIFY-01, AUTH-OUTPUT-01, EDIT-DEST-01, INTEGRATION-01, DOC-02.
- Success Criteria:
  - Metadata-only, symlink, directory, regular-file, and template previews come from forced chezmoi builtin diff output with the established target-to-local direction.
  - First-column-only chezmoi history drift does not enter the destination/target reconciliation list.
  - Scripts are reported separately and cannot be selected as ordinary sync actions.
  - Invalid type/action combinations are rejected before confirmation.
  - A target changed after review is not mutated, and a successful command does not remove a target that remains dirty.
  - Relative `cm edit` targets work with a non-home chezmoi destination.
- Plans: 1
  - [x] phase-1-authoritative-reconciliation-01

### Phase 2: Unified File Workspace

- Goal: Find and fully inspect relevant configuration in cm even when it is clean, unmanaged, or ignored.
- Depends on: The verified Phase 1 baseline; preserve it as a reviewable change before starting implementation. Tagging/release is a separate ship decision.
- Requirements: WORKSPACE-01, WORKSPACE-02, WORKSPACE-03, WORKSPACE-04, TEMPLATE-01.
- Scope:
  - Add an explicit `cm ui [path...]` entrypoint. Keep bare `cm` as read-only status and `cm sync [target...]` as focused reconciliation; share the workbench implementation rather than create two TUIs.
  - Browse managed-clean, managed-dirty, unmanaged, and ignored targets. Limit unmanaged discovery to an explicitly selected directory; no automatic home-wide scan or add.
  - Provide tree/flat views, path search, category filters, and clear type/attribute/path information. Chezmoi remains authoritative for membership, ignore rules, source mappings, and attributes.
  - Keep authoritative unified diff as the default; add full-screen preview, horizontal/vertical navigation, search, and hunk jumps so truncated lines are inspectable.
  - Offer opt-in source-template and destination/target content views with explicit labels. Do not automatically reveal decrypted content or template data.
- Success Criteria:
  - A clean managed file remains visible and inspectable in `cm ui`; clean state does not close the workbench.
  - A new file in a selected managed directory appears as an unmanaged candidate without modifying source state; ignored entries are visibly distinguished.
  - Switching tree/flat/filter views retains the same selected target when it remains visible, with a deterministic fallback when it does not.
  - A long line's tail and every diff hunk are reachable; metadata-only and symlink changes remain visible.
  - Bare `cm` and focused `cm sync` retain their existing safety and exit contracts.
- Verification: Real chezmoi fixture with a clean file, dirty file, new file, ignored file, template, symlink, and directory; actual PTY navigation at normal and narrow terminal sizes; cancellation must leave both sides unchanged.
- Not Doing: Builtin editor/merge/Git implementation, automatic filesystem watching, source filename or ignore-rule reimplementation, and inline editable hunks.
- Plans: 1
  - [x] phase-2-unified-workspace-01

### Phase 2.1: Workspace Usability

- Goal: Make the completed read-only workspace easy to scan and operate without external documentation.
- Depends on: Phase 2 workspace inventory, preview, and explicit mode boundaries.
- Requirements: UX-01.
- Scope: Semantic state/type/attribute palette, selection/focus hierarchy, compact legend, responsive footer, preview feedback, and a keyboard-isolated `?` help overlay.
- Not Doing: Themes, persistent preferences, icon fonts, mouse controls, localization, or workspace mutations.
- Plans: 1
  - [x] phase-2.1-workspace-usability-01

### Phase 2.2: Directory Browser

- Goal: Replace recursive disclosure rows with direct-child browsing that gives directory context without alignment ambiguity.
- Depends on: Phase 2 inventory/preview safety and Phase 2.1 semantic presentation.
- Requirements: UX-02; supersedes Phase 2's tree/flat/collapse interaction details.
- Scope: Virtual ancestor directories, Parent/Current/Preview responsive panes, local directory summaries, focus-sensitive `h`/`l` navigation, global path locating, and ancestor-preserving filters.
- Not Doing: File operations, selection mode, tabs, mouse controls, themes, custom layouts, or additional filesystem inventory.
- Plans: 1
  - [x] phase-2.2-directory-browser-01

### Phase 3: Mature Tool Handoffs

- Goal: Complete editing, merging, and source Git work from cm while retaining familiar external tools and returning to a consistent workspace.
- Depends on: Phase 2 file identity, selection, filters, and preview state.
- Requirements: TEMPLATE-02, TEMPLATE-03, TOOLS-01, TOOLS-02, TOOLS-03, TOOLS-04.
- Scope:
  - Provide distinct local-edit and source-edit actions. Use the configured editor; route managed source/template/encrypted edits through `chezmoi edit`.
  - Explain destination/source/rendered-target roles before invoking the configured chezmoi merge tool. After edit or merge, return to refreshed review rather than implicitly apply.
  - Add a Git entry with source-repository summary and launch lazygit in the correct working context. Stage, commit, branches, fetch, pull, push, and conflict resolution remain lazygit's responsibility.
  - Suspend cm rendering while the external tool owns the terminal. Restore cm on success, cancellation, or failure; show actionable missing-tool errors without making optional tools startup requirements.
  - On return, refresh source Git and relevant managed/local state, including membership changes. Restore surviving selection/filter state and retain pending actions only after confirming their reviewed state remains valid.
  - Honor the same chezmoi source/destination/config context across commands. In the workbench source-edit flow, explicitly prevent inherited edit auto-apply/watch settings from bypassing review.
- Success Criteria:
  - Edit a clean template from cm, return to its rendered diff, and explicitly apply the reviewed result without manually finding source paths.
  - Launch merge, return, and see which differences remain; successful merge alone is not labeled fully synchronized.
  - Launch lazygit, perform real source-repository operations, exit, and recover selection and refreshed state in cm.
  - A source change or branch switch in lazygit invalidates incompatible pending actions and handles renamed/deleted selected targets.
  - Missing tools, nonzero exits, resize, and cancellation restore a usable terminal; no Git operation is automatically run on startup or return.
- Verification: Actual PTY round trips with configured editor, merge tool, and lazygit using disposable chezmoi state and a local bare Git remote; exercise source edits, branch changes, failure, and cancellation without production pushes.
- Not Doing: Reimplementing Git/editor/merge features, embedding a continuously running child terminal in a subpane, automatically committing/pulling/pushing, or starting helpers while cm is already mutating a target.
- Plans: TBD — expand the execution plan after Phase 2's state contract is stable.
  - [ ] phase-3-tool-handoffs-01

### Phase 4: Management and Diagnostics

- Goal: Handle routine target lifecycle operations and diagnose failures without leaving cm to remember chezmoi command syntax.
- Depends on: Phase 2's file inventory and Phase 3's handoff/refresh behavior.
- Requirements: MANAGE-01, MANAGE-02, MANAGE-03, TEMPLATE-04, DIAG-01.
- Scope:
  - Add a searchable command palette and contextual menus. Menus and shortcuts call the same app operation paths, with unavailable actions explained rather than dispatched blindly.
  - Add confirmed unmanaged-file add and managed-file forget flows. Preview affected targets and source effects; verify that forget preserves destination content.
  - Delegate private/executable/encrypted attribute changes and template conversion to chezmoi. Preserve current attributes unless an explicit change is selected; distinguish permissions (`private`) from encryption.
  - Permit multi-selection for independent, supported target operations with an exact target list, clear scope, confirmation, per-target outcomes, and no implicit recursive expansion.
  - Add diagnostics for source/destination/config context, tool availability, native chezmoi doctor output, current errors, and session-local last operation results. Sensitive raw config/data/output is opt-in, not automatically copied to timing logs.
  - Show pending scripts as a separate automation inventory; keep script execution outside ordinary file sync.
- Success Criteria:
  - Discover a new file, choose its intended source treatment, review/add it, then open lazygit to commit from the same cm session.
  - Forget a managed file and verify that its destination content remains unchanged and inventory now reflects the new membership.
  - Change a supported attribute or convert a plain managed file to a template through chezmoi, then inspect the new source/target result.
  - A multi-target operation lists the exact affected targets and reports completed, failed, skipped, and remaining work without claiming atomic rollback.
  - Diagnose a missing helper, template error, or permission failure in cm; excluded/uninspected targets are never reported as verified clean.
- Verification: Real chezmoi lifecycle fixtures for add/forget/chattr, custom destination, encrypted state, ignored targets, and partial failure; PTY palette/menu workflows; inspect that diagnostics do not copy sensitive output into timing logs.
- Not Doing: Destroy/purge and bulk script execution, a backup/undo database, generic configuration transformation or secret-redaction engines, and application-syntax validators in this three-step plan.
- Plans: TBD — expand the execution plan before enabling new mutating operations.
  - [ ] phase-4-management-diagnostics-01

## Delivery Gates

- Each step must deliver the complete user flow described above before the next step becomes active; new implementation needs its own concrete phase execution plan.
- Preserve authoritative chezmoi behavior, fingerprint checks, postflight checks, and explicit secret/script boundaries across all adapters.
- Verify interactive work against an actual PTY and real external tools, not only command-construction fakes. Run focused regressions and the existing Go quality gates before claiming implementation complete.
- Keep release numbering, tagging, and remote publication separate from the three-step development sequence.
- The former four template requirements remain tracked: source views belong to Phase 2, edit/merge round trips to Phase 3, and template conversion to Phase 4.

## Progress

| Phase | Plans Complete | Status | Completed |
|---|---:|---|---|
| Phase 1: Authoritative Reconciliation | 1/1 | complete | 2026-09-05 |
| Phase 2: Unified File Workspace | 1/1 | complete | 2026-09-07 |
| Phase 2.1: Workspace Usability | 1/1 | complete | 2026-09-07 |
| Phase 2.2: Directory Browser | 1/1 | complete | 2026-09-07 |
| Phase 3: Mature Tool Handoffs | 0/TBD | planned | - |
| Phase 4: Management and Diagnostics | 0/TBD | planned | - |
