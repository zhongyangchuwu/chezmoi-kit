# Codebase Conventions

## Go style and naming

- The repository uses standard Go package scoping: lowercase package names (`cli`, `chezmoi`, `syncdiff`, `reconcile`, `ui`, `process`, `build`) and small files grouped by responsibility under `internal/`.
- `cmd/cm/main.go` is intentionally tiny: `main()` only calls `os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))`.
- Public exported types are reserved for cross-package contracts or adapters, for example `chezmoi.Client`, `chezmoi.StatusEntry`, `syncdiff.Differ`, `process.Runner`, `reconcile.Action`, and `build.Info`.
- Private helpers stay unexported and local to their package: examples include `newRootCommand`, `renderStatus`, `targetSummary`, `parseSourceStatus`, `formatStderr`, `visibleEntries`, `truncate`, and `padLines`.
- Receiver names are short but meaningful and conventional: `c Client`, `s chezmoiService`, `d Differ`, `m syncTUIModel`, `k ActionKind`, and pointer receivers only where mutation is needed such as `func (m *syncTUIModel) applyDiff(...)`.
- Enum-like values use a private integer type plus exported constants when used across packages: `reconcile.ActionKind` with `ActionAdd`, `ActionApply`, and `ActionMerge` in `internal/reconcile/service.go`.
- Small value structs carry domain words directly: `StatusEntry{Code, Path}`, `sourceEntry{Code, Path}`, `reconcile.Action{Target, Kind}`, and `process.IO{Stdin, Stdout, Stderr, Dir}`.
- Slices are preallocated when the final size is known or bounded, e.g. `make([]string, 0, 2+len(targets))` in `internal/chezmoi/client.go` and `make([]sourceEntry, 0, len(lines))` in `internal/cli/service.go`.

### Concrete style examples

```go
// `cmd/cm/main.go`
func main() {
	os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
```

```go
// `internal/reconcile/service.go`
type Action struct {
	Target string
	Kind   ActionKind
}
```

```go
// `internal/chezmoi/client.go`
args := make([]string, 0, 2+len(targets))
args = append(args, "status", "--path-style=absolute")
args = append(args, targets...)
```


## Error propagation

- Expected failures return `error`; there are no panic-based control paths in the inspected application code.
- Low-level boundaries usually return raw errors when they cannot add context, for example `chezmoi.Client.Output`, `chezmoi.Client.Run`, and `RunSyncTUI` returning `service.Status` errors unchanged.
- Service boundaries add action and target context with `%w`: `fmt.Errorf("git status %s: %w%s", sourceDir, err, formatStderr(stderr.Bytes()))`, `fmt.Errorf("lazygit %s: %w", sourceDir, err)`, `fmt.Errorf("add %s: %w", targetSummary(addTargets), err)`, and `fmt.Errorf("merge %s: %w", target, err)` in `internal/cli/service.go`.
- Parser errors include exact input context, for example `ParseStatus` in `internal/chezmoi/status.go` returns `fmt.Errorf("malformed chezmoi status line %d: %q", i+1, line)`.
- External command-not-found is normalized at the process boundary: `internal/process/runner.go` maps `exec.ErrNotFound` to `fmt.Errorf("%s not found in PATH", command)`.
- Benign domain cases are rendered as data instead of errors: `syncdiff.DiffBytes` returns `"binary file differs: ...\n"`, `"file too large to diff: ...\n"`, or `"no diff: ...\n"` as bytes.
- CLI execution handles command errors once in `internal/cli/cli.go`: `run` writes the error to `stderr` and returns exit code `1`; Cobra has `SilenceUsage: true` and `SilenceErrors: true`.

### Concrete error examples

```go
// `internal/cli/service.go`
return fmt.Errorf("git status %s: %w%s", sourceDir, err, formatStderr(stderr.Bytes()))
```

```go
// `internal/chezmoi/status.go`
return nil, fmt.Errorf("malformed chezmoi status line %d: %q", i+1, line)
```


## Interfaces, adapters, and fakes

- Interfaces live at package boundaries where tests or alternate implementations need seams, not as broad abstractions.
- `internal/process/runner.go` owns the subprocess seam: `Runner` exposes `Output(command string, args []string, io IO)` and `Run(command string, args []string, io IO)`; `ExecRunner` is the real implementation.
- `internal/chezmoi/client.go` aliases the process seam as `type RunnerIO = process.IO` and `type Runner = process.Runner`, keeping test fakes close to the chezmoi API without redefining a second process contract.
- `internal/syncdiff/diff.go` defines the content boundary as `ContentSource` with `TargetContent` and `LocalContent`; `internal/chezmoi/content.go` implements it via `ContentLoader`.
- `internal/reconcile/service.go` defines `ReviewService` as the TUI-facing contract: `Status`, `DiffOutput`, and `Execute`.
- `internal/cli/service.go` splits command needs into narrow interfaces (`statusService`, `diffService`, `targetCommandService`, `sourceGitService`, `editService`) and then groups them in `commandServices`.
- Test doubles are hand-written fakes in `_test.go`, named by behavior rather than framework: `fakeRunner` in `internal/chezmoi/client_test.go`, `recordingRunner` in `internal/cli/service_test.go`, `fakeService` in `internal/cli/cli_test.go`, `fakeContentSource` in `internal/syncdiff/diff_test.go`, and `fakeReviewService` in `internal/ui/sync_test.go`.
- Fakes copy mutable slices before storing or returning them, for example `append([]string(nil), targets...)`, `append([]chezmoi.StatusEntry(nil), entries...)`, and `bytes.Clone(...)`; this prevents tests from passing through accidental aliasing.

