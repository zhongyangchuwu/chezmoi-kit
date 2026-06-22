# State

## Project Reference

See: `.planning/PROJECT.md` (updated 2026-06-22)

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** `v0.2.0` template operation model discussion; no implementation scope committed yet.

## Current Position

- Phase: Phase 1 — Template Operation Model Discussion
- Plan: None
- Status: discussing
- Active Artifact: `.planning/phases/phase-1-template-operation-model-discussion/CONTEXT.md`
- Last activity: 2026-06-22 — archived `v0.1.1`, opened template operation model discussion, and recorded current chezmoi template/edit/merge findings.
- Progress: `v0.1.1` is published and archived; template support decisions are under discussion.

## Accumulated Context

### Decisions

- MIT license is the project license.
- Use GoReleaser for release artifacts.
- GoReleaser signing, notarization, package-manager publishing, and Docker images remain out of scope until a future phase adds them.
- Semantic report documents, not Markdown strings, are the internal output model; Markdown remains an output renderer option.
- Completed release history belongs in `.planning/archive/releases/<version>/`; root planning docs should describe current and future work only.
- Version policy: `v0.1.Z` is for safe additions/fixes that preserve defaults and mental model; `v0.Y.0` is for behavior, safety-model, status-model, output-contract, or command-semantics changes; `v1.0.0` remains deferred.
- Chezmoi template support should delegate rendering, data loading, source naming, editing, and merge mechanics to chezmoi rather than duplicating them in `cm`.

### Blockers/Concerns

- None for current planning state.

## Recent Evidence

- `v0.1.1` tag points to commit `b47b8fa498745bd716de94e1a6ea10339bc06608`.
- GitHub Release `v0.1.1` is published, non-draft, and non-prerelease.
- GitHub Actions CI run `27951265044` passed for `main` at commit `b47b8fa498745bd716de94e1a6ea10339bc06608`.
- GitHub Actions Release run `27951505582` passed for tag `v0.1.1`.
- Release-level evidence is recorded in `.planning/archive/releases/v0.1.1/VERIFICATION.md`.
- Chezmoi docs confirm `source-path <target>` prints a target's source state path, `edit <target>` edits source state and checks template syntax, and `merge <target>` performs destination/source/target three-way merge with two-way fallback when target state cannot be computed.

## Session Continuity

- Last session: 2026-06-22
- Stopped at: template operation model discussion opened; implementation scope intentionally not committed.
- Next Action: decide TUI layout for rendered diff plus source template view, and decide which chezmoi edit/merge/template actions belong in `cm sync`.
- Resume file: None

## Updated

- 2026-06-22
