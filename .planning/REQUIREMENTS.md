# Requirements: cm

**Defined:** 2026-06-19
**Updated:** 2026-09-06
**Core Value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Current Requirements

### v0.2.0 — Authoritative Reconciliation

- [x] AUTH-DIFF-01 — `cm diff` and sync use bounded, forced chezmoi builtin reverse diff output with all target types included, so content, mode, symlink, directory, and other target-state changes match chezmoi semantics.
- [x] AUTH-STATUS-01 — Reconciliation work is derived from chezmoi's second status column; first-column-only history drift is not presented as a destination/target mismatch.
- [x] AUTH-TYPE-01 — The app layer owns reconciliation entry and review types, and the TUI offers only actions valid for the reviewed target type and template state.
- [x] AUTH-SCRIPT-01 — Pending chezmoi scripts are reported separately and are never executed through ordinary `cm sync` file actions.
- [x] AUTH-REVIEW-01 — Every pending sync action records the reviewed fingerprint; changed targets are rejected and refreshed before mutation.
- [x] AUTH-VERIFY-01 — A successful subprocess removes a target only after post-execution status confirms it is reconciled; unresolved targets return to review.
- [x] AUTH-OUTPUT-01 — Successful non-interactive action output and warnings remain visible to the user instead of being silently discarded.
- [x] EDIT-DEST-01 — `cm edit` resolves relative targets against chezmoi's configured destination directory and accepts absolute targets unchanged.
- [x] INTEGRATION-01 — Real chezmoi integration coverage exercises regular content, metadata-only drift, symlinks, directories, script exclusion, templates, encrypted re-add, and custom destination editing where practical.
- [x] DOC-02 — README, usage, design, development, and changelog describe the authoritative diff, script boundary, stale-review protection, postflight behavior, and custom destination support.

## Planned Requirements — Next Three Steps

These requirements define the requested roadmap, not implemented behavior. Release numbers and detailed execution plans remain separate decisions. The four prior template requirement IDs are retained and reassigned rather than dropped.

### Step 1 / Phase 2 — Unified File Workspace

- [x] WORKSPACE-01 — An explicit `cm ui [path...]` opens a persistent management workspace without changing bare `cm` status or focused `cm sync` contracts.
- [x] WORKSPACE-02 — Users can browse managed-clean, managed-dirty, scoped unmanaged candidates, and ignored entries using chezmoi's own classification.
- [x] WORKSPACE-03 — Tree/flat views, path search, and category filters preserve a surviving selected target and use a predictable fallback when it disappears.
- [x] WORKSPACE-04 — Users can inspect complete bounded previews through full-screen mode, horizontal/vertical navigation, search, and hunk jumps; metadata and symlink changes remain visible.
- [x] TEMPLATE-01 — Users can explicitly switch between rendered target diff, destination content, and source-template detail without exposing template data or decrypted content automatically.

### Phase 2.1 — Workspace Usability

- [x] UX-01 — `cm ui` distinguishes file state/type/attributes with semantic color and retained textual markers, provides responsive contextual guidance, and offers keyboard-isolated in-TUI help with a complete legend.

### Step 2 / Phase 3 — Mature Tool Handoffs

- [ ] TEMPLATE-02 — Managed source/template edits delegate to `chezmoi edit`, return to refreshed review, and do not implicitly apply or watch through inherited editor configuration.
- [ ] TEMPLATE-03 — Merge launch explains destination/source/target roles and delegates to the configured merge tool; return shows remaining differences before any explicit apply.
- [ ] TOOLS-01 — A distinct local-edit action launches the configured editor for destination content rather than silently editing source state.
- [ ] TOOLS-02 — Users open lazygit in the source-repository context from cm and return afterward; cm does not reimplement Git mutations.
- [ ] TOOLS-03 — External-tool success, cancellation, and failure restore the cm terminal, surviving selection, filters, and view position; missing optional tools do not block startup.
- [ ] TOOLS-04 — Returning from an external tool refreshes source Git and relevant managed/local state, revalidates pending actions, and handles changed membership and renamed/deleted targets.

### Step 3 / Phase 4 — Management and Diagnostics

