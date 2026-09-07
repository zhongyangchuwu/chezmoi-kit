# Summary: Phase 2 Unified File Workspace

## Completed Changes

- Added persistent, read-only `cm ui [path...]` without changing bare `cm` status or focused `cm sync` behavior.
- Added chezmoi-backed managed path mappings, typed membership, source-ignored inventory, secret-skipping status, scoped unmanaged discovery, and bounded target/source content adapters.
- Added `WorkspaceSnapshot`, `WorkspaceEntry`, `WorkspacePreview`, and `WorkspaceService` at the app boundary.
- Added merged managed clean/dirty/script, source-ignored, and explicitly scoped unmanaged inventory with recursive candidate discovery caps.
- Generalized the Bubble Tea model into explicit sync and workspace modes. Workspace keeps clean rows and has no mutation path.
- Added tree/flat projections, collapse, filters, path search, stable absolute-path selection restoration, labeled previews, explicit reveal, full preview, horizontal/vertical scrolling, preview search, and diff hunk navigation.
- Added real chezmoi integration coverage for clean/dirty/template/encrypted/symlink/directory/script/ignored/nested-unmanaged inventory and explicit secret-template reveal.
- Added public workspace documentation, changelog notes, and refreshed codebase maps.

## Files Changed

- Runtime: `internal/app/workspace.go`, `internal/app/workspace_service.go`, `internal/chezmoi/workspace.go`, `internal/tui/workspace_list.go`, and shared CLI/TUI/service files.
- Tests: workspace adapter, app real-tool integration, TUI interaction, and CLI command coverage.
- User docs: `README.md`, `docs/usage.md`, `docs/design.md`, `docs/development.md`, and `CHANGELOG.md`.
- Planning: `.planning/codebase/*` and this phase directory.

## Deviations

- Native `chezmoi unmanaged` may return a directory rather than nested files. The workspace recursively re-queries immediate children through chezmoi rather than treating filesystem traversal as membership classification.
- Secret-backed templates omitted by `--skip-secrets` are represented as `uninspected`, never as clean. An explicit diff reveal obtains the authoritative clean/dirty state.
- A stale preview completion originally reset active preview state. The TUI now caches the stale result without changing the active view's message, match set, or scroll position.

## Evidence

- Focused `internal/app`, `internal/chezmoi`, `internal/tui`, and `internal/cli` tests passed during implementation.
- Isolated real chezmoi workspace integration passed, including explicit secret-provider nonexecution before reveal and invocation only after reveal.
- Actual PTY smoke exercised clean persistence, scoped unmanaged and ignored labels, tree/filter/search, long-line horizontal scrolling, diff hunk jumps, full preview, template target withholding/reveal, and no-mutation exit.

## Unresolved Risks

- Explicit reveal intentionally exposes rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is capped at 10,000 entries and 2,000 batches; users must narrow unusually large scopes.
- Remote CI, pull request, merge, and release publication remain separate user-directed ship actions.
