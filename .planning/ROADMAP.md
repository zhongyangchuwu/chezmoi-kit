# Roadmap: cm v0.1.0 Release Preparation

## Overview

The v0.1.0 release path starts by stabilizing runtime behavior and module state, then cleans public documentation and release metadata, adds CI/release automation, and verifies the release gate. Post-release architecture phases then normalize package boundaries, move application services into `internal/app`, and introduce semantic reports for reusable CLI/TUI output formatting.

## Phases

- [x] **Phase 1: Runtime Stability** — Fixed TUI lifecycle, execution semantics, parser strictness, edit tests, and module tidy state.
- [x] **Phase 2: Release Metadata and Documentation Cleanup** — Deleted stale docs, added release metadata, and rewrote public docs to match current behavior.
- [x] **Phase 3: CI and Release Automation** — Added GitHub workflows and GoReleaser config for validation and tag-based release artifacts.
- [x] **Phase 4: Release Verification** — Completed local release gate, smoke-tested safe commands, documented external terminal/CI/tag gates, and prepared v0.1.0 release inputs.
- [x] **Phase 5: Package Architecture Cleanup** — Normalized package names, removed undersized packages, and split mixed TUI files before behavior moves.
- [x] **Phase 6: App Services and Test Ownership** — Moved application services and sync action ownership into `internal/app`; shrank CLI/TUI tests to their real boundaries.
- [x] **Phase 7: Semantic Reports and Palette Rendering** — Added semantic output documents and renderers for status, diff, version, CLI, and TUI reuse.

## Phase Details

### Phase 1: Runtime Stability

**Goal:** Make existing code safe and tidy enough to become the base for release documentation and CI.

**Depends on:** Nothing (first phase)

**Requirements:** SAFE-01, SAFE-02, SAFE-03, CLI-01, CLI-02, REL-01

**Success Criteria** (what must be TRUE):

1. `cm sync` uses Bubble Tea signal handling in production so normal Ctrl+C cleanup is not disabled.
2. `modeExecuting` no longer advertises or processes a fake quit/cancel path for confirmed actions.
3. Source git status malformed output returns an actionable error instead of being silently skipped.
4. `cm edit <target>` has command wiring test coverage and the fake service implements the edit boundary.
5. `go mod tidy -diff`, `go test ./...`, and targeted package checks pass.

**Plans:** 1

**Plans:**

- [x] `1-runtime-fixes` — TUI lifecycle, parser, command test, and module tidy fixes.

### Phase 2: Release Metadata and Documentation Cleanup

**Goal:** Make repository docs and metadata clean, current, and suitable for public GitHub release.

**Depends on:** Phase 1

**Requirements:** REL-02, DOC-01, DOC-02, DOC-03

**Success Criteria** (what must be TRUE):

1. `docs/superpowers/` is removed from public release documentation.
2. Repository has `LICENSE`, `.gitignore`, `CHANGELOG.md`, and versioned build/install targets.
3. README and `docs/usage.md` document prerequisites, install, all commands including `cm edit`, and current sync semantics.
4. `docs/design.md` and `docs/development.md` reflect current architecture and release workflow.

**Plans:** 1

**Plans:**

- [x] `2-release-docs` — Metadata, documentation rewrite, and stale doc removal.

### Phase 3: CI and Release Automation

**Goal:** Let GitHub enforce release checks and build artifacts from tags.

**Depends on:** Phase 2

**Requirements:** REL-03, REL-04

**Success Criteria** (what must be TRUE):

1. Push and pull request workflow runs Go tests, race tests, vet, tidy diff, module verify, vulnerability scan, and build.
2. Tag workflow for `v*` builds versioned release artifacts for Linux, macOS, and Windows target platforms.
3. Release workflow injects the configured version ldflag into release builds.
4. Workflow files are documented in development/release instructions.

**Plans:** 1

**Plans:**

- [x] `3-ci-release` — GitHub CI and release workflow setup.

### Phase 4: Release Verification

**Goal:** Prove the repository is ready to tag `v0.1.0`.