- [ ] MANAGE-01 — A searchable command palette and contextual menus reuse the same app operations as shortcuts and explain unavailable actions.
- [ ] MANAGE-02 — Users review and confirm adding unmanaged files or forgetting managed files through chezmoi; forget preserves destination content.
- [ ] MANAGE-03 — Supported private/executable/encrypted attribute changes and independent multi-target operations show exact scope, preserve attributes unless explicitly changed, and report per-target outcomes.
- [ ] TEMPLATE-04 — Existing files convert to templates through an explicit chezmoi-backed operation such as `chattr +template`, followed by refreshed source/target inspection.
- [ ] DIAG-01 — A diagnostic view shows active context, tool availability, native doctor results, errors, and session-local operation results; sensitive details are opt-in and skipped targets are not reported as verified clean.

## Deferred

- Context-aware cancellation of already-started mutating subprocesses.
- Embedded child-terminal subpanes; first use suspend/run/restore handoffs.
- An in-process Git client, text editor, merge engine, or interactive patch editor.
- Destroy/purge workflows and script execution UI; pending scripts remain visible outside file sync.
- Automatic backup/undo storage, generic configuration transformations, application-syntax validators, and share-oriented secret filtering.
- JSON status output and drift-sensitive exit codes.
- GoReleaser signing, notarization, package-manager publishing, and Docker distribution.
- Issue templates, PR templates, and security policy files.
- Template execute/debug commands beyond opt-in inspection and delegated editing.

## Out of Scope

| Feature | Reason |
|---|---|
| Replacing chezmoi | `cm` delegates target-state calculation and mutation to chezmoi. |
| Automatic git commit/push/pull | Explicit Git work is delegated to lazygit; cm does not perform Git mutations on startup or return. |
| Daemon/watch mode | Background mutation increases surprise. |
| Persistent state database | Chezmoi remains authoritative; reviewed fingerprints are session state only. |
| Ordinary sync execution of scripts | Pending scripts stay visible, but script execution is not a file reconciliation action. |

## Traceability

| Requirement | Phase | Status |
|---|---|---|
| AUTH-DIFF-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-STATUS-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-TYPE-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-SCRIPT-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-REVIEW-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-VERIFY-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| AUTH-OUTPUT-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| EDIT-DEST-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| INTEGRATION-01 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| DOC-02 | v0.2.0 Phase 1 — Authoritative Reconciliation | Complete |
| WORKSPACE-01 | Phase 2 — Unified File Workspace | Complete |
| WORKSPACE-02 | Phase 2 — Unified File Workspace | Complete |
| WORKSPACE-03 | Phase 2 — Unified File Workspace | Complete |
| WORKSPACE-04 | Phase 2 — Unified File Workspace | Complete |
| TEMPLATE-01 | Phase 2 — Unified File Workspace | Complete |
| UX-01 | Phase 2.1 — Workspace Usability | Complete |
| TEMPLATE-02 | Phase 3 — Mature Tool Handoffs | Planned |
| TEMPLATE-03 | Phase 3 — Mature Tool Handoffs | Planned |
| TOOLS-01 | Phase 3 — Mature Tool Handoffs | Planned |
| TOOLS-02 | Phase 3 — Mature Tool Handoffs | Planned |
| TOOLS-03 | Phase 3 — Mature Tool Handoffs | Planned |
| TOOLS-04 | Phase 3 — Mature Tool Handoffs | Planned |
| MANAGE-01 | Phase 4 — Management and Diagnostics | Planned |
| MANAGE-02 | Phase 4 — Management and Diagnostics | Planned |
| MANAGE-03 | Phase 4 — Management and Diagnostics | Planned |
| TEMPLATE-04 | Phase 4 — Management and Diagnostics | Planned |
| DIAG-01 | Phase 4 — Management and Diagnostics | Planned |

## Coverage Summary

- Completed baseline, workspace, and usability requirements: 16
- Planned next-step requirements: 11 (6 tool handoff, 5 management/diagnostics)
- Total mapped requirements: 27
- Unmapped requirements: 0
- Phase 2.1 implementation satisfies UX-01; the remaining 11 requirements are planned.
