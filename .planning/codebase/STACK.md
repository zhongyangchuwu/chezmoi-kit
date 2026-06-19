# Stack Map

## Language and runtime

- The project is a Go CLI module declared as `module github.com/zhongyangchuwu/cm` in `go.mod`.
- `go.mod` sets the language/toolchain target to `go 1.26`.
- The process entrypoint is `cmd/cm/main.go`; it passes `os.Args[1:]`, `os.Stdin`, `os.Stdout`, and `os.Stderr` into `internal/cli.Main` and exits with the returned status code.
- Runtime build metadata is read in `internal/build/info.go` via `runtime/debug.ReadBuildInfo()`.
- `internal/build/info.go` exposes a package variable `Version = "dev"`, so release builds can override it with `-ldflags` as documented in `docs/development.md`.
- The CLI is intentionally local/terminal-oriented: there is no server, daemon, database, HTTP API, Dockerfile, or generated-code setup visible in the repository markers.

## Module and package layout

- `cmd/cm/main.go` is a thin binary wrapper around internal application code.
- `internal/cli/cli.go` wires the command tree with Cobra and owns commands such as `status`, `diff`, `sync`, `add`, `apply`, `merge`, `edit`, `git`, `version`, and `completion`.
- `internal/cli/service.go` adapts CLI commands to the chezmoi client, source-git status, sync diff generation, and reconciliation execution.
- `internal/chezmoi/client.go` wraps the external `chezmoi` binary behind a small client and a runner interface.
- `internal/chezmoi/status.go` parses `chezmoi status` output into status entries.
- `internal/chezmoi/content.go` loads rendered target content through `chezmoi cat` and local content through `os.Open`.
- `internal/process/runner.go` is the subprocess boundary using `os/exec` and provides injectable `Output` and `Run` methods.
- `internal/syncdiff/diff.go` generates in-process unified diffs and caps files at `DefaultMaxFileSize = 1 << 20` bytes.
- `internal/reconcile/service.go` defines reconciliation actions and the review service contract consumed by the UI.
- `internal/ui/*.go` owns the Bubble Tea terminal interface, key handling, path display, lazy diff loading, and final confirmation flow.
- `internal/build/info.go` owns `cm version` build information formatting.

## Direct dependencies from `go.mod`

- `charm.land/bubbles/v2 v2.1.0` — used by `internal/ui/model.go` and `internal/ui/keys.go` for help/key binding UI pieces.
- `charm.land/bubbletea/v2 v2.0.7` — used by `internal/ui/tui.go`, `internal/ui/update.go`, `internal/ui/view.go`, `internal/ui/diff.go`, and `internal/ui/confirm.go` as the TUI runtime.
- `charm.land/lipgloss/v2 v2.0.4` — used by `internal/ui/view.go` and `internal/ui/keys.go` for terminal styling.
- `github.com/fatih/color v1.19.0` — used by `internal/cli/status.go` for colored status headings and indicators.
- `github.com/spf13/cobra v1.10.2` — used by `internal/cli/cli.go` for command parsing and shell completion generation.
- `golang.org/x/term v0.44.0` — used by `internal/ui/tui.go` to detect whether output is a terminal.

## Indirect dependencies recorded in `go.mod`

- Charm terminal support stack: `github.com/charmbracelet/colorprofile v0.4.3`, `github.com/charmbracelet/ultraviolet v0.0.0-20260525132238-948f4557a654`, `github.com/charmbracelet/x/ansi v0.11.7`, `github.com/charmbracelet/x/term v0.2.2`, `github.com/charmbracelet/x/termios v0.1.1`, and `github.com/charmbracelet/x/windows v0.2.2`.
- Text width and segmentation support: `github.com/clipperhouse/displaywidth v0.11.0`, `github.com/clipperhouse/uax29/v2 v2.7.0`, `github.com/mattn/go-runewidth v0.0.23`, and `github.com/rivo/uniseg v0.4.7`.
- Terminal/color compatibility support: `github.com/lucasb-eyer/go-colorful v1.4.0`, `github.com/mattn/go-colorable v0.1.14`, `github.com/mattn/go-isatty v0.0.20`, `github.com/muesli/cancelreader v0.2.2`, and `github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e`.
- Cobra support dependency: `github.com/spf13/pflag v1.0.9`; Windows console/mousetrap support appears as `github.com/inconshreveable/mousetrap v1.1.0`.
- Diff and internal helper dependency: `github.com/rogpeppe/go-internal v1.15.0`, used directly as `github.com/rogpeppe/go-internal/diff` in `internal/syncdiff/diff.go` even though it is listed in the indirect block.
- Go extended packages: `golang.org/x/sync v0.20.0` and `golang.org/x/sys v0.46.0`.
- `go.sum` is present and records module checksums for reproducible module resolution.