**Depends on:** Phase 3

**Requirements:** VER-01, VER-02

**Success Criteria** (what must be TRUE):

1. Local release gate passes with observed output.
2. Remote CI status is documented; tag publication remains gated on CI observation and user approval.
3. Manual smoke checks cover status, diff, sync, edit, git, completion, and version or document any environment-limited gaps.
4. `CHANGELOG.md` and release notes are ready for `v0.1.0`.

**Plans:** 1

**Plans:**

- [x] `4-release-gate` — Local gate passed; remote CI/manual smoke/tag publication documented as external follow-ups.

### Phase 5: Package Architecture Cleanup

**Goal:** Make package names and file boundaries express one architecture language before moving more behavior: CLI/TUI adapters, app use cases, and infrastructure capabilities.

**Depends on:** Phase 4

**Requirements:** ARCH-01, TUI-01

**Success Criteria** (what must be TRUE):

1. `internal/ui` is renamed to `internal/tui` and imports/docs match.
2. `internal/build` is merged into `internal/app` as version metadata and release ldflags are updated.
3. `internal/testutil` is removed; tiny test helpers are local or made unnecessary by ownership boundaries.
4. TUI diff state/loading and diff rendering live in separate files.
5. `go test ./...`, `go vet ./...`, and `go mod tidy -diff` pass.

**Plans:** 1

**Plans:**

- [x] `5-package-architecture` — Package rename, build/app merge, testutil removal, diff file split, docs update, and verification complete.

### Phase 6: App Services and Test Ownership

**Goal:** Promote application use cases and sync action ownership into `internal/app`, leaving CLI as command wiring and TUI as terminal interaction.

**Depends on:** Phase 5

**Requirements:** ARCH-02, ARCH-03, TEST-01

**Success Criteria** (what must be TRUE):

1. `internal/app` exposes the service graph and owns status, diff, sync, source git, targets, edit, options, and version metadata.
2. `internal/cli` contains no concrete chezmoi/git/lazygit/diff engine orchestration.
3. `internal/tui` depends on app sync action/service contracts instead of owning them.
4. Terminal merge command construction is tested once in app, not duplicated across CLI/TUI tests.
5. Existing command behavior is preserved under full test/vet/tidy checks.

**Plans:** 1

**Plans:**

- [x] `6-app-services` — App service graph, use cases, sync contracts, test ownership, docs, and verification complete.

### Phase 7: Semantic Reports and Palette Rendering

**Goal:** Add a semantic information model so status, diff, version, and future outputs separate content meaning from plain, ANSI, Markdown, and TUI presentation.

**Depends on:** Phase 6

**Requirements:** OUT-01, OUT-02

**Success Criteria** (what must be TRUE):

1. App-owned status, diff, and version outputs are represented as semantic report documents, not ad hoc strings.
2. Plain and ANSI renderers use a semantic palette for strong, muted, warning, path, command, and diff-line roles.
3. Diff line classification is defined once and reused by report rendering and TUI diff rendering.
4. Color behavior respects TTY detection and non-empty `NO_COLOR`.
5. Markdown output is implemented or explicitly deferred with rationale.

**Plans:** 1

**Plans:**

- [x] `7-semantic-reports` — Semantic report model, renderers, color policy, shared diff classification, docs, and verification complete.

## Progress

| Phase | Plans Complete | Status | Completed |
|---|---:|---|---|
| 1. Runtime Stability | 1/1 | Complete | 2026-06-19 |
| 2. Release Metadata and Documentation Cleanup | 1/1 | Complete | 2026-06-19 |
| 3. CI and Release Automation | 1/1 | Complete | 2026-06-19 |
| 4. Release Verification | 1/1 | Complete | 2026-06-20 |
| 5. Package Architecture Cleanup | 1/1 | Complete | 2026-06-20 |
| 6. App Services and Test Ownership | 1/1 | Complete | 2026-06-20 |
| 7. Semantic Reports and Palette Rendering | 1/1 | Complete | 2026-06-20 |
