# Codebase Structure

## Repository shape
- `go.mod` declares module `github.com/zhongyangchuwu/cm` and Go `1.26`.
- `cmd/cm/main.go` is the only executable entry point currently present under `cmd/`.
- `cmd/cm/main.go` delegates immediately to `internal/cli.Main` and exits with its return code.
- `internal/` contains all application packages; there are no exported library packages outside `internal/`.
- `internal/cli/` owns command construction, command rendering, and service adapters for external tools.
- `internal/chezmoi/` wraps the `chezmoi` executable and parses chezmoi-oriented output formats.
- `internal/syncdiff/` computes local-vs-chezmoi file diffs using byte content sources.
- `internal/reconcile/` defines shared reconciliation action types and the review service boundary.
- `internal/ui/` owns the interactive sync TUI, its model, key bindings, views, diff loading, and confirmation flow.
- `internal/process/` provides the process execution abstraction used by the chezmoi client and CLI services.
- `internal/build/` reads build metadata and formats version output.
- `docs/` contains human-facing repository documentation such as `docs/usage.md`, `docs/development.md`, and `docs/design.md`, but runtime code does not import from it.
- Tests live beside implementation files, for example `internal/cli/cli_test.go`, `internal/cli/service_test.go`, `internal/ui/sync_test.go`, and `internal/syncdiff/diff_test.go`.

## Directory layout
- `cmd/cm/main.go` is the binary package for the `cm` executable.
- `internal/build/info.go` is the build metadata package used by command version output.
- `internal/chezmoi/client.go`, `internal/chezmoi/status.go`, and `internal/chezmoi/content.go` form the chezmoi integration package.
- `internal/cli/cli.go`, `internal/cli/status.go`, `internal/cli/diff.go`, and `internal/cli/service.go` form the command package.
- `internal/process/runner.go` is the process execution package and has no CLI-specific command knowledge.
- `internal/reconcile/service.go` is a small domain contract package shared by command services and UI.
- `internal/syncdiff/diff.go` is the diff engine package.
- `internal/ui/model.go`, `internal/ui/tui.go`, `internal/ui/update.go`, `internal/ui/view.go`, `internal/ui/diff.go`, `internal/ui/confirm.go`, `internal/ui/keys.go`, and `internal/ui/path.go` form the sync TUI package.
- `internal/*/*_test.go` files stay co-located with the package they verify.
- `docs/` and `README.md` are documentation surfaces, not runtime packages.

## Repository support files
- `justfile` defines an `install` recipe that installs `./cmd/cm` and generates zsh completion via `go run ./cmd/cm completion zsh`.
- `go.sum` locks transitive module checksums for dependencies declared in `go.mod`.
- `.planning/codebase/` contains generated brownfield planning artifacts for release preparation and is not part of runtime code.
- `.planning/codebase/MAP.md` summarizes this codebase map set but may need refresh if line counts change after this structure pass.

## Entry points and command surface
- `cmd/cm/main.go` imports `github.com/zhongyangchuwu/cm/internal/cli` and passes `os.Args[1:]`, `os.Stdin`, `os.Stdout`, and `os.Stderr` into `cli.Main`.
- `internal/cli/cli.go` defines `Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int` as the testable application entry point.
- `internal/cli/cli.go` builds a Cobra root command with `Use: "cm"`, short description `Chezmoi reconciliation manager`, and version from `internal/build`.
- Running `cm` with no subcommand executes the same read-only status rendering path as `cm status`.
- `cm status [target...]` is defined in `internal/cli/cli.go` and rendered by `internal/cli/status.go`.
- `cm diff [target...]` is defined in `internal/cli/cli.go` and rendered by `internal/cli/diff.go` using internal sync diffs.
- `cm sync [target...]` is defined in `internal/cli/cli.go` and launches `ui.RunSyncTUI` from `internal/ui/tui.go`.
- `cm add [target...]` forwards to `commandServices.Target.AddTargets` in `internal/cli/service.go`.
- `cm apply [target...]` forwards to `commandServices.Target.ApplyTargets` in `internal/cli/service.go`.
- `cm merge [target...]` forwards to `commandServices.Target.MergeTargets` in `internal/cli/service.go`.
- `cm edit <target>` forwards to `commandServices.Edit.EditTarget` and offers managed-file shell completion through `ManagedFiles`.
- `cm git` forwards to `commandServices.SourceGit.OpenSourceGit` and opens `lazygit` in the chezmoi source repository.
- `cm version` prints detailed build metadata from `internal/build/info.go`.
- `cm completion [bash|zsh|fish|powershell]` emits Cobra-generated shell completion scripts.