## Go module and package tooling

- Dependency resolution is Go Modules only: `go.mod` declares the module path and versions, while `go.sum` pins checksums.
- No `vendor/` directory or `go.work` workspace file was found in the repository tree, so local development appears to use normal module mode rather than vendoring or a multi-module workspace.
- The project keeps application code under `internal/`, which means packages such as `internal/cli`, `internal/chezmoi`, `internal/ui`, and `internal/process` are private to this module by Go's import rules.
- Tests follow Go's standard colocated `*_test.go` convention under the same packages, for example `internal/cli/cli_test.go`, `internal/chezmoi/client_test.go`, `internal/syncdiff/diff_test.go`, and `internal/ui/sync_test.go`.
- There is no separate assertion/mocking test library in `go.mod`; tests visible in the tree use the standard `testing` package plus hand-written fakes.
- There is no checked-in formatter/linter configuration; formatting and vet/lint policy are therefore not encoded outside the standard Go toolchain files.

## Config, build, and tooling files

- `go.mod` and `go.sum` are the only Go module configuration files found at the repo root.
- No `go.work`, `.go-version`, `Makefile`, `Dockerfile`, `.github/`, YAML/TOML/JSON config, or visible lint configuration file was found by filename lookup.
- `justfile` is the only task runner file found; it sets `bash -eu -o pipefail -c` as the recipe shell.
- `justfile` defines one recipe, `install`, which runs `go install ./cmd/cm`, creates `${HOME}/.zfunc`, runs `go run ./cmd/cm completion zsh > ${HOME}/.zfunc/_cm`, and prints the installed binary/completion paths.
- `docs/development.md` documents development requirements as Go, chezmoi, and just.
- `docs/development.md` documents project-wide tests and targeted test commands, but this stack mapping did not run any commands.
- `docs/development.md` documents local execution examples using `go run ./cmd/cm status`, `go run ./cmd/cm sync`, and `go run ./cmd/cm diff ~/.zshrc`.
- `docs/development.md` documents release-version override with `go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/build.Version=v0.1.0"`.

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
- `internal/cli/cli.go` generates shell completions through Cobra for `bash`, `zsh`, `fish`, and `powershell`.

## External command matrix

- `cm status` and bare `cm` flow through `internal/cli/cli.go` to `internal/cli/service.go`, then use `chezmoi status --path-style=absolute` plus `git status --porcelain=v1` in the directory returned by `chezmoi source-path`.
- `cm diff [target...]` is not a shell-out to `chezmoi diff`; `internal/cli/service.go` calls `internal/syncdiff/diff.go`, which compares `chezmoi cat <target>` output to bytes read from the local target path.
- `cm sync [target...]` uses the same status and diff plumbing as `status`/`diff`, then executes selected actions through `chezmoi add`, `chezmoi apply --force`, or per-target `chezmoi merge` from `internal/cli/service.go`.
- `cm add`, `cm apply`, and `cm merge` are direct wrappers around the matching chezmoi subcommands, implemented by `runTargets` in `internal/cli/service.go`.
- `cm edit <target>` converts the provided target to a home-rooted absolute path with `os.UserHomeDir()` and `filepath.Join()` before running `chezmoi edit`.
- `cm git` does not run git commands directly beyond status; it launches `lazygit` with the working directory set to the chezmoi source path.
- `cm completion` is implemented entirely through Cobra generation methods in `internal/cli/cli.go`; `justfile` consumes this by generating zsh completion during `just install`.

## Release-adjacent stack notes

