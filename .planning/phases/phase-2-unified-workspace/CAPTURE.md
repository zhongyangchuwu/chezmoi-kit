# Capture: Phase 2 Unified File Workspace

## Durable Docs Updated

- `README.md`
- `docs/usage.md`
- `docs/design.md`
- `docs/development.md`
- `CHANGELOG.md`
- `.planning/codebase/MAP.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/CONCERNS.md`
- `.planning/codebase/STACK.md`
- `.planning/codebase/STRUCTURE.md`

## Planning Records Updated

- `.planning/PROJECT.md`
- `.planning/REQUIREMENTS.md`
- `.planning/ROADMAP.md`
- `.planning/STATE.md`
- Phase `SUMMARY.md`, `REVIEW.md`, `VERIFICATION.md`, and `CAPTURE.md`

## Learnings

- `chezmoi managed --path-style=all --format=json` provides authoritative destination/source mappings without parsing source filenames.
- `chezmoi unmanaged` can return an unmanaged directory rather than its descendant files. Recursive child queries preserve chezmoi classification while exposing nested candidates.
- `--skip-secrets` skips secret-backed templates instead of producing a safe clean value. Omitted template status must become `uninspected` and require explicit diff reveal.
- A shared TUI model is viable when the behavior split is explicit: sync owns confirmation/execution/clean removal, while workspace owns persistent inventory and read-only previews.
- Preview cache keys must include target, view, and reveal state; stale asynchronous completion must not disturb the active view state.

## Ship Inputs

- Branch: `feat/unified-workspace`
- Phase verification: passed.
- Required before merge/release: review the accumulated working tree, create a conventional commit/PR if desired, and observe remote CI.
- Release notes source: `CHANGELOG.md` Unreleased section.
- No tag, push, PR, merge, or release was performed.