### Concrete seam examples

```go
// `internal/process/runner.go`
type Runner interface {
	Output(command string, args []string, io IO) ([]byte, error)
	Run(command string, args []string, io IO) error
}
```

```go
// `internal/reconcile/service.go`
type ReviewService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	Execute(actions []Action) error
}
```


## Test file conventions

- Tests live beside implementation files and use the same package name (`package cli`, `package ui`, `package chezmoi`, etc.), so they can exercise unexported helpers like `run`, `newSyncTUIModel`, `togglePending`, and `ParseManagedFiles` without external test packages.
- Test names are behavior sentences: `TestRunDefaultsToReadOnlyStatus`, `TestChezmoiServiceExecuteBatchesAddAndApplyBeforeSequentialMerges`, `TestDiffBytesSkipsBinaryContent`, and `TestExecutePendingPreflightDropsCleanTargets`.
- Assertions use the standard library only (`testing`, `reflect`, `strings`, `bytes`, `errors`); there is no testify or generated mock convention.
- Table tests are used where one behavior has many command variants, e.g. `TestRunMutatingWrappersForwardToChezmoi` in `internal/cli/cli_test.go`.
- Output tests assert meaningful substrings or exact protocol output, not renderer internals: status output checks for `"local:"`, `"run cm sync"`, and absence of raw `"MM /home/me/.zshrc"`.
- Command-forwarding tests record exact argv slices such as `[]string{"chezmoi", "apply", "--force", "/home/me/.zshrc"}` and `[]string{"git", "status", "--porcelain=v1"}`.
- TUI tests call model methods and Bubble Tea commands directly instead of spawning a terminal: examples include `model.updateReview(keyPress(tea.KeyTab))`, `cmd().(syncDiffMsg)`, and `model.executeActions(... )().(executeMsg)` in `internal/ui/sync_test.go`.

### Concrete test examples

```go
// `internal/cli/cli_test.go`
if !reflect.DeepEqual(service.commands, [][]string{{"diff-output", ".zshrc"}}) {
	t.Fatalf("commands = %#v", service.commands)
}
```

```go
// `internal/ui/sync_test.go`
msg := cmd().(syncDiffMsg)
if msg.target != "/home/me/.gitconfig" {
	t.Fatalf("diff target = %q", msg.target)
}
```


## Package boundaries

- `cmd/cm` owns only process startup.
- `internal/cli` owns Cobra command wiring, command rendering, service aggregation, and translating user commands to domain services.
- `internal/chezmoi` owns interaction with the `chezmoi` binary plus parsing of chezmoi output (`Status`, `ParseStatus`, `ManagedFiles`, `ContentLoader`).
- `internal/process` owns direct `os/exec` usage; higher packages call `Runner` instead of importing `os/exec`.
- `internal/syncdiff` owns internal diff generation from rendered chezmoi target content to local file content, including binary/large-file guards.
- `internal/reconcile` owns the shared sync action model (`ActionKind`, `Action`, `ReviewService`) and deliberately keeps UI presentation out.
- `internal/ui` owns Bubble Tea model state, key handling, diff display, confirmation, and final TUI rendering; it depends on `reconcile.ReviewService`, not on `chezmoi.Client` or git.
- `internal/build` owns runtime build metadata and formatting for `cm version`.

## Command rendering patterns

- Cobra commands are assembled in `newRootCommand` in `internal/cli/cli.go`; root `cm` and `cm status` both call `renderStatus(stdout, services.Status, args)`.
- Command constructors inject `stdin`, `stdout`, and `stderr`; commands do not write to global `os.Stdout` directly.
- Read-only commands return renderer errors directly: `renderStatus`, `renderDiff`, `version`, and `completion` all write to the injected writer and return write/generation errors.
- `renderStatus` in `internal/cli/status.go` renders `clean` when both local and source statuses are empty, otherwise it emits separate `local:` and `chezmoi:` sections.
- Plain CLI status uses `github.com/fatih/color` for headings and markers, while preserving simple text lines such as `"  local config differs from chezmoi source"` and `"  run cm sync"`.
- `renderDiff` in `internal/cli/diff.go` expands no-target invocations by calling `svc.Status(nil)` and diffing each dirty `entry.Path`; explicit targets bypass that status lookup.
- Completion is a normal Cobra subcommand: `completion [bash|zsh|fish|powershell]` switches to `root.GenBashCompletion`, `GenZshCompletion`, `GenFishCompletion`, or `GenPowerShellCompletion`, and returns `fmt.Errorf("unsupported shell %q", args[0])` for unknown shells.

