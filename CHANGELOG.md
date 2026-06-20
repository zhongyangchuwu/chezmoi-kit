# Changelog

All notable changes to `cm` are documented here.

## Unreleased

### Added

- Global `--debug` writes `cm sync` timing diagnostics to a temporary log file and prints the path to stderr.

### Changed

- Package boundaries now use `internal/tui`, `internal/app`, and `internal/diff`; version ldflags target `internal/app.Version`.

### Fixed

- `cm sync` now loads the selected file diff by default.
- `cm sync` confirmed actions now run one target at a time and return to review mode when files remain.
- `cm sync` now prints a distinct completion message when all files are resolved.
- `cm sync` add/apply subprocesses no longer inherit the active TUI terminal, avoiding slow terminal-stdio interactions.

## v0.1.0 - 2026-06-19

### Added

- Initial `cm` CLI for personal chezmoi reconciliation.
- Read-only `cm` and `cm status [target...]` commands.
- Internal `cm diff [target...]` rendering from chezmoi target content to local file content.
- Interactive `cm sync [target...]` two-pane review TUI with pending actions and confirmation.
- Direct mutating wrappers: `cm add`, `cm apply`, `cm merge`, and `cm edit`.
- `cm git` helper for opening `lazygit` in the chezmoi source repository.
- Shell completion generation for bash, zsh, fish, and PowerShell.
- `cm version` build metadata output.

### Fixed

- Bubble Tea signal handling is no longer disabled in production TUI runs.
- Confirmed sync execution no longer advertises or handles `q` as a fake cancellation path.
- Malformed source git status output is reported instead of silently ignored.
- `cm edit` command wiring is covered by tests.
- Go module metadata is tidy for release checks.
- `cm sync` now loads the selected file diff by default.
- `cm sync` confirmed actions now run one target at a time and return to review mode when files remain.
- `cm sync` now prints a distinct completion message when all files are resolved.

### Documentation

- Added MIT license.
- Added release changelog and ignore rules.
- Removed stale historical planning docs from public documentation.
