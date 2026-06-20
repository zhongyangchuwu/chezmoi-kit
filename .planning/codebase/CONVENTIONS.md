# Codebase Conventions

## Go style and naming

- The repository uses standard Go package scoping: lowercase package names (`app`, `cli`, `chezmoi`, `diff`, `tui`, `process`) and small files grouped by responsibility under `internal/`.
- `cmd/cm/main.go` is intentionally tiny: `main()` only calls `os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))`.
- Public exported types are reserved for cross-package contracts or adapters, for example `app.Info`, `app.Services`, `app.Action`, `app.SyncService`, `chezmoi.Client`, `chezmoi.StatusEntry`, `diff.Differ`, and `process.Runner`.
- Private helpers stay unexported and local to their package: examples include `newRootCommand`, `renderStatus`, `targetSummary`, `formatStderr`, `visibleEntries`, `truncate`, and `padLines`.
- Receiver names are short but meaningful and conventional: `c Client`, `s service`, `d Differ`, `m syncTUIModel`, `k ActionKind`, and pointer receivers only where mutation is needed such as `func (m *syncTUIModel) applyDiff(...)`.
- Enum-like values use a private integer type plus exported constants when used across packages: `app.ActionKind` with `ActionAdd`, `ActionApply`, and `ActionMerge` in `internal/app/services.go`.
- Small value structs carry domain words directly: `StatusEntry{Code, Path}`, `app.Action{Target, Kind}`, and `process.IO{Stdin, Stdout, Stderr, Dir}`.
- Slices are preallocated when the final size is known or bounded, e.g. `make([]string, 0, 2+len(targets))` in `internal/chezmoi/client.go`.

## Error propagation

- Expected failures return `error`; there are no panic-based control paths in the inspected application code.
- Low-level boundaries usually return raw errors when they cannot add context, for example `chezmoi.Client.Output`, `chezmoi.Client.Run`, and `RunSyncTUI` returning `service.Status` errors unchanged.
- Service boundaries add action and target context with `%w`, such as source git and lazygit errors in `internal/app/services.go`.
- Parser errors include exact input context, for example `ParseStatus` in `internal/chezmoi/status.go` returns `fmt.Errorf("malformed chezmoi status line %d: %q", i+1, line)`.
- External command-not-found is normalized at the process boundary: `internal/process/runner.go` maps `exec.ErrNotFound` to `fmt.Errorf("%s not found in PATH", command)`.
- Benign domain cases are rendered as data instead of errors: `diff.DiffBytes` returns `"binary file differs: ...\n"`, `"file too large to diff: ...\n"`, or `"no diff: ...\n"` as bytes.
- CLI execution handles command errors once in `internal/cli/root.go`: `run` writes the error to `stderr` and returns exit code `1`; Cobra has `SilenceUsage: true` and `SilenceErrors: true`.

## Interfaces, adapters, and fakes

- Interfaces live at package boundaries where tests or alternate implementations need seams, not as broad abstractions.
- `internal/process/runner.go` owns the subprocess seam: `Runner` exposes `Output(command string, args []string, io IO)` and `Run(command string, args []string, io IO)`; `ExecRunner` is the real implementation.
- `internal/chezmoi/client.go` aliases the process seam as `type RunnerIO = process.IO` and `type Runner = process.Runner`, keeping test fakes close to the chezmoi API without redefining a second process contract.
- `internal/diff/diff.go` defines the content boundary as `ContentSource` with `TargetContent` and `LocalContent`; `internal/chezmoi/content.go` implements it via `ContentLoader`.
- `internal/app/services.go` defines the app service graph and boundary contracts: `StatusService`, `DiffService`, `TargetCommandService`, `SourceGitService`, `EditService`, and `SyncService`.
- `internal/report` defines semantic output contracts, renderers, color policy, and shared diff line classification.
- `internal/cli/service.go` aliases `app.Services` so command wiring has one service graph type without owning concrete orchestration.
- Fakes copy mutable slices before storing or returning them, for example `append([]string(nil), targets...)`, `append([]chezmoi.StatusEntry(nil), entries...)`, and `bytes.Clone(...)`; this prevents tests from passing through accidental aliasing.

## Test file conventions

