# Architecture Map

## High-level design
- The application is a Go command-line wrapper around `chezmoi` with extra reconciliation views and commands.
- `cmd/cm/main.go` is a process bootstrapper only; it does not parse flags, perform I/O directly beyond passing standard streams, or own business behavior.
- `internal/cli/cli.go` is the composition root for CLI execution: it constructs a `chezmoi.Client`, wraps it in `chezmoiService`, creates Cobra commands, and returns an exit code.
- Cobra command handlers in `internal/cli/cli.go` are intentionally thin and delegate to renderer functions, target services, or the TUI.
- `internal/cli/service.go` is the main adapter layer between command intent and external tools.
- The core external dependency is the `chezmoi` executable, wrapped by `internal/chezmoi/client.go`.
- Process execution is abstracted by `internal/process/runner.go`, which allows tests to record commands without invoking real binaries.
- The interactive workflow is isolated in `internal/ui/` and communicates through `reconcile.ReviewService` from `internal/reconcile/service.go`.
- Internal diffs are computed in `internal/syncdiff/diff.go` from abstract file content rather than by shelling out to a diff command.
- Build/version metadata is isolated in `internal/build/info.go` and consumed only by the CLI root/version commands.

## Module and package boundaries
- `internal/cli` owns command surface, user-visible command routing, terminal text rendering, and concrete service composition.
- `internal/chezmoi` owns the protocol for calling `chezmoi`, including `status`, `source-path`, `cat`, `managed`, and direct mutating invocations.
- `internal/process` owns the lowest-level `exec.Command` integration and default stdin/stdout/stderr behavior.
- `internal/reconcile` owns only domain vocabulary: action kinds, action labels, markers, and the review service interface.
- `internal/ui` owns stateful interaction: selected file, focused pane, pending actions, confirmation mode, diff cache, and execution messages.
- `internal/syncdiff` owns comparison rules: max file size, binary detection, local missing-file behavior, and unified diff output.
- `internal/build` owns build metadata extraction and formatting.
- The package graph points inward from UI/CLI to contracts and adapters; there is no dependency from `internal/chezmoi` back into `internal/cli` or `internal/ui`.
- The data contract crossing most package boundaries is `chezmoi.StatusEntry` from `internal/chezmoi/status.go`.
- The action contract crossing the UI/service boundary is `reconcile.Action` from `internal/reconcile/service.go`.

## Dependency direction
- `cmd/cm/main.go` depends on `internal/cli` and nothing else in the application.
- `internal/cli/cli.go` depends on Cobra plus internal packages for build info, chezmoi access, and UI startup.
- `internal/cli/service.go` depends on `internal/chezmoi`, `internal/process`, `internal/reconcile`, and `internal/syncdiff` to implement command services.
- `internal/ui` depends on Bubble Tea libraries, `internal/reconcile`, and `internal/chezmoi` status entries.
- `internal/syncdiff` depends on a content interface and the `rogpeppe/go-internal/diff` formatter, not on CLI packages.
- `internal/chezmoi` depends on `internal/process` but does not depend on UI or Cobra.
- `internal/reconcile` depends only on `internal/chezmoi` for status entry types and avoids taking a dependency on command or process code.
- `internal/process` sits at the bottom of the application dependency graph.

## Interface boundaries
- `statusService` in `internal/cli/service.go` exposes `Status` and `SourceStatus` for read-only status rendering.
- `diffService` in `internal/cli/service.go` exposes `Status` and `DiffOutput` for the diff command's target expansion and diff rendering.
- `targetCommandService` in `internal/cli/service.go` exposes explicit mutating wrappers for add, apply, and merge.
- `sourceGitService` in `internal/cli/service.go` isolates the `cm git` behavior behind `OpenSourceGit`.
- `editService` in `internal/cli/service.go` groups edit execution and completion source lookup.
- `reconcile.ReviewService` in `internal/reconcile/service.go` is intentionally broader than read-only status because the TUI must status, diff, and execute confirmed actions.
- `syncdiff.ContentSource` in `internal/syncdiff/diff.go` is intentionally byte-oriented so diffing does not know whether content came from chezmoi, disk, or a fake source.
- `process.Runner` in `internal/process/runner.go` has separate `Output` and `Run` methods to distinguish captured-output calls from streaming command execution.

