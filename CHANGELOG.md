# Changelog

All notable changes to `cm` are documented here.

## Unreleased

### Added

- `cm ui [path...]` provides a persistent, read-only workspace for managed clean/dirty files, source-ignored entries, and explicitly scoped unmanaged candidates.
- Workspace category filters, global path/preview search, full-screen preview, horizontal scrolling, diff hunk navigation, and labeled destination/target/source views.
- Explicit per-view reveal for uninspected diffs, rendered template/encrypted targets, and encrypted source content; no-scope workspace inventory never scans unmanaged `$HOME` paths.
- `cm ui` now uses semantic state/type/attribute colors, a compact legend, responsive key hints, and `?` in-TUI help while retaining text-only fallback under `NO_COLOR=1`.
- `cm ui` now uses a Yazi-inspired Parent/Current/Preview directory browser with virtual path ancestors, direct-child directory navigation, local directory summaries, and aligned trailing-`/` directory rows.

### Changed

- `cm diff` and `cm sync` now force chezmoi's builtin diff with reverse direction and no pager, preserving target-to-destination review semantics while correctly representing permissions, symlinks, directories, templates, and other target types.
- `cm status` reports pending chezmoi scripts in a separate `automation:` block; targetless `cm diff` and `cm sync` no longer treat scripts as ordinary file reconciliation.
- Reconciliation entries now follow chezmoi's second status column, so first-column-only history drift is not shown as destination/target work.
- Sync actions are gated by reviewed target type and template state, and pending actions carry a fingerprint of the exact review.
- Sync execution rechecks the review before mutation and verifies target status afterward; stale or unresolved targets return to review instead of producing false completion.
- Successful buffered chezmoi output and warnings remain visible in sync results.
- CI and release workflows install pinned chezmoi v2.72.1 and run isolated real-tool coverage for content, metadata, symlink, directory, remove, script, template, encrypted re-add, and custom destination behavior.

### Fixed

- `cm edit` now resolves relative targets against chezmoi's configured destination directory and preserves absolute targets instead of assuming `$HOME`.
- Confirmed sync add continues to use `chezmoi re-add`, preserving encrypted source attributes.

## v0.1.1 - 2026-06-22

### Added

- Optional report controls for `cm`, `cm status`, `cm diff`, `cm doctor`, and `cm version`: `--output plain|ansi|markdown` and `--color auto|always|never`.
- Read-only `cm doctor` diagnostics for required `chezmoi`/`git` prerequisites, chezmoi source path readability, source repository status, and optional `lazygit` availability.

### Fixed

- `cm sync` TUI truncation now uses display width instead of byte length for wide-character file paths and diff lines.


## v0.1.0 - 2026-06-20

### Added

- Initial `cm` CLI for personal chezmoi reconciliation.
- Read-only `cm` and `cm status [target...]` commands.
- Internal `cm diff [target...]` rendering from chezmoi target content to local file content.
- Interactive `cm sync [target...]` two-pane review TUI with pending actions and confirmation.
- Direct mutating wrappers: `cm add`, `cm apply`, `cm merge`, and `cm edit`.
- `cm git` helper for opening `lazygit` in the chezmoi source repository.
- Shell completion generation for bash, zsh, fish, and PowerShell.
- `cm version` build metadata output.
- Global `--debug` writes `cm sync` timing diagnostics to a temporary log file.
- Semantic report documents and palette renderers for status, diff, version, and diff line classification (plain, ANSI with TTY/NO_COLOR detection, Markdown).

### Changed

- Package architecture normalized: `internal/tui` (terminal UI), `internal/app` (service graph, use cases, sync contracts), `internal/report` (semantic output), `internal/diff`, `internal/chezmoi`, `internal/process`.
- Version ldflags target `internal/app.Version` instead of `internal/build`.
- `internal/cli` reduced to Cobra wiring, flag parsing, and stream plumbing only.
- Release notes now reference curated `CHANGELOG.md` instead of auto-generated commit hashes.

### Fixed

- Bubble Tea signal handling is no longer disabled in production TUI runs.
- Confirmed sync execution no longer advertises or handles `q` as a fake cancellation path.
- Malformed source git status output is reported instead of silently ignored.
- `cm edit` command wiring is covered by tests.
- Go module metadata is tidy for release checks.
- `cm sync` now loads the selected file diff by default.
- `cm sync` confirmed actions now run one target at a time and return to review mode when files remain.
- `cm sync` now prints a distinct completion message when all files are resolved.
- `cm sync` add/apply subprocesses no longer inherit the active TUI terminal.

### Documentation

- Added MIT license.
- Added release changelog and ignore rules.
- Removed stale historical planning docs from public documentation.
- Added `docs/design.md` and `docs/development.md` covering architecture and workflow.
