# Codebase Structure

## Repository shape

- `go.mod` declares module `github.com/zhongyangchuwu/cm` and Go `1.26`.
- `cmd/cm/main.go` is the only executable entry point currently present under `cmd/`.
- `cmd/cm/main.go` delegates immediately to `internal/cli.Main` and exits with its return code.
- `internal/` contains all application packages; there are no exported library packages outside `internal/`.
- `internal/app/` owns process options and version metadata.
- `internal/cli/` owns command construction, command rendering, and service adapters for external tools.
- `internal/chezmoi/` wraps the `chezmoi` executable and parses chezmoi-oriented output formats.
- `internal/diff/` computes local-vs-chezmoi file diffs using byte content sources.
- `internal/tui/` owns the interactive sync TUI, its model, key bindings, views, diff loading, and confirmation flow.
- `internal/process/` provides the process execution abstraction used by the chezmoi client and CLI services.
- `docs/` contains human-facing repository documentation such as `docs/usage.md`, `docs/development.md`, and `docs/design.md`, but runtime code does not import from it.
- Tests live beside implementation files, for example `internal/cli/cli_test.go`, `internal/cli/service_test.go`, `internal/tui/sync_test.go`, and `internal/diff/diff_test.go`.

## Directory layout

- `cmd/cm/main.go` is the binary package for the `cm` executable.
- `internal/app/options.go` contains process-wide CLI options.
- `internal/app/version.go` is the version metadata package used by command version output.
- `internal/chezmoi/client.go`, `internal/chezmoi/status.go`, and `internal/chezmoi/content.go` form the chezmoi integration package.
- `internal/cli/root.go`, `internal/cli/status.go`, `internal/cli/diff_cmd.go`, and `internal/cli/service.go` form the command package.
- `internal/process/runner.go` is the process execution package and has no CLI-specific command knowledge.
- `internal/diff/diff.go` is the diff engine package.
- `internal/tui/model.go`, `internal/tui/tui.go`, `internal/tui/update.go`, `internal/tui/view.go`, `internal/tui/diff_state.go`, `internal/tui/diff_view.go`, `internal/tui/confirm.go`, `internal/tui/keys.go`, and `internal/tui/path.go` form the sync TUI package.
- `internal/*/*_test.go` files stay co-located with the package they verify.
- `docs/` and `README.md` are documentation surfaces, not runtime packages.

## Entry points and command surface

- `cmd/cm/main.go` imports `github.com/zhongyangchuwu/cm/internal/cli` and passes `os.Args[1:]`, `os.Stdin`, `os.Stdout`, and `os.Stderr` into `cli.Main`.
- `internal/cli/root.go` defines `Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int` as the testable application entry point.
- `internal/cli/root.go` builds a Cobra root command with `Use: "cm"`, short description `Chezmoi reconciliation manager`, and version from `internal/app`.
- Running `cm` with no subcommand executes the same read-only status rendering path as `cm status`.
- `cm status [target...]` is rendered by `internal/cli/status.go`.
- `cm diff [target...]` is rendered by `internal/cli/diff_cmd.go` using internal diffs.
- `cm sync [target...]` launches `tui.RunSyncTUI` from `internal/tui/tui.go`.
- `cm add [target...]`, `cm apply [target...]`, and `cm merge [target...]` forward to command service methods in `internal/cli/service.go`.
- `cm edit <target>` forwards to `EditTarget` and offers managed-file shell completion through `ManagedFiles`.
- `cm git` forwards to `OpenSourceGit` and opens `lazygit` in the chezmoi source repository.
- `cm version` prints detailed build metadata from `internal/app/version.go`.
- `cm completion [bash|zsh|fish|powershell]` emits Cobra-generated shell completion scripts.

## Runtime entry and support entry points

- Runtime process entry is `main()` in `cmd/cm/main.go`.
- Testable CLI entry is `Main` in `internal/cli/root.go`.
- Cobra command construction entry is `newRootCommand` in `internal/cli/root.go`.
- Concrete service wiring entry is `commandServicesFor` in `internal/cli/service.go`.
- Interactive TUI entry is `RunSyncTUI` in `internal/tui/tui.go`.
- Diff engine entry is `Differ.Diff` in `internal/diff/diff.go`.
- Chezmoi process wrapper entries are `Client.Status`, `Client.Output`, `Client.OutputLimit`, `Client.Run`, and `Client.ManagedFiles` in `internal/chezmoi/client.go`.
- Build metadata entry is `Current` in `internal/app/version.go`.

## Boundary notes

- `cmd/cm/main.go` is intentionally thin; all CLI behavior is inside `internal/cli`.
- `internal/cli` depends on `internal/app`, `internal/chezmoi`, `internal/diff`, `internal/process`, and `internal/tui`.
- `internal/tui` owns terminal state and uses the review service interface rather than knowing concrete process commands.
- `internal/diff` depends on a content-source interface rather than directly on `chezmoi.Client`.
- `internal/chezmoi` depends on `internal/process` so external command execution can be replaced in tests.
- Release preparation should treat CLI commands and generated completion behavior as the compatibility surface rather than Go package APIs.