## Data flow: status command
- `cmd/cm/main.go` calls `cli.Main` in `internal/cli/cli.go`.
- `cli.Main` creates `chezmoi.Client{Stdin, Stdout, Stderr}` and wraps it as `chezmoiService`.
- The root command or `status` subcommand calls `renderStatus` in `internal/cli/status.go`.
- `renderStatus` calls `Status(targets)` on the injected `statusService` interface.
- `chezmoiService.Status` in `internal/cli/service.go` delegates to `chezmoi.Client.Status`.
- `chezmoi.Client.Status` in `internal/chezmoi/client.go` runs `chezmoi status --path-style=absolute` and parses output with `chezmoi.ParseStatus`.
- `renderStatus` also calls `SourceStatus` to include source repository git changes.
- `chezmoiService.SourceStatus` runs `chezmoi source-path`, then runs `git status --porcelain=v1` in that directory through `process.Runner`.
- `parseSourceStatus` in `internal/cli/service.go` converts porcelain lines into `sourceEntry` values.
- `renderStatus` prints `clean` when both local status and source git status are empty.

## Data flow: diff command
- `cm diff [target...]` routes through `renderDiff` in `internal/cli/diff.go`.
- If no targets are provided, `renderDiff` calls `Status(nil)` and uses each `StatusEntry.Path` as a diff target.
- For each target, `renderDiff` calls `DiffOutput(target)` on the injected `diffService`.
- `chezmoiService.DiffOutput` constructs `syncdiff.Differ{Source: chezmoi.ContentLoader{Client: s.client}}` in `internal/cli/service.go`.
- `syncdiff.Differ.Diff` in `internal/syncdiff/diff.go` reads target content from `chezmoi cat <target>` through `ContentLoader.TargetContent`.
- `ContentLoader.LocalContent` in `internal/chezmoi/content.go` reads the local target path directly with a byte limit.
- `syncdiff.DiffBytes` rejects oversized files, reports binary differences, returns `no diff` for identical content, or emits a unified diff.
- `renderDiff` writes raw diff bytes to the command output writer.

## Data flow: interactive sync
- `cm sync [target...]` routes from `internal/cli/cli.go` to `ui.RunSyncTUI` in `internal/ui/tui.go`.
- `RunSyncTUI` first calls `ReviewService.Status(targets)` and exits early with `clean` when there are no entries.
- `newSyncTUIModel` in `internal/ui/model.go` stores the service, status entries, diff cache, help model, and home directory.
- `internal/ui/update.go` maps navigation keys, pending action keys, diff refresh, confirmation, and quit behavior.
- `internal/ui/diff.go` loads a selected target diff by calling `ReviewService.DiffOutput(target)` asynchronously through a Bubble Tea command.
- `internal/ui/confirm.go` rechecks fresh status for selected targets before execution to avoid applying actions to targets that became clean.
- Confirmed actions are sent to `ReviewService.Execute(actions)` and represented as `executeMsg` values.
- `chezmoiService.Execute` in `internal/cli/service.go` batches add actions, batches forced apply actions, and runs merge actions one target at a time.
- Apply actions in `Execute` call `chezmoi apply --force <targets...>`, while direct `cm apply` uses `chezmoi apply <targets...>` through `runTargets`.
- The TUI view in `internal/ui/view.go` renders file, diff, confirm, footer, and message panes from model state.

## Data flow: mutating wrappers and tools
- `cm add`, `cm apply`, and `cm merge` call `runTargets` in `internal/cli/service.go`, which constructs `[]string{command, targets...}` and calls `chezmoi.Client.Run`.
- `cm edit <target>` in `internal/cli/service.go` resolves the current user home directory and calls `chezmoi edit <home>/<target>`.
- `cm git` calls `chezmoi source-path`, validates the path is not empty, and runs `lazygit` in that directory with the CLI standard streams.
- Shell completion for `cm edit` calls `ManagedFiles`, which is implemented by `chezmoi.Client.ManagedFiles` in `internal/chezmoi/client.go`.
- `chezmoi.Client.ManagedFiles` parses command output through `ParseManagedFiles` in `internal/chezmoi/status.go`.

## Design patterns and testing seams
- Dependency injection is done with small interfaces in `internal/cli/service.go`: `statusService`, `diffService`, `targetCommandService`, `sourceGitService`, and `editService`.
- `commandServices` groups those interfaces so `newRootCommand` can be tested with fakes instead of real `chezmoi` or `git` processes.
- `reconcile.ReviewService` in `internal/reconcile/service.go` is a package-level boundary that keeps the TUI independent from concrete CLI service implementation.
- `syncdiff.ContentSource` in `internal/syncdiff/diff.go` decouples diffing logic from the chezmoi client and local filesystem details.
- `process.Runner` in `internal/process/runner.go` decouples command execution from command construction and gives tests a recording seam.
- Most functions accept `io.Reader` and `io.Writer` values, visible in `cmd/cm/main.go`, `internal/cli/cli.go`, and `internal/ui/tui.go`, which keeps CLI behavior testable without global standard streams.
- Errors are propagated up to Cobra, and `run` in `internal/cli/cli.go` prints command errors to stderr and returns `1`.
- The code favors composition over global state; the notable global-style value is `build.Version` in `internal/build/info.go` for release-time version injection.
- Interactive state is modeled as immutable-ish value updates on `syncTUIModel` methods in `internal/ui/model.go` and `internal/ui/update.go`.
- External command responsibilities are centralized, so future release hardening can focus on `internal/cli/service.go`, `internal/chezmoi/client.go`, and `internal/process/runner.go`.