## Command-to-file map
- Root `cm`: `internal/cli/cli.go` builds the root command and `internal/cli/status.go` renders the default status behavior.
- `cm status [target...]`: `internal/cli/cli.go` routes to `renderStatus` in `internal/cli/status.go`.
- `cm diff [target...]`: `internal/cli/cli.go` routes to `renderDiff` in `internal/cli/diff.go`.
- `cm sync [target...]`: `internal/cli/cli.go` routes to `RunSyncTUI` in `internal/ui/tui.go`.
- `cm add [target...]`: `internal/cli/cli.go` routes to `AddTargets` in `internal/cli/service.go`.
- `cm apply [target...]`: `internal/cli/cli.go` routes to `ApplyTargets` in `internal/cli/service.go`.
- `cm merge [target...]`: `internal/cli/cli.go` routes to `MergeTargets` in `internal/cli/service.go`.
- `cm edit <target>`: `internal/cli/cli.go` routes to `EditTarget` and `ManagedFiles` in `internal/cli/service.go`.
- `cm git`: `internal/cli/cli.go` routes to `OpenSourceGit` in `internal/cli/service.go`.
- `cm version`: `internal/cli/cli.go` prints metadata formatted by `internal/build/info.go`.
- `cm completion [bash|zsh|fish|powershell]`: `internal/cli/cli.go` delegates to Cobra completion generators.

## Runtime entry and support entry points
- Runtime process entry is `main()` in `cmd/cm/main.go`.
- Testable CLI entry is `Main` in `internal/cli/cli.go`.
- Cobra command construction entry is `newRootCommand` in `internal/cli/cli.go`.
- Concrete service wiring entry is `commandServicesFor` in `internal/cli/service.go`.
- Interactive UI entry is `RunSyncTUI` in `internal/ui/tui.go`.
- Diff engine entry is `Differ.Diff` in `internal/syncdiff/diff.go`.
- Chezmoi process wrapper entries are `Client.Status`, `Client.Output`, `Client.OutputLimit`, `Client.Run`, and `Client.ManagedFiles` in `internal/chezmoi/client.go`.
- Build metadata entry is `Current` in `internal/build/info.go`.

