# Stack Map

## Language and runtime

- The project is a Go CLI module declared as `module github.com/zhongyangchuwu/cm` in `go.mod`.
- `go.mod` sets the language/toolchain target to `go 1.26`.
- The process entrypoint is `cmd/cm/main.go`; it passes `os.Args[1:]`, `os.Stdin`, `os.Stdout`, and `os.Stderr` into `internal/cli.Main` and exits with the returned status code.
- Runtime build metadata is read in `internal/app/version.go` via `runtime/debug.ReadBuildInfo()`.
- `internal/app/version.go` exposes a package variable `Version = "dev"`, so release builds can override it with `-ldflags` as documented in `docs/development.md`.
- The CLI is intentionally local/terminal-oriented: there is no server, daemon, database, HTTP API, Dockerfile, or generated-code setup visible in the repository markers.

## Module and package layout

- `cmd/cm/main.go` is a thin binary wrapper around internal application code.
- `internal/app/options.go` and `internal/app/version.go` own process options and version metadata.
- `internal/cli/root.go` wires the command tree with Cobra and owns commands such as `status`, `diff`, `sync`, `add`, `apply`, `merge`, `edit`, `git`, `version`, and `completion`.
- `internal/cli/service.go` adapts CLI commands to the chezmoi client, source-git status, diff generation, and sync execution.
- `internal/chezmoi/client.go` wraps the external `chezmoi` binary behind a small client and a runner interface.
- `internal/chezmoi/status.go` parses `chezmoi status` output into status entries.
- `internal/chezmoi/content.go` loads rendered target content through `chezmoi cat` and local content through `os.Open`.
- `internal/process/runner.go` is the subprocess boundary using `os/exec` and provides injectable `Output` and `Run` methods.
- `internal/diff/diff.go` generates in-process unified diffs and caps files at `DefaultMaxFileSize = 1 << 20` bytes.
- `internal/tui/*.go` owns the Bubble Tea terminal interface, key handling, path display, lazy diff loading, and final confirmation flow.

## Direct dependencies from `go.mod`

- `charm.land/bubbles/v2 v2.1.0` — used by `internal/tui/model.go` and `internal/tui/keys.go` for help/key binding UI pieces.
- `charm.land/bubbletea/v2 v2.0.7` — used by `internal/tui/tui.go`, `internal/tui/update.go`, `internal/tui/view.go`, `internal/tui/diff_state.go`, and `internal/tui/confirm.go` as the TUI runtime.
- `charm.land/lipgloss/v2 v2.0.4` — used by `internal/tui/view.go` and `internal/tui/keys.go` for terminal styling.
- `github.com/fatih/color v1.19.0` — used by `internal/cli/status.go` for colored status headings and indicators.
- `github.com/spf13/cobra v1.10.2` — used by `internal/cli/root.go` for command parsing and shell completion generation.
- `golang.org/x/term v0.44.0` — used by `internal/tui/tui.go` to detect whether output is a terminal.

## Indirect dependencies recorded in `go.mod`

- Charm terminal support stack: `github.com/charmbracelet/colorprofile`, `github.com/charmbracelet/ultraviolet`, `github.com/charmbracelet/x/ansi`, `github.com/charmbracelet/x/term`, `github.com/charmbracelet/x/termios`, and `github.com/charmbracelet/x/windows`.
- Text width and segmentation support: `github.com/clipperhouse/displaywidth`, `github.com/clipperhouse/uax29/v2`, `github.com/mattn/go-runewidth`, and `github.com/rivo/uniseg`.
- Terminal/color compatibility support: `github.com/lucasb-eyer/go-colorful`, `github.com/mattn/go-colorable`, `github.com/mattn/go-isatty`, `github.com/muesli/cancelreader`, and `github.com/xo/terminfo`.
- Cobra support dependency: `github.com/spf13/pflag`; Windows console/mousetrap support appears as `github.com/inconshreveable/mousetrap`.
- Diff and internal helper dependency: `github.com/rogpeppe/go-internal`, used directly as `github.com/rogpeppe/go-internal/diff` in `internal/diff/diff.go` even though it is listed in the indirect block.
- Go extended packages: `golang.org/x/sync` and `golang.org/x/sys`.
- `go.sum` is present and records module checksums for reproducible module resolution.

