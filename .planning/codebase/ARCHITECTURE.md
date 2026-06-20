# Architecture Map

## High-level design

- The application is a Go command-line wrapper around `chezmoi` with extra reconciliation views and commands.
- `cmd/cm/main.go` is a process bootstrapper only; it passes process args and standard streams to `internal/cli.Main`.
- `internal/cli/root.go` constructs the Cobra command tree, creates a `chezmoi.Client`, wraps it with `app.NewServices`, and returns an exit code.
- `internal/cli/service.go` aliases the app service graph for command wiring only.
- `internal/app/services.go` owns the service graph, status/diff/sync/source-git/edit/direct-wrapper use cases, and sync action contract.
- `internal/app` also owns process-wide options and runtime version metadata.
- The core external dependency is the `chezmoi` executable, wrapped by `internal/chezmoi/client.go`.
- Process execution is abstracted by `internal/process/runner.go`, which allows tests to record commands without invoking real binaries.
- The interactive workflow is isolated in `internal/tui/` and consumes `app.SyncService` plus `app.Action` values.
- Internal diffs are computed in `internal/diff/diff.go` from abstract file content rather than by shelling out to a diff command.

## Module and package boundaries

- `internal/cli` owns command surface, user-visible command routing, stream plumbing, and exit behavior.
- `internal/app` owns process-wide runtime options, `cm version` metadata, app service composition, app use cases, and sync action/domain contracts.
- `internal/chezmoi` owns the protocol for calling `chezmoi`, including `status`, `source-path`, `cat`, `managed`, and direct mutating invocations.
- `internal/process` owns the lowest-level `exec.Command` integration and default stdin/stdout/stderr behavior.
- `internal/tui` owns stateful terminal interaction: selected file, focused pane, pending actions, confirmation mode, diff cache, and execution messages.
- `internal/diff` owns comparison rules: max file size, binary detection, local missing-file behavior, and unified diff output.
- The package graph points inward from CLI/TUI adapters to app services and infrastructure capabilities; there is no dependency from `internal/chezmoi` back into `internal/cli`, `internal/app`, or `internal/tui`.
- The data contract crossing most package boundaries is `chezmoi.StatusEntry` from `internal/chezmoi/status.go`.
- The sync action contract crossing the TUI/service boundary is `app.Action` from `internal/app/services.go`.

## Dependency direction

- `cmd/cm/main.go` depends on `internal/cli` and nothing else in the application.
- `internal/cli/root.go` depends on Cobra plus internal packages for app version info, app service construction, chezmoi access, and TUI startup.
- `internal/app/services.go` depends on `internal/chezmoi`, `internal/process`, and `internal/diff` to implement command use cases.
- `internal/tui` depends on Bubble Tea libraries, `internal/app` contracts, and `internal/chezmoi` status entries.
- `internal/diff` depends on a content interface and the `rogpeppe/go-internal/diff` formatter, not on CLI packages.
- `internal/chezmoi` depends on `internal/process` but does not depend on TUI or Cobra.
- `internal/process` sits at the bottom of the application dependency graph.

## Interface boundaries

- `app.StatusService` exposes semantic status reports for read-only command output.
- `app.DiffService` exposes semantic diff reports for command output.
- `app.TargetCommandService` exposes explicit mutating wrappers for add, apply, and merge.
- `app.SourceGitService` isolates the `cm git` behavior behind `OpenSourceGit`.
- `app.EditService` groups edit execution and completion source lookup.
- `app.SyncService` is intentionally broader than read-only status because the TUI must status, diff, and execute confirmed actions.
- `report.Document` is the semantic output contract for status, diff, version, Markdown, and ANSI/plain CLI rendering.
- `diff.ContentSource` in `internal/diff/diff.go` is byte-oriented so diffing does not know whether content came from chezmoi, disk, or a fake source.
- `process.Runner` in `internal/process/runner.go` has separate `Output` and `Run` methods to distinguish captured-output calls from streaming command execution.