- Version output is built from Go build info in `internal/build/info.go`; it reports version, VCS revision, VCS time, dirty state, and Go version when present.
- The documented release override in `docs/development.md` targets `github.com/zhongyangchuwu/cm/internal/build.Version`, so the module path in `go.mod` is part of the release build contract.
- Because `justfile` installs with `go install ./cmd/cm`, the binary name comes from the `cmd/cm` directory.
- The install recipe discovers the binary destination with `go env GOBIN`, falling back to `go env GOPATH` plus `/bin` when `GOBIN` is empty.
- The repository has no CI configuration in `.github/`, so release checks are not encoded as GitHub Actions in the current tree.
- The repository has no visible package-manager files beyond Go Modules and `justfile`; there is no npm, Python, Rust, or container build stack present in the root-level file set.

## Platform and environment assumptions

- The runtime assumes a POSIX-like or PATH-based shell environment where `chezmoi`, `git`, and optionally `lazygit` are discoverable by `exec.Command`.
- `justfile` assumes `bash` and Unix-style environment variables such as `${HOME}`, `${GOBIN}`, and `${GOPATH}`.
- `justfile` installs zsh completion under `${HOME}/.zfunc/_cm`, so the default install path is zsh-oriented even though the CLI can generate other shell completions.
- `internal/cli/service.go` and `internal/ui/path.go` use `os.UserHomeDir()` and `filepath` to resolve/display home-relative paths.
- `internal/ui/path.go` formats paths under the user home as `~` plus `os.PathSeparator`, so displayed separators follow the build platform.
- `internal/ui/tui.go` changes behavior based on `golang.org/x/term.IsTerminal`; non-terminal output disables the Bubble Tea renderer and prints the final view string.
- `internal/process/runner.go` passes stdin/stdout/stderr through to external commands, which makes interactive commands such as `chezmoi merge` and `lazygit` depend on a usable terminal.
- `docs/usage.md` and `README.md` describe a personal configuration-management workflow centered on local files, the chezmoi source repository, and terminal interaction.

## Standard-library runtime surfaces

- `cmd/cm/main.go` uses `os.Exit`, `os.Args`, and standard streams directly, so process exit behavior is centralized at the binary boundary.
- `internal/process/runner.go` uses `os/exec` rather than a shell wrapper; command arguments are passed as argv slices from Go code.
- `internal/process/runner.go` defaults missing stdin/stdout/stderr to `os.Stdin`, `os.Stdout`, and `os.Stderr`, preserving interactive subprocess behavior when callers do not inject streams.
- `internal/chezmoi/content.go` reads local target files with `os.Open` and a bounded buffer, not by invoking `cat` or another external file command.
- `internal/syncdiff/diff.go` uses byte-level checks for oversized and binary content before invoking the Go diff library.
- `internal/cli/service.go` trims `chezmoi source-path` output with `strings.TrimSpace` before using it as a subprocess working directory.
- `internal/cli/service.go` parses source repository status from `git status --porcelain=v1` lines using the two-character status code and path suffix.
- `internal/ui/tui.go` passes `tea.WithoutSignals()`, so the Bubble Tea program is not relying on Bubble Tea signal handling for this CLI path.

## Configuration surface summary

- The stack has no repository-level application configuration file; runtime behavior is driven by CLI args, the user's chezmoi configuration, Git state in the chezmoi source repository, PATH lookup, and terminal/stdstream capabilities.
- `README.md` frames `cm` as a read-only-by-default helper around local config, chezmoi source state, and chezmoi git history.
- `docs/usage.md` documents mutation boundaries: `cm status`, `cm diff`, `cm version`, and completion generation are non-mutating, while `cm sync`, `cm add`, `cm apply`, `cm merge`, and `cm git` can mutate user state through external tools.
- No Go build tags were found in `cmd/` or `internal/`, so there are no checked-in platform-specific Go source variants in the current code paths.
- No vendored dependencies were found in the repository file lookup; the dependency graph is expected to resolve from module metadata instead of checked-in source copies.
- No generated-code markers or generator config were found in the stack files inspected; command completions are generated at runtime/install time by Cobra rather than committed as generated files.
- The planning artifact for this stack map is `.planning/codebase/STACK.md`; no source, documentation outside `.planning/codebase/`, `go.mod`, `justfile`, or tests were edited for this assignment.