## Go module and package tooling

- Dependency resolution is Go Modules only: `go.mod` declares the module path and versions, while `go.sum` pins checksums.
- The project keeps application code under `internal/`, which means packages such as `internal/cli`, `internal/chezmoi`, `internal/tui`, and `internal/process` are private to this module by Go's import rules.
- Tests follow Go's standard colocated `*_test.go` convention under the same packages, for example `internal/cli/cli_test.go`, `internal/chezmoi/client_test.go`, `internal/diff/diff_test.go`, and `internal/tui/sync_test.go`.
- There is no separate assertion/mocking test library in `go.mod`; tests visible in the tree use the standard `testing` package plus hand-written fakes.
- There is no checked-in formatter/linter configuration; formatting and vet/lint policy are therefore not encoded outside the standard Go toolchain files.

## Config, build, and tooling files

- `go.mod` and `go.sum` are the Go module configuration files at the repo root.
- `justfile` defines local install and release build recipes.
- `docs/development.md` documents development requirements, local execution, release build, and verification commands.
- `docs/development.md` documents release-version override with `go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/app.Version=v0.1.0"`.

## External command integrations

- `internal/process/runner.go` executes external commands through `exec.Command`; missing binaries are wrapped as `<command> not found in PATH`.
- `internal/chezmoi/client.go` defaults the managed binary name to `chezmoi` when `Client.Binary` is empty.
- `internal/chezmoi/client.go` runs `chezmoi status --path-style=absolute` for local reconciliation status.
- `internal/chezmoi/client.go` runs `chezmoi managed` for shell completion of editable managed files.
- `internal/chezmoi/content.go` obtains rendered target content through `chezmoi cat <target>` via `Client.OutputLimit`.
- `internal/cli/service.go` runs `chezmoi source-path` to locate the chezmoi source repository.
- `internal/cli/service.go` runs `git status --porcelain=v1` with `process.IO.Dir` set to the chezmoi source path.
- `internal/cli/service.go` runs `lazygit` in the chezmoi source directory for the `cm git` command.
- `internal/cli/service.go` delegates direct mutation commands to chezmoi: `chezmoi add`, `chezmoi apply`, `chezmoi merge`, and `chezmoi edit`.
- `internal/cli/service.go` uses `chezmoi apply --force` for confirmed apply actions in the sync flow.
- `internal/cli/root.go` generates shell completions through Cobra for `bash`, `zsh`, `fish`, and `powershell`.

## Release-adjacent stack notes

- Version output is built from Go build info in `internal/app/version.go`; it reports version, VCS revision, VCS time, dirty state, and Go version when present.
- The documented release override in `docs/development.md` targets `github.com/zhongyangchuwu/cm/internal/app.Version`, so the module path in `go.mod` is part of the release build contract.
- Because `justfile` installs with `go install ./cmd/cm`, the binary name comes from the `cmd/cm` directory.
- The install recipe discovers the binary destination with `go env GOBIN`, falling back to `go env GOPATH` plus `/bin` when `GOBIN` is empty.

## Platform and environment assumptions

- The runtime assumes a POSIX-like or PATH-based shell environment where `chezmoi`, `git`, and optionally `lazygit` are discoverable by `exec.Command`.
- `justfile` assumes `bash` and Unix-style environment variables such as `${HOME}`, `${GOBIN}`, and `${GOPATH}`.
- `justfile` installs zsh completion under `${HOME}/.zfunc/_cm`, so the default install path is zsh-oriented even though the CLI can generate other shell completions.
- `internal/cli/service.go` and `internal/tui/path.go` use `os.UserHomeDir()` and `filepath` to resolve/display home-relative paths.
- `internal/tui/path.go` formats paths under the user home as `~` plus `os.PathSeparator`, so displayed separators follow the build platform.
- `internal/tui/tui.go` changes behavior based on `golang.org/x/term.IsTerminal`; non-terminal output disables the Bubble Tea renderer and prints the final view string.
- `internal/process/runner.go` passes stdin/stdout/stderr through to external commands, which makes interactive commands such as `chezmoi merge` and `lazygit` depend on a usable terminal.
