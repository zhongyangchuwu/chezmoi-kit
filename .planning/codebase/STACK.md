# Stack

**Mapped:** 2026-09-07

## Language and Toolchain

- Go module: `github.com/zhongyangchuwu/cm`
- Go language version: `1.26`
- Entrypoint: `cmd/cm/main.go`
- Formatting: `gofmt`
- Tests: standard `testing`
- Static analysis: `go vet`
- Vulnerability scan: `govulncheck`
- Release: GoReleaser v2

## Direct Runtime Dependencies

| Module | Purpose |
|---|---|
| `charm.land/bubbletea/v2` | TUI runtime and terminal command handoff. |
| `charm.land/bubbles/v2` | Help and key-binding models. |
| `charm.land/lipgloss/v2` | Pane layout, styles, display-width handling. |
| `github.com/charmbracelet/x/ansi` | Display-width-safe ANSI truncation. |
| `github.com/spf13/cobra` | CLI command tree and completion generation. |
| `golang.org/x/term` | Terminal detection. |

The former `github.com/rogpeppe/go-internal/diff` dependency is removed because chezmoi now produces authoritative builtin diff output.

## External Runtime Commands

- `chezmoi` — required for managed status, diff, metadata, edit, and mutation.
- `git` — required for source repository status.
- `lazygit` — optional except for `cm git`.

## Chezmoi Commands Used

- `status --include=all --exclude=none --path-style=absolute` and `--skip-secrets` inventory status
- forced builtin `diff --include=all --exclude=none --reverse`
- `managed --include=all --exclude=none --path-style=all --format=json`
- typed `managed` membership plus NUL-delimited `ignored` and scoped `unmanaged`
- `cat`, `decrypt`, `target-path`, `source-path`
- `dump --include=all --exclude=none --format=json --recursive=false`
- `re-add`, `apply --force`, `merge`, `add`, `edit`

## Output Formats

- plain text
- ANSI text with TTY/`NO_COLOR` handling
- Markdown
- Bubble Tea TUI

JSON output remains deferred because report presentation blocks are not the correct domain schema.

## Build and Release

- `just install` installs the binary and zsh completion.
- `just build-release` produces a local version-injected binary.
- GoReleaser builds Linux, macOS, and Windows archives for amd64/arm64.
- CI and release workflows download the `chezmoi-linux-amd64` v2.72.1 release asset for real-tool tests.
- CI runs tidy diff, module verification, tests, race tests, vet, govulncheck, GoReleaser check, and snapshot release.

## Repository Shape

- No Viper.
- No generated application code.
- No custom diff engine.
- No external assertion/mocking library.
- No runtime database owned by cm.