- Tests live beside implementation files and use the same package name (`package cli`, `package tui`, `package chezmoi`, etc.), so they can exercise unexported helpers like `run`, `newSyncTUIModel`, `togglePending`, and `ParseManagedFiles` without external test packages.
- Test names are behavior sentences: `TestRunDefaultsToReadOnlyStatus`, `TestDiffBytesSkipsBinaryContent`, and `TestExecuteCurrentActionDropsCleanTarget`.
- Assertions use the standard library only (`testing`, `reflect`, `strings`, `bytes`, `errors`); there is no testify or generated mock convention.
- Table tests are used where one behavior has many command variants, e.g. `TestRunMutatingWrappersForwardToChezmoi` in `internal/cli/cli_test.go`.
- Output tests assert meaningful substrings or exact protocol output, not renderer internals: status output checks for `"local:"`, `"run cm sync"`, and absence of raw `"MM /home/me/.zshrc"`.
- Command-forwarding tests record exact argv slices such as `[]string{"chezmoi", "apply", "--force", "/home/me/.zshrc"}` and `[]string{"git", "status", "--porcelain=v1"}`.
- TUI tests call model methods and Bubble Tea commands directly instead of spawning a terminal: examples include `model.updateReview(keyPress(tea.KeyTab))` and `cmd().(syncDiffMsg)` in `internal/tui/sync_test.go`.

## Package boundaries

- `cmd/cm` owns only process startup.
- `internal/app` owns process options, runtime version metadata, service composition, command use cases, semantic report construction, and sync action contracts.
- `internal/cli` owns Cobra command wiring, report rendering to streams, completion, and exit behavior.
- `internal/chezmoi` owns interaction with the `chezmoi` binary plus parsing of chezmoi output (`Status`, `ParseStatus`, `ManagedFiles`, `ContentLoader`).
- `internal/process` owns direct `os/exec` usage; higher packages call `Runner` instead of importing `os/exec`.
- `internal/diff` owns internal diff generation from rendered chezmoi target content to local file content, including binary/large-file guards.
- `internal/tui` owns Bubble Tea model state, key handling, diff display, confirmation, and final TUI rendering.

## Command rendering patterns

- Cobra commands are assembled in `newRootCommand` in `internal/cli/root.go`; root `cm` and `cm status` both call `renderStatus(stdout, services.Status, args)`.
- Command constructors inject `stdin`, `stdout`, and `stderr`; commands do not write to global `os.Stdout` directly.
- Read-only commands return renderer errors directly: `renderStatus`, `renderDiff`, `version`, and `completion` all write to the injected writer and return write/generation errors.
- `internal/app` builds semantic status, diff, and version `report.Document` values; it does not hand-build final ANSI/plain strings.
- `internal/cli/render.go` renders reports with auto ANSI color for TTY stdout and plain text otherwise; non-empty `NO_COLOR` disables ANSI.
- `internal/report/diff.go` classifies diff lines once for CLI report rendering and TUI diff styling.
- Completion is a normal Cobra subcommand: `completion [bash|zsh|fish|powershell]` switches to Cobra completion generators and returns `fmt.Errorf("unsupported shell %q", args[0])` for unknown shells.

## TUI patterns and testing style

- `RunSyncTUI` in `internal/tui/tui.go` preloads status, prints `clean` for no entries, then builds a Bubble Tea program with `tea.WithInput`, `tea.WithOutput`, and terminal-aware renderer options.
- Non-terminal output uses `tea.WithoutRenderer()` and then prints `final.viewString()`, which makes CLI-driven tests deterministic.
- TUI state is immutable-by-value for most transitions: methods such as `togglePending`, `clearPending`, `scrollDiff`, `toggleFocus`, `handleUp`, and `handleDown` return an updated `syncTUIModel`.
- Side effects are isolated behind Bubble Tea commands: `loadDiffCmd` returns `syncDiffMsg`, and execution commands return `executeMsg`.
- Diffs load by default on entry and when moving between files; `startDiffLoad(false)` skips cached loaded diffs.
- Diff state/loading lives in `internal/tui/diff_state.go`; diff line rendering lives in `internal/tui/diff_view.go`.
- Key bindings are centralized in `defaultSyncKeys`: `a` add, `p` apply, `m` merge, `s` skip, `d` refresh diff, `tab` focus, `enter` confirm, `y` execute, `esc` back, and `q`/`ctrl+c` quit before execution.
- Before executing, command logic rechecks selected targets with `m.service.Status(targets)` and drops actions whose targets are already clean.

## Release and development command conventions

- Module identity and Go version are declared in `go.mod`: `module github.com/zhongyangchuwu/cm` and `go 1.26`.
- Runtime dependencies are intentionally focused on CLI/TUI concerns: Bubble Tea/Bubbles/Lip Gloss, Cobra, `golang.org/x/term`, and `rogpeppe/go-internal/diff`.
- `justfile` contains local install and build-release workflows and injects release versions with `internal/app.Version`.
- `docs/development.md` documents project-wide tests and targeted packages including `internal/app`, `internal/cli`, `internal/chezmoi`, `internal/diff`, and `internal/tui`.
- Local run examples are documented as `go run ./cmd/cm status`, `go run ./cmd/cm sync`, and `go run ./cmd/cm diff ~/.zshrc` in `docs/development.md`.
- Release versioning is built around `internal/app.Version`; `docs/development.md` documents `go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/app.Version=v0.1.0"`.
