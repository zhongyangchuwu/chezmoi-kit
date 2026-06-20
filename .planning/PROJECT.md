# cm

## What This Is

`cm` is a small Go CLI for safer personal reconciliation of chezmoi-managed configuration files. It keeps chezmoi as the authority while adding a simpler status model, internal diffs, direct wrappers, and an interactive review TUI for sync decisions.

The current project goal is to prepare a clean `v0.1.0` GitHub release with stable runtime behavior, accurate documentation, CI, and reproducible release steps.

## Core Value

Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Requirements

### Validated

- ✓ Read-only default command renders local chezmoi mismatch and source repository git status — existing.
- ✓ `cm status [target...]` uses the same read-only status path as the root command — existing.
- ✓ `cm diff [target...]` renders internal target-vs-local diffs without shelling out to external diff — existing.
- ✓ `cm sync [target...]` provides a two-pane Bubble Tea review workflow with pending actions and confirmation — existing.
- ✓ Direct wrappers exist for `cm add`, `cm apply`, `cm merge`, `cm edit`, and `cm git` — existing.
- ✓ Shell completion generation exists for bash, zsh, fish, and PowerShell — existing.
- ✓ Build/version information is available through `cm version` — existing.
- ✓ TUI signal handling, executing-mode semantics, source git parser strictness, `cm edit` command wiring tests, and module tidy state are release-ready — Phase 1 complete.
- ✓ Public release metadata and documentation are clean for v0.1.0 — Phase 2 complete.
- ✓ GitHub CI and GoReleaser release automation are configured and locally verified — Phase 3 complete.
- ✓ Release readiness is locally verified and remaining external terminal/CI/tag gates are documented — Phase 4 complete.
- ✓ Package boundaries now use CLI/TUI adapters, app support, and infrastructure capability packages — Phase 5 complete.
- ✓ Application service ownership now lives in `internal/app`, with CLI/TUI tests aligned to package ownership — Phase 6 complete.
- ✓ Semantic report documents and palette renderers now back status, diff, version, and shared diff classification — Phase 7 complete.

### Active

- [x] Make TUI signal handling and execution-state behavior safe for `v0.1.0`.
- [x] Keep Go module metadata tidy and reproducible.
- [x] Cover exposed command behavior with focused tests.
- [x] Remove stale planning docs from public documentation.
- [x] Provide release metadata: license, changelog, ignore rules, versioned build commands.
- [x] Document current installation, requirements, usage, architecture, and release workflow.
- [x] Add GitHub CI and release workflows for tag-based `v0.1.0` publishing.
- [x] Close release readiness with local gates and documented external ship gates.
- [x] Normalize package boundaries into CLI/TUI adapters, app services, and infrastructure capabilities.
- [x] Move application service ownership into `internal/app` and make CLI/TUI tests match real package ownership.
- [x] Introduce semantic report documents and palette renderers for reusable status, diff, version, and TUI output.

### Out of Scope

- Replacing chezmoi — `cm` remains a wrapper and review layer.
- Automatic git commit, push, pull, or source repository automation.
- Daemon/watch mode or background synchronization.
- Persistent state database beyond chezmoi's own state.
- Full cancellation of already-started mutating subprocesses for `v0.1.0`; executing mode will use conservative non-cancellable semantics unless a later phase explicitly changes the runner contract.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package manager publishing.

## Context

- Brownfield codebase map lives in `.planning/codebase/MAP.md`.
- Historical `docs/superpowers/` files were deleted from public docs in Phase 2.
- Release verification is complete locally; `v0.1.0` publication remains gated on observed remote CI, real-terminal smoke where available, and explicit tag approval.

## Constraints

- Preserve a clean cutover: no compatibility shims for stale pre-release behavior.
- Prefer small, direct Go changes over broad abstractions.
- Keep tests behavior-focused at command/service/model boundaries.
- Do not expand `v0.1.0` scope into large new features; `cm doctor` is deferred unless explicitly pulled into release scope.
- Release docs must match observed code, not historical plans.
- User-facing non-interactive output should be modeled semantically in app and rendered by explicit plain/ANSI/Markdown/TUI renderers.
- CI must run the same checks expected before tagging.

## Key Decisions

| Date | Decision | Rationale | Outcome |
|---|---|---|---|
| 2026-06-19 | Prepare `v0.1.0` through phased workflow artifacts. | Release work spans runtime, docs, CI, and process; durable state is needed. | Initialized `.planning/`. |
| 2026-06-19 | Delete `docs/superpowers/` before release. | Historical artifacts contradict current behavior and mislead contributors. | Scheduled in Phase 2. |
| 2026-06-19 | Fix runtime behavior before release metadata and CI. | CI should validate a stable target, not freeze known release risks. | Phase 1 focuses on runtime stability. |
| 2026-06-19 | Use conservative executing-mode semantics for `v0.1.0`. | Full subprocess cancellation requires broader runner/context changes. | Executing mode should not promise quit/cancel. |
| 2026-06-19 | Use GoReleaser for v0.1.0 release artifacts. | User explicitly preferred GoReleaser over hand-written artifact upload. | Added `.goreleaser.yaml` and GoReleaser workflows. |
| 2026-06-19 | Completed Phase 1 runtime stability before documentation/CI work. | Automated checks passed and release-blocking runtime issues were fixed. | Phase 2 can start on metadata and docs. |
| 2026-06-19 | Use MIT license for v0.1.0. | User approved MIT; it is simple and suitable for a small public CLI. | Added `LICENSE`. |
| 2026-06-19 | Completed Phase 2 documentation and release metadata cleanup. | Public docs now match current code and release helpers inject explicit versions. | Phase 3 can start on CI/release workflows. |
| 2026-06-19 | Completed Phase 3 CI and release automation. | Local GoReleaser check and snapshot builds passed. | Phase 4 can start final release verification. |
| 2026-06-20 | Completed Phase 4 release verification closure. | Local gates, snapshot artifacts, and safe CLI smoke passed; external terminal/CI/tag gates are documented as ship inputs. | Phase 5 package architecture cleanup can proceed. |
| 2026-06-20 | Plan architecture cleanup as Phases 5-7 after release verification. | Package names and tests currently mix adapters, app behavior, and infrastructure. | Added package architecture, app services, and semantic report phases. |
| 2026-06-20 | Use semantic report documents rather than Markdown as the internal output model. | Markdown is useful output, but it loses app semantics needed by ANSI, TUI, and future formats. | Phase 7 will add semantic reports plus plain/ANSI/Markdown renderers. |
| 2026-06-20 | Completed Phase 5 package architecture cleanup. | Package names and file boundaries now match the planned architecture language before moving app services. | Phase 6 app services and test ownership can proceed. |
| 2026-06-20 | Completed Phase 6 app service ownership. | CLI/TUI packages no longer own concrete app orchestration or sync action domain contracts. | Phase 7 semantic reports and palette rendering can proceed. |
| 2026-06-20 | Completed Phase 7 semantic reports and palette rendering. | Status, diff, and version output now preserve semantic roles before rendering. | Planned post-release architecture phases are complete; remaining release blockers are external gates. |

_Last updated: 2026-06-20_
