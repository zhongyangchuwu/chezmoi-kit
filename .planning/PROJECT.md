# cm

## What This Is

`cm` is a Go CLI/TUI whose product goal is a complete chezmoi management entrypoint. It owns file discovery, browsing, review, confirmation, and state refresh while reusing chezmoi for state/mutation, lazygit for Git, and configured editors/merge tools for specialized work.

The implemented baseline includes authoritative reconciliation and the verified Phase 2 unified file workspace. The remaining steps integrate mature-tool round trips, then add contextual management/diagnostics; they are planned and have no assigned release numbers.

## Core Value

Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Requirements

### Validated

- ✓ `cm` and `cm status [target...]` render local chezmoi mismatch and source repository git status.
- ✓ `cm diff [target...]` renders target-vs-local diffs without invoking a user-configured external diff tool.
- ✓ `cm sync [target...]` provides a two-pane Bubble Tea review workflow with pending actions and confirmation.
- ✓ Confirmed sync add uses `chezmoi re-add`, preserving encrypted source attributes.
- ✓ Direct wrappers exist for `cm add`, `cm apply`, `cm merge`, `cm edit`, and `cm git`.
- ✓ Shell completion generation exists for bash, zsh, fish, and PowerShell.
- ✓ Build/version information, semantic report rendering, doctor diagnostics, and display-width-safe TUI truncation are available.
- ✓ `v0.1.0` and `v0.1.1` were published with passing CI and GoReleaser workflows.
- ✓ `v0.2.0` authoritative reconciliation implementation is complete and verified on `feat/authoritative-reconciliation`: forced chezmoi builtin diff, second-column status routing, typed reviews/action gating, script separation, reviewed fingerprints, postflight verification, visible action output, custom destination edit, and real-tool integration coverage.
- ✓ `cm ui [path...]` is a persistent, read-only workspace for managed clean/dirty, source-ignored, and explicitly scoped unmanaged entries; its bounded previews require explicit reveal for sensitive content.

### Active

- Phase 2 unified workspace is complete and verified; it preserves bare `cm` status and focused `cm sync` contracts.
- Next step 2 / Phase 3: Mature Tool Handoffs — configured editor, chezmoi merge, and lazygit with terminal restoration and refreshed state.
- Next step 3 / Phase 4: Management and Diagnostics — contextual operations, target lifecycle/attributes, and actionable diagnostics.
- `.planning/ROADMAP.md` defines remaining scope, exclusions, dependencies, and acceptance scenarios. `.planning/REQUIREMENTS.md` retains template requirements across the remaining steps.

### Out of Scope

- Replacing chezmoi — `cm` remains a review and orchestration layer.
- Automatic git commit, push, pull, or source repository automation. Explicit user-driven Git operations belong in lazygit launched from cm.
- Daemon/watch mode or background synchronization.
- Persistent state beyond chezmoi's own state.
- Full cancellation of already-started mutating subprocesses until runner/client/app contracts explicitly support it.
- GoReleaser signing, notarization, package-manager publishing, and Docker images.
- Reimplementing chezmoi template rendering, source naming, data loading, ignore rules, encryption, diff, or merge semantics.
- Reimplementing a Git client, text editor, or merge tool solely to keep all UI code inside cm.
- Embedded PTY subpanes, destroy/purge UI, script-execution workflows, generic config transformations, backup/undo storage, and application validators in the next three steps.

## Context

- Brownfield codebase maps live under `.planning/codebase/`.
- Completed release history lives under `.planning/archive/releases/`.
- Pre-phase isolated checks against chezmoi v2.72.1 reproduced metadata-only `no diff`, symlink dereference, directory failure, script/file semantic mismatch, and custom destination edit failure in the former implementation.
- The verified implementation now delegates forced builtin diff semantics, uses bounded target metadata and template filters, and protects execution with reviewed fingerprints and postflight status.
- CI and release workflows download the official pinned chezmoi v2.72.1 Linux asset for isolated integration coverage.
- Product discussion selected integration over novelty: proven frontend features are valuable when they make cm the daily entrypoint, but mature Git/editor/merge workflows should be reused rather than recreated.
- Existing `OpenSourceGit`, `EditTarget`, and terminal merge delegation provide starting points. The missing work is full inventory, in-TUI entrypoints, context-preserving round trips, and consistent refresh—not new Git or editor engines.

## Constraints

- Preserve a clean cutover; remove obsolete custom diff paths rather than retaining parallel implementations.
- Prefer small direct Go types and methods over a generic reconciliation framework.
- Keep tests behavior-focused at CLI, app service, TUI state-machine, parser, and real-tool integration boundaries.
- Never execute scripts through the ordinary file reconciliation action model.
- Never log rendered target contents, template contents, or action subprocess output to the timing log.
- User-facing non-interactive output remains semantic app data rendered through explicit plain/ANSI/Markdown adapters.
- CI must run the same checks expected before tagging.
- A unified entrypoint does not require a single implementation: use suspend/run/restore for mature terminal tools and restore the workspace on normal exit, cancellation, and failure.
- External tools may change source files and membership; refresh managed/local/Git state and revalidate pending reviews after returning.
- Preserve bare `cm` as read-only status and focused `cm sync` behavior; add the general workbench explicitly instead of changing defaults silently.
- Browse/filter/preview features do not authorize mutation. New lifecycle and attribute actions require explicit scope and confirmation.

## Key Decisions

| Date | Decision | Rationale | Outcome |
|---|---|---|---|
| 2026-06-19 | Use MIT license and GoReleaser release artifacts. | Small public CLI with conventional release automation. | License and release workflows shipped. |
| 2026-06-20 | Use semantic report documents rather than Markdown strings internally. | ANSI, Markdown, TUI, and future formats need retained semantics. | `internal/report` owns semantic documents and renderers. |
| 2026-06-22 | Reserve minor releases for safety-model, status-model, output-contract, and command-semantics changes. | Users should consciously absorb low-surprise contract changes. | Authoritative reconciliation is scoped as `v0.2.0`. |
| 2026-06-22 | Delegate template evaluation, editing, and merge mechanics to chezmoi. | Duplicating source naming, data, encryption, and template behavior would be unsafe. | Template research favors chezmoi-backed operations. |
| 2026-09-05 | Prioritize authoritative target semantics before template-specific UI. | Regular-file assumptions misrepresented metadata, symlinks, directories, and scripts. | `v0.2.0` implementation and verification completed; template workspace remains future scope. |
| 2026-09-05 | Bind confirmed actions to a reviewed fingerprint and verify postconditions. | A dirty-only preflight can execute against unreviewed content, and command success does not prove reconciliation. | Stale actions defer without mutation; only verified-clean targets leave the TUI. |
| 2026-09-06 | Make cm the complete management entrypoint by integrating proven workflows. | Users benefit from one place to find, inspect, edit, reconcile, and manage source history; originality is not a prerequisite for value. | Three planned steps cover workspace, tool handoffs, and management/diagnostics. |
| 2026-09-06 | Delegate Git to lazygit and editing/merging to configured tools. | Rebuilding mature specialized interfaces adds maintenance without improving the daily workflow. | cm owns launch context, terminal lifecycle, state restoration, and review invalidation. |

_Last updated: 2026-09-07_