### Concrete rendering examples

```go
// `internal/cli/cli.go`
RunE: func(cmd *cobra.Command, args []string) error {
	return renderStatus(stdout, services.Status, args)
},
```

```go
// `internal/cli/diff.go`
if _, err := w.Write(diff); err != nil {
	return err
}
```


## TUI patterns and testing style

- `RunSyncTUI` in `internal/ui/tui.go` preloads status, prints `clean` for no entries, then builds a Bubble Tea program with `tea.WithInput`, `tea.WithOutput`, and `tea.WithoutSignals`.
- Non-terminal output uses `tea.WithoutRenderer()` and then prints `final.viewString()`, which makes CLI-driven tests deterministic.
- TUI state is immutable-by-value for most transitions: methods such as `togglePending`, `clearPending`, `scrollDiff`, `toggleFocus`, `handleUp`, and `handleDown` return an updated `syncTUIModel`.
- Side effects are isolated behind Bubble Tea commands: `loadDiffCmd` returns `syncDiffMsg`, and `executeActions` returns `executeMsg`.
- Diffs load lazily: `startDiffLoad(false)` skips cached loaded diffs and only calls `DiffOutput` when focus moves to the diff pane or refresh is requested.
- The view code keeps layout helpers local (`rect`, `visibleEntry`, `clamp`, `truncate`, `padLines`, `splitLines`) and uses Lip Gloss styles declared in `internal/ui/keys.go`.
- Key bindings are centralized in `defaultSyncKeys`: `a` add, `p` apply, `m` merge, `s` skip, `d` refresh diff, `tab` focus, `enter` confirm, `y` execute, `esc` back, and `q`/`ctrl+c` quit.
- Before executing, `executeActions` rechecks selected targets with `m.service.Status(targets)` and drops actions whose targets are already clean; this behavior is covered by `TestExecutePendingPreflightDropsCleanTargets`.

### Concrete TUI examples

```go
// `internal/ui/diff.go`
return syncDiffMsg{target: target, diff: string(out), err: err}
```

```go
// `internal/ui/tui.go`
options := []tea.ProgramOption{
	tea.WithInput(input),
	tea.WithOutput(output),
	tea.WithoutSignals(),
}
```


## Release and development command conventions

- Module identity and Go version are declared in `go.mod`: `module github.com/zhongyangchuwu/cm` and `go 1.26`.
- Runtime dependencies are intentionally focused on CLI/TUI concerns: Bubble Tea/Bubbles/Lip Gloss, Cobra, `fatih/color`, `golang.org/x/term`, and `rogpeppe/go-internal/diff`.
- `justfile` contains the local install workflow only: `go install ./cmd/cm`, generate zsh completion with `go run ./cmd/cm completion zsh > "${HOME}/.zfunc/_cm"`, then print the installed binary and completion paths.
- `docs/development.md` documents project-wide tests and targeted packages (`go test ./internal/chezmoi`, `go test ./internal/ui`, `go test ./cmd/cm`) but this mapping task did not run any commands.
- Local run examples are documented as `go run ./cmd/cm status`, `go run ./cmd/cm sync`, and `go run ./cmd/cm diff ~/.zshrc` in `docs/development.md`.
- Release versioning is built around `internal/build.Version`; `docs/development.md` documents `go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/build.Version=v0.1.0"`.
- `internal/build/info.go` prefers module build info when available and formats stable output through `Info.FormatShort("cm")` and `Info.FormatDetailed("cm")`; version output includes `cm`, `commit`, `built`, `dirty`, and `go` fields.

### Concrete release and dev examples

```just
# `justfile`
install:
    go install ./cmd/cm
    go run ./cmd/cm completion zsh > "${HOME}/.zfunc/_cm"
```

```bash
# `docs/development.md`
go test ./...
go test ./internal/chezmoi
go test ./internal/ui
go test ./cmd/cm
```

```bash
# `docs/development.md`
go run ./cmd/cm status
go run ./cmd/cm sync
go run ./cmd/cm diff ~/.zshrc
```

```bash
# `docs/development.md`
go build -ldflags "-X github.com/zhongyangchuwu/cm/internal/build.Version=v0.1.0"
```

```go
// `internal/build/info.go`
var Version = "dev"

func (i Info) FormatDetailed(name string) string {
	return fmt.Sprintf("%s: %s\ncommit: %s\nbuilt: %s\ndirty: %s\ngo: %s\n", name, i.Version, fallback(i.Commit, "unknown"), fallback(i.Time, "unknown"), fallback(i.Modified, "unknown"), fallback(i.GoVersion, "unknown"))
}
```
