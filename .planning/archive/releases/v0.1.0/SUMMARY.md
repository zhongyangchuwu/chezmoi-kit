# v0.1.0 Summary

## Released
2026-06-20

## What Shipped
`cm` v0.1.0 shipped as the initial public GitHub release for safer personal chezmoi reconciliation. The release includes the CLI command surface, interactive sync review TUI, release metadata, CI and GoReleaser automation, app-oriented package boundaries, and semantic report rendering for status, diff, and version output.

## Included Phases
- `phase-1-runtime-stability` — Runtime Stability.
- `phase-2-release-docs` — Release Metadata and Documentation Cleanup.
- `phase-3-ci-release` — CI and Release Automation.
- `phase-4-release-gate` — Release Verification.
- `phase-5-package-architecture` — Package Architecture Cleanup.
- `phase-6-app-services` — App Services and Test Ownership.
- `phase-7-semantic-reports` — Semantic Reports and Palette Rendering.

## Completed Scope

### Roadmap Items
- Phase 1: Fixed TUI lifecycle, execution semantics, parser strictness, edit tests, and module tidy state.
- Phase 2: Removed stale public docs, added MIT license, `.gitignore`, changelog, versioned build/install helpers, and rewrote README/design/usage/development docs.
- Phase 3: Added GitHub CI workflow, tag-based release workflow, and GoReleaser v2 configuration for Linux, macOS, and Windows artifacts.
- Phase 4: Completed local release gate, GoReleaser snapshot validation, safe binary smoke checks, changelog dating, and release-target correction.
- Phase 5: Normalized architecture names by moving `internal/ui` to `internal/tui`, merging build metadata into `internal/app`, removing `internal/testutil`, and splitting TUI diff state/view code.
- Phase 6: Moved application service graph, use cases, sync action contracts, and terminal command ownership into `internal/app`; reduced CLI/TUI to adapter boundaries.
- Phase 7: Added semantic report documents, plain/ANSI/Markdown renderers, color policy, shared diff classification, and CLI/TUI report rendering integration.

### Requirements
- SAFE-01 — `cm sync` preserves terminal state by delegating signal handling to Bubble Tea.
- SAFE-02 — `cm sync` no longer advertises fake cancellation once confirmed execution starts.
- SAFE-03 — malformed source repository git status output surfaces an actionable error.
- CLI-01 — exposed v0.1.0 command behavior has focused command/service coverage.
- CLI-02 — `cm edit <target>` remains wired to `chezmoi edit`; completion uses managed file names.
- REL-01 — Go module metadata is tidy and reproducible under release checks.
- REL-02 — repository release metadata exists: MIT license, changelog, ignore rules, and versioned build commands.
- REL-03 — GitHub CI enforces the local release gate on pushes and pull requests.
- REL-04 — tag-based GitHub release workflow builds versioned release artifacts.
- DOC-01 — public documentation describes current requirements, installation, usage, commands, and release workflow.
- DOC-02 — stale `docs/superpowers/` public documentation was removed.
- DOC-03 — architecture documentation matches implemented package boundaries and behavior.
- VER-01 — release readiness was proven by observed local checks before tagging.
- VER-02 — smoke scenarios for status, diff, sync, edit, git, completion, and version were completed or documented with environment-limited gaps.
- ARCH-01 — package boundaries use CLI/TUI adapters, app use cases/services, and infrastructure capabilities.
- ARCH-02 — `internal/app` owns the service graph, sync action contract, command use cases, and version metadata.
- ARCH-03 — `internal/cli` owns Cobra wiring, flag parsing, stream plumbing, and exit behavior only.
- TUI-01 — terminal UI package is `internal/tui`, with diff state/loading split from diff rendering.
- OUT-01 — status, diff, and version outputs use semantic report documents.
- OUT-02 — report rendering uses semantic palette roles with TTY detection and `NO_COLOR` handling.
- TEST-01 — tests align with package ownership; app-owned terminal command construction is tested once in app.

## Notable Decisions
- Use phased `.planning/` workflow for v0.1.0 because runtime, docs, CI, verification, and architecture work needed durable state; outcome: seven verified phases archived here.
- Delete `docs/superpowers/` because historical planning artifacts contradicted current public behavior; outcome: stale public docs removed in Phase 2.
- Use conservative executing-mode semantics because broad subprocess cancellation was outside v0.1.0; outcome: confirmed execution no longer promises cancel/quit.
- Use MIT license for the first public release; outcome: `LICENSE` added in Phase 2.
- Use GoReleaser for release artifacts; outcome: `.goreleaser.yaml`, release workflow, and published `v0.1.0` GitHub release.
- Set GoReleaser release target to `zhongyangchuwu/chezmoi-kit`; outcome: Phase 4 corrected the repository target before publication.
- Treat semantic report documents, not Markdown strings, as the internal output model; outcome: Phase 7 added report documents and renderers.

## Follow-ups
- `cm doctor` prerequisite checker.
- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and app sync services.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.
- Issue templates, PR templates, and security policy files.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.
- User-facing `--output` and `--color` flags for existing internal report renderers.
- TUI file/diff layout display-width handling.
- Real-terminal smoke for `cm sync` Ctrl+C cleanup and `cm git` lazygit startup remains environment-dependent evidence outside this harness.
