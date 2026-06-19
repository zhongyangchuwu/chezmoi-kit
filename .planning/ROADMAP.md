# Roadmap: cm v0.1.0 Release Preparation

## Overview

The v0.1.0 release path starts by stabilizing runtime behavior and module state, then cleans public documentation and release metadata, then adds CI/release automation, and finally verifies the whole release gate before tagging.

## Phases

- [x] **Phase 1: Runtime Stability** — Fixed TUI lifecycle, execution semantics, parser strictness, edit tests, and module tidy state.
- [x] **Phase 2: Release Metadata and Documentation Cleanup** — Deleted stale docs, added release metadata, and rewrote public docs to match current behavior.
- [x] **Phase 3: CI and Release Automation** — Added GitHub workflows and GoReleaser config for validation and tag-based release artifacts.
- [ ] **Phase 4: Release Verification** — Run local and CI gates, smoke test release commands, and prepare v0.1.0 tag inputs.

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
3. Release workflow injects the tag version into `internal/build.Version`.
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
2. CI passes on the release preparation branch.
3. Manual smoke checks cover status, diff, sync, edit, git, completion, and version or document any environment-limited gaps.
4. `CHANGELOG.md` and release notes are ready for `v0.1.0`.

**Plans:** 1

**Plans:**

- [ ] `4-release-gate` — Local gate passed; waiting on remote CI and interactive smoke before tag.

## Progress

| Phase | Plans Complete | Status | Completed |
|---|---:|---|---|
| 1. Runtime Stability | 1/1 | Complete | 2026-06-19 |
| 2. Release Metadata and Documentation Cleanup | 1/1 | Complete | 2026-06-19 |
| 3. CI and Release Automation | 1/1 | Complete | 2026-06-19 |
| 4. Release Verification | 0/1 | Blocked on remote CI/manual smoke | - |
