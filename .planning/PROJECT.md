# cm

## What This Is

`cm` is a small Go CLI for safer personal reconciliation of chezmoi-managed configuration files. It keeps chezmoi as the authority while adding a simpler status model, internal diffs, direct wrappers, and an interactive review TUI for sync decisions.

The current project state is post-`v0.1.0`: the initial GitHub release is published and the completed release planning history is archived under `.planning/archive/releases/v0.1.0/`.

## Core Value

Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Requirements

### Validated

- ✓ `cm` and `cm status [target...]` render local chezmoi mismatch and source repository git status.
- ✓ `cm diff [target...]` renders internal target-vs-local diffs without shelling out to external diff.
- ✓ `cm sync [target...]` provides a two-pane Bubble Tea review workflow with pending actions and confirmation.
- ✓ Direct wrappers exist for `cm add`, `cm apply`, `cm merge`, `cm edit`, and `cm git`.
- ✓ Shell completion generation exists for bash, zsh, fish, and PowerShell.
- ✓ Build/version information is available through `cm version`.
- ✓ TUI signal handling, executing-mode semantics, source git parser strictness, `cm edit` command wiring tests, and module tidy state are release-ready.
- ✓ Public release metadata and documentation are clean for v0.1.0.
- ✓ GitHub CI and GoReleaser release automation are configured and have passed for `v0.1.0`.
- ✓ Package boundaries use CLI/TUI adapters, app services, infrastructure capability packages, and semantic report rendering.
- ✓ Semantic report documents and palette renderers back status, diff, version, and shared diff classification.
- ✓ `v0.1.0` is published as a non-draft, non-prerelease GitHub Release.
- ✓ `cm`, `cm status`, `cm diff`, and `cm version` support optional `--output plain|ansi|markdown` and `--color auto|always|never` report controls.
- ✓ `cm doctor` reports read-only prerequisite diagnostics for required and optional external tools.
- ✓ `cm sync` TUI truncation uses display width for file names and diff lines, preserving valid UTF-8 for wide characters.

### Active

- `v0.1.1` is planned as a safe usability-polish release.
- Current scope must not change default command usage, the README mental model, existing command semantics, or mutation behavior.
- Remaining planned requirement: None; `v0.1.1` planned implementation is complete.

### Out of Scope

- Replacing chezmoi — `cm` remains a wrapper and review layer.
- Automatic git commit, push, pull, or source repository automation.
- Daemon/watch mode or background synchronization.
- Persistent state database beyond chezmoi's own state.
- Full cancellation of already-started mutating subprocesses until the runner/client/app contracts explicitly support it.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.

## Context

- Brownfield codebase map lives in `.planning/codebase/MAP.md`.
- Completed `v0.1.0` planning history lives in `.planning/archive/releases/v0.1.0/`.
- GitHub Release `v0.1.0` was published on 2026-06-20.
- GitHub Actions CI and Release workflows passed for commit `f8e80afb408b93a53e743a58be52d31f9cf842f4`.
- Real-terminal smoke for `cm sync` Ctrl+C cleanup and `cm git` lazygit startup remains environment-dependent evidence not available in the non-interactive harness.
- `v0.1.1` planning should keep README's current mental model intact; README changes belong after implementation verifies new optional behavior.

## Constraints

- Preserve a clean cutover: no compatibility shims for stale pre-release behavior.
- Prefer small, direct Go changes over broad abstractions.
- Keep tests behavior-focused at command/service/model boundaries.
- Do not expand scope into large new features without a new phase.
- Release docs must match observed code, not historical plans.
- User-facing non-interactive output should be modeled semantically in app and rendered by explicit plain/ANSI/Markdown/TUI renderers.
- CI must run the same checks expected before tagging.

## Key Decisions

| Date | Decision | Rationale | Outcome |
|---|---|---|---|
| 2026-06-19 | Use MIT license for v0.1.0. | User approved MIT; it is simple and suitable for a small public CLI. | Added `LICENSE`. |
| 2026-06-19 | Use GoReleaser for v0.1.0 release artifacts. | User explicitly preferred GoReleaser over hand-written artifact upload. | Added `.goreleaser.yaml` and GoReleaser workflows. |
| 2026-06-20 | Use semantic report documents rather than Markdown as the internal output model. | Markdown is useful output, but it loses app semantics needed by ANSI, TUI, and future formats. | Added semantic reports plus plain/ANSI/Markdown renderers. |
| 2026-06-20 | Publish `v0.1.0` from repository `zhongyangchuwu/chezmoi-kit`. | GoReleaser must target the actual remote repository. | GitHub Release `v0.1.0` published successfully. |
| 2026-06-22 | Use `v0.1.Z` for safe additions/fixes and reserve `v0.Y.0` for behavior or mental-model changes before `v1.0.0`. | `cm` values low-surprise behavior; version numbers should signal whether users need to relearn defaults or safety semantics. | Planned `v0.1.1` as safe usability polish: output controls, doctor diagnostics, and TUI width fixes. |
| 2026-06-21 | Archive completed `v0.1.0` planning history. | Root planning docs should describe only current and future work after release. | Release archive created under `.planning/archive/releases/v0.1.0/`. |

_Last updated: 2026-06-22_
