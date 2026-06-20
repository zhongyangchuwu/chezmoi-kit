# Changelog

All notable changes to `cm` are documented here.

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