## Data flow: status command

- `cmd/cm/main.go` calls `cli.Main`.
- `cli.Main` creates `chezmoi.Client{Stdin, Stdout, Stderr}` and passes it to `app.NewServices`.
- The root command or `status` subcommand calls `renderStatus` in `internal/cli/status.go`.
- `renderStatus` renders `app.StatusService.StatusReport(targets)` with `internal/report` plain/ANSI rendering.
- The app service delegates local status to `chezmoi.Client.Status`.
- `chezmoi.Client.Status` runs `chezmoi status --path-style=absolute` and parses output with `chezmoi.ParseStatus`.
- The app service also includes source repository git changes.
- `app.service.SourceStatus` runs `chezmoi source-path`, then runs `git status --porcelain=v1` in that directory through `process.Runner`.
- Status output prints `clean` when both local status and source git status are empty.

## Data flow: diff command

- `cm diff [target...]` routes through `renderDiff` in `internal/cli/diff_cmd.go`.
- `renderDiff` renders `app.DiffService.DiffReport(targets)` with `internal/report` plain/ANSI rendering.
- If no targets are provided, the app diff service calls `Status(nil)` and uses each `StatusEntry.Path` as a diff target.
- For each target, the app diff service constructs `diff.Differ{Source: chezmoi.ContentLoader{Client: s.client}}`.
- `diff.Differ.Diff` reads target content from `chezmoi cat <target>` through `ContentLoader.TargetContent`.
- `ContentLoader.LocalContent` in `internal/chezmoi/content.go` reads the local target path directly with a byte limit.
- `diff.DiffBytes` rejects oversized files, reports binary differences, returns `no diff` for identical content, or emits a unified diff.
- `report.ClassifyDiffLine` classifies unified diff lines once for CLI report rendering and TUI diff styling.

## Data flow: interactive sync

- `cm sync [target...]` routes from `internal/cli/root.go` to `tui.RunSyncTUI` in `internal/tui/tui.go`.
- `RunSyncTUI` first calls `app.SyncService.Status(targets)` and exits early with `clean` when there are no entries.
- `newSyncTUIModel` in `internal/tui/model.go` stores the service, status entries, diff cache, help model, and home directory.
- `internal/tui/update.go` maps navigation keys, pending action keys, diff refresh, confirmation, and quit behavior.
- `internal/tui/diff_state.go` loads a selected target diff by calling `app.SyncService.DiffOutput(target)` asynchronously through a Bubble Tea command.
- `internal/tui/diff_view.go` styles loaded diff lines for display.
- `internal/tui/confirm.go` rechecks fresh status for selected targets before execution to avoid applying actions to targets that became clean.
- Confirmed actions are sent to `app.SyncService.ExecuteNonInteractive` or `app.SyncService.TerminalCommand` and represented as `executeMsg` values.
- The TUI view in `internal/tui/view.go` renders file, diff, confirm, footer, and message panes from model state.

## Design patterns and testing seams

- Dependency injection is done with app-owned interfaces in `internal/app/services.go` and command-level fakes in adapter tests.
- `app.Services` groups those interfaces so `newRootCommand` can be tested with fakes instead of real `chezmoi`, `git`, or `lazygit` processes.
- `app.SyncService` is the package-level boundary that keeps the TUI independent from concrete service implementation.
- `diff.ContentSource` decouples diffing logic from the chezmoi client and local filesystem details.
- `process.Runner` decouples command execution from command construction and gives tests a recording seam.
- `report.Options` controls ANSI color: auto color requires TTY output, `NO_COLOR` disables ANSI when non-empty, and explicit always/never modes are available internally for tests and future CLI flags.
- Most functions accept `io.Reader` and `io.Writer` values, visible in `cmd/cm/main.go`, `internal/cli/root.go`, and `internal/tui/tui.go`, which keeps CLI behavior testable without global standard streams.
