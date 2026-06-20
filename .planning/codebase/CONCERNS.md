# Codebase Concerns

## Release and workflow risks

- `v0.1.0` publication still requires observed remote GitHub CI and explicit tag approval.
- Interactive smoke for `cm sync` Ctrl+C cleanup and `cm git` lazygit startup still needs a real terminal when preparing the tag.
- `cm git` can mutate the chezmoi source repository through lazygit even though `cm` itself only launches the tool; docs should keep that distinction clear.

## Untested paths and behavior gaps

- Shell-specific completion behavior is only lightly smoke-tested; Cobra generation is wired, but shells differ in runtime integration.
- `cmd/cm/main.go` only forwards `cli.Main` to `os.Exit`; command-level behavior is tested through `internal/cli`, but the process entrypoint itself has no visible smoke coverage.
- `internal/tui/view.go` truncates by byte length, not display width or rune boundaries; Unicode paths or diff lines may render poorly.
- `internal/chezmoi/content.go` reads local target files with a bounded buffer; direct local file read edge cases are lower-level than current command tests.

## TUI lifecycle risks

- `internal/tui/tui.go` changes behavior based on terminal detection; command tests mostly exercise non-terminal output.
- Tiny terminal sizes can severely truncate the two-pane layout.
- Diffs are cached per target for the session; browsing many changed files retains loaded diff lines in memory.
- Confirmed add/apply commands run as non-interactive subprocesses, while merge stays terminal-bound; real merge tools can block or fail in environment-specific ways.

## Subprocess and external-tool risks

- Missing `chezmoi`, `git`, or `lazygit` all pass through the process runner; user-facing error quality depends on the path that adds context.
- Source git status errors include source-directory context and trimmed stderr; useful for debugging but may expose local filesystem paths in shared logs.
- `cm diff` can display managed configuration contents in terminal scrollback, including secrets present in dotfiles.

## Performance and memory risks

- `cm diff` builds a full target list before diffing all dirty files; large status sets are not streamed.
- `internal/diff` caps each side at 1 MiB plus sentinel, bounding common diff memory while still allowing large output for many targets.
- TUI view rendering allocates visible slices and rendered strings per render; acceptable for small terminal views but worth revisiting if large diffs feel sluggish.

## Maintainability gaps

- Report rendering now uses a small semantic model; future output expansion should avoid adding generic table/layout abstractions until a command needs them.
- TUI file/diff layout still truncates by byte length, not display width or rune boundaries; Unicode paths or diff lines may render poorly.