## In-memory state shape
- `syncTUIModel` in `internal/ui/model.go` stores `entries []chezmoi.StatusEntry` as the ordered list of dirty targets presented to the user.
- `syncTUIModel.cursor` in `internal/ui/model.go` indexes the active entry and resets diff scroll when movement succeeds.
- `syncTUIModel.pending` in `internal/ui/model.go` maps target paths to selected `reconcile.ActionKind` values.
- `syncTUIModel.diffs` in `internal/ui/model.go` caches diff state per target path, including loaded lines, loading state, and error.
- `syncTUIModel.mode` in `internal/ui/model.go` separates review, confirmation, and execution phases.
- `syncTUIModel.focus` in `internal/ui/model.go` separates file-list navigation from diff-pane scrolling.
- `pendingActions` in `internal/ui/model.go` emits actions in status-entry order rather than Go map iteration order.
- `limitedBuffer` in `internal/chezmoi/client.go` keeps one byte beyond the requested limit so `syncdiff` can detect oversize content.

## State and side-effect model
- Read-only command paths are `cm`, `cm status`, `cm diff`, `cm version`, and `cm completion` as wired in `internal/cli/cli.go`.
- Mutating command paths are `cm add`, `cm apply`, `cm merge`, `cm edit`, and confirmed actions from `cm sync`.
- Tool-launching side effects are `cm git` in `internal/cli/service.go` and Bubble Tea terminal control in `internal/ui/tui.go`.
- `chezmoi.Client.Output` in `internal/chezmoi/client.go` captures stdout for read-style calls.
- `chezmoi.Client.Run` in `internal/chezmoi/client.go` streams through configured process I/O for command-style calls.
- `chezmoi.Client.OutputLimit` in `internal/chezmoi/client.go` uses `limitedBuffer` to bound target content capture for diffing.
- `syncTUIModel` in `internal/ui/model.go` stores pending actions separately from status entries until confirmation.
- `executeActions` in `internal/ui/confirm.go` validates pending targets against fresh status before calling `Execute`.

## Error and output behavior
- Cobra usage and error printing are silenced in `internal/cli/cli.go`; the local `run` function owns error output and nonzero exit conversion.
- Status output in `internal/cli/status.go` combines local divergence and source git divergence into one report.
- Diff output in `internal/cli/diff.go` writes bytes produced by `syncdiff` without additional per-line formatting.
- TUI diff display in `internal/ui/diff.go` applies presentation styling after loading the same diff bytes used by the `diff` command.
- Missing local files during diffing are tolerated in `internal/syncdiff/diff.go` when the error is `os.IsNotExist`.
- Oversized or binary file differences return explanatory text from `internal/syncdiff/diff.go` rather than failing the command.
- External process not-found errors are normalized in `internal/process/runner.go` as `<command> not found in PATH`.
- Source git status errors in `internal/cli/service.go` add source-directory context and trimmed stderr text when present.

## Release-readiness implications
- Command behavior can be reviewed mostly from `internal/cli/cli.go` plus the renderer and service files it delegates to.
- User-file mutation paths are concentrated in `internal/cli/service.go`, which makes it the most important file for release safety review.
- Terminal lifecycle risk is concentrated in `internal/ui/tui.go`, especially Bubble Tea program options and non-terminal rendering behavior.
- Diff correctness risk is concentrated in `internal/syncdiff/diff.go` and `internal/chezmoi/content.go`.
- External dependency availability risk is concentrated in `internal/process/runner.go` and the binary names chosen by `internal/chezmoi/client.go` and `internal/cli/service.go`.
- Version and release metadata behavior is concentrated in `internal/build/info.go` and `internal/cli/cli.go`.
- The Go package API is intentionally private under `internal/`, so v0.1.0 compatibility should focus on CLI commands, flags, output, completion, and side effects.
- The existing architecture has narrow seams for focused verification without requiring source edits outside the package under review.