## Package inventory
- `internal/cli/cli.go` contains Cobra command construction and dependency injection through `commandServices`.
- `internal/cli/status.go` formats local chezmoi differences and source repository git changes for terminal output.
- `internal/cli/diff.go` expands an empty target list into status paths, then writes each computed diff.
- `internal/cli/service.go` adapts the `chezmoi.Client` to command interfaces used by Cobra handlers and the TUI.
- `internal/chezmoi/client.go` defines `Client`, default binary resolution, command execution, output limits, and runner wiring.
- `internal/chezmoi/client.go` runs `chezmoi status --path-style=absolute` so status entries carry absolute target paths.
- `internal/chezmoi/status.go` parses `chezmoi status` lines into `StatusEntry` values and parses managed file lists.
- `internal/chezmoi/content.go` implements `ContentLoader` for `syncdiff.ContentSource` by reading chezmoi target content and local files.
- `internal/syncdiff/diff.go` defines `ContentSource`, `Differ`, size limits, binary detection, and unified diff generation.
- `internal/reconcile/service.go` defines `ActionKind`, `Action`, string/marker labels, and the `ReviewService` interface.
- `internal/ui/model.go` defines the sync TUI state model, pending actions, cursor movement, focus state, and target selection.
- `internal/ui/tui.go` starts the Bubble Tea program and handles terminal-vs-nonterminal rendering behavior.
- `internal/ui/update.go` maps key presses to review, confirm, execute, cursor, and diff-loading state transitions.
- `internal/ui/diff.go` loads diffs asynchronously through the review service and styles diff lines for display.
- `internal/ui/confirm.go` rechecks status before executing pending actions and skips actions for clean targets.
- `internal/ui/keys.go` centralizes TUI key bindings and Lip Gloss styles.
- `internal/ui/path.go` shortens paths under the current home directory to `~` for display.
- `internal/ui/view.go` renders panes, footer help, visible entries, truncation, and line padding.
- `internal/process/runner.go` defines `Runner`, `IO`, and the `ExecRunner` implementation over `os/exec`.
- `internal/build/info.go` collects VCS and Go build metadata via `runtime/debug.ReadBuildInfo`.

## Boundary notes
- `cmd/cm/main.go` is intentionally thin; all CLI behavior is inside `internal/cli`.
- `internal/cli` depends on `internal/build`, `internal/chezmoi`, `internal/reconcile`, `internal/syncdiff`, `internal/process`, and `internal/ui`.
- `internal/ui` depends only on the reconciliation contract and chezmoi status entries rather than on Cobra.
- `internal/reconcile` has no process or UI implementation; it is the shared data contract between CLI services and UI.
- `internal/syncdiff` depends on a content-source interface rather than directly on `chezmoi.Client`.
- `internal/chezmoi` depends on `internal/process` so external command execution can be replaced in tests.
- `internal/process` is the lowest-level package for launching external processes.
- External libraries in use include Cobra for commands, Bubble Tea/Bubbles/Lip Gloss for TUI, Fatih Color for status output, and `rogpeppe/go-internal/diff` for diff text.

## Notable implementation files
- `internal/cli/service_test.go` shows the CLI service layer is intended to be exercised with a recording runner rather than real `chezmoi`, `git`, or `lazygit` invocations.
- `internal/cli/cli_test.go` exercises command routing through fake services, confirming `newRootCommand` is the command-surface seam.
- `internal/ui/sync_test.go` covers the sync TUI interaction model rather than placing that behavior in Cobra command tests.
- `internal/chezmoi/client_test.go` and `internal/chezmoi/status_test.go` keep external command parsing and client behavior close to the wrapper package.
- `internal/syncdiff/diff_test.go` keeps diff behavior isolated from the CLI and TUI packages.
- There is no separate `pkg/` tree, plugin tree, or multiple binary tree; release structure is currently a single `cm` binary backed by internal packages.

## Release-facing structure observations
- The public executable surface is concentrated in `cmd/cm/main.go` and `internal/cli/cli.go`, which makes command review for v0.1.0 localized.
- User-visible terminal strings are spread across `internal/cli/cli.go`, `internal/cli/status.go`, `internal/cli/diff.go`, `internal/ui/view.go`, and `internal/ui/keys.go`.
- External binary names are concentrated in `internal/chezmoi/client.go` for `chezmoi`, `internal/cli/service.go` for `git` and `lazygit`, and Cobra shell names in `internal/cli/cli.go`.
- File-size and binary-diff behavior lives in `internal/syncdiff/diff.go`, not in command handlers.
- Home-directory display behavior lives in `internal/ui/path.go`, while home-directory edit target expansion lives in `internal/cli/service.go`.
- The codebase currently has one clear command boundary, one UI boundary, one process boundary, and one domain action boundary.
- Because implementation packages are under `internal/`, external consumers cannot import these packages as public Go APIs.
- Release preparation should treat CLI commands and generated completion behavior as the compatibility surface rather than Go package APIs.
