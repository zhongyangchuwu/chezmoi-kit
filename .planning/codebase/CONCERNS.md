# Release Concerns

Scope: release-risk map for v0.1.0 preparation, based on repository file inspection only; no project commands, tests, builds, linters, formatters, package-manager tasks, or runtime smoke checks were run.

## Stale or inconsistent documentation

- `README.md:61` has a typo: "No remplacer for chezmoi." This is user-visible release copy.
- `docs/superpowers/specs/2026-06-13-cm-mvp-design.md:24-35` still describes the first MVP status model with explanations/recommendations and sequential per-entry sync behavior; current code and docs implement a two-pane TUI with simplified `!` status output instead (`internal/cli/status.go:23-39`, `internal/ui/view.go:35-64`).
- `docs/superpowers/specs/2026-06-13-cm-mvp-design.md:36-42` says "Full-screen TUI" is out of scope, while current release behavior is a Bubble Tea full-screen-style sync interface (`internal/ui/tui.go:35-46`). Treat that spec as historical, not release truth.
- `docs/superpowers/specs/2026-06-13-cm-mvp-design.md:83-91` says `cm diff` calls `chezmoi diff`; current code uses internal rendered target vs local file diff (`internal/cli/diff.go:24-31`, `internal/cli/service.go:119-121`, `internal/syncdiff/diff.go:42-56`).
- `docs/development.md:18-20` lists targeted test packages but omits recently visible packages with tests: `internal/syncdiff`, `internal/build`, and `internal/cli` service behavior.
- `docs/development.md:134-138` gives a manual release build ldflags example, but there is no release packaging/checklist artifact beside it.
- `docs/usage.md:161` says `just install` generates zsh completion automatically; `justfile:6-15` confirms it writes only zsh completion, while release docs also advertise bash/fish/powershell scripts at `docs/usage.md:150-158` without install coverage.
- `docs/superpowers/plans/2026-06-13-cm-mvp.md:7-9` is also historical: it says standard library only and no TUI framework, but current code uses Cobra and Bubble Tea-related dependencies in `go.mod:6-11`.

## Missing release artifacts

- No `.github/` workflow was found, so CI status, release builds, and dependency checks are not represented in the repo.
- No `CHANGELOG` file was found; v0.1.0 release notes currently have no durable artifact.
- No `LICENSE` file was found; package consumers cannot verify redistribution terms from the repository.
- No goreleaser or equivalent release configuration was found; `justfile:3-15` only supports local installation.
- No issue/PR templates or security policy files were found under `.github/`; pre-release contribution and vulnerability intake are undefined.
- `internal/build/info.go:8` defaults `Version` to `dev`; without an explicit release build path, `cm version` can ship ambiguous version text.
- No binary distribution manifest or checksum/signing plan was found; release consumers would need to build from source or trust ad hoc local artifacts.

## Untested paths and behavior gaps

- `internal/process/runner.go` has no visible tests; missing `PATH` command errors, stderr wiring, stdin/stdout fallback behavior, and working-directory execution are release-critical because all external commands use this boundary.
- `internal/cli/service.go:83-88` parses source git status but tests cover only a happy path; error wrapping with `formatStderr` at `internal/cli/service.go:189-194` lacks visible coverage.
- `internal/cli/service.go:174-186` silently skips source git status lines shorter than four bytes; malformed or quoted porcelain paths are not surfaced and not visibly tested.
- `internal/cli/service.go:209-214` joins `$HOME` with a user-supplied edit target; visible tests cover `.zshrc` only, not absolute targets, `..`, or paths already returned by `chezmoi managed`.
- `internal/cli/cli.go:131-150` completion subcommands are only smoke-tested for bash output containing `cm`; shell-specific completion behavior is otherwise unverified.
- `internal/ui/update.go:46-58` confirm/executing key handling is lightly covered through one command-level happy path; quit during confirm, escape/back behavior, execute errors, and no-pending execute states need release attention.
- `internal/ui/tui.go:50-52` silently ignores a final model with an unexpected type; no test appears to assert this cannot hide a Bubble Tea lifecycle error.
- `internal/ui/diff.go:27-31` loads diffs synchronously inside a Bubble Tea command; stale diff responses for a previously selected target are stored by target and not reconciled with current cursor state.
- `internal/chezmoi/content.go:20-35` only reads one chunk of `limit+1` bytes; large file sentinel behavior is tested indirectly for the client fake, but local file read edge cases and read errors are not visible in tests.
- `internal/chezmoi/client.go:106-117` tests visible in `internal/chezmoi/client_test.go:53-64` cover truncation through a fake runner, not the real limited writer's short-write interaction with `exec.Cmd`.
- `internal/ui/view.go:203-214` truncates by byte length, not display width or rune boundaries; visible tests do not cover Unicode paths/diff lines even though the UI depends on terminal display-width libraries indirectly.
- `cmd/cm/main.go:9-10` only forwards `cli.Main` to `os.Exit`; command-level behavior is tested through `internal/cli`, but the process entrypoint itself has no visible smoke coverage.

## TUI lifecycle risks

- `internal/ui/tui.go:38` uses `tea.WithoutSignals()`, so interrupt/terminal signal handling is deliberately disabled; release UX for Ctrl-C, terminal resize, and process termination should be checked manually before v0.1.0.
- `internal/ui/tui.go:40-43` disables the renderer for non-terminal output, then prints the final full view at `internal/ui/tui.go:57-59`; command tests can pass while actual terminal rendering regressions go unnoticed.
- `internal/ui/view.go:41-53` forces minimum pane sizes even when terminal width/height are tiny; severe truncation or layout overflow is possible on small terminals.
- `internal/ui/view.go:107-108` renders an unloaded diff as `no diff`, which can be confused with a loaded identical diff; this is release UX risk because diffs load lazily.
- `internal/ui/view.go:153-164` footer advertises `q quit` during executing mode, and `internal/ui/keys.go:34-35` defines executing help, but `internal/ui/update.go:46-61` does not route keypresses to any executing-mode handler while a command runs.
- `internal/ui/confirm.go:10-40` executes all selected actions in one Bubble Tea command; long-running `chezmoi add/apply/merge` blocks UI updates until completion.
- `internal/ui/tui.go:72-74` returns no initial command; this differs from the two-pane plan at `docs/superpowers/plans/2026-06-16-sync-two-pane-tui.md:50-55`, which expected current diff loading on `Init` and cursor movement when cached behavior required it.
- `internal/ui/keys.go:59` binds both `q` and `ctrl+c` to quit, but `internal/ui/tui.go:38` disables Bubble Tea signal handling; Ctrl-C is therefore treated as a key only while the program is actively reading input.

## Subprocess and external-tool risks

- `internal/process/runner.go:25-38` captures stdout for command output but lets stderr default to process stderr unless overridden; some errors may be printed before they are wrapped at higher layers.
- `internal/process/runner.go:41-52` has no context, timeout, or cancellation path; hung `chezmoi`, `git`, merge tools, or `lazygit` will hang `cm`.
- `internal/cli/service.go:100-106` assumes `lazygit` is available; the missing-binary message comes from `ExecRunner`, but release docs list `chezmoi` and `just` requirements only (`docs/development.md:3-8`).
- `internal/cli/service.go:147-153` batches confirmed apply actions into one `chezmoi apply --force`; if one target fails, the error summary is coarse (`targetSummary`) and partial effects are not reported.
- `internal/cli/service.go:155-158` runs merge targets sequentially, but there is no retry/recovery state if the second merge fails after earlier actions already mutated files.
- `internal/chezmoi/client.go:45-60` uses a custom output limit via `Run` and `limitedBuffer`; commands that exit with `errLimitExceeded` are treated as successful truncation, so real command errors after writing the limit may be hidden if indistinguishable from the sentinel.
- `internal/cli/cli.go:72-93` allows zero arguments for direct mutating wrappers because all three use `cobra.ArbitraryArgs`; a bare `cm add`, `cm apply`, or `cm merge` forwards a whole-command mutation to chezmoi.
- `internal/cli/service.go:111-117` trims `chezmoi source-path`; whitespace-only output becomes the same empty-source case as no source path, which may hide unexpected chezmoi output.

## Security and dependency concerns

- `go.mod:3` declares `go 1.26`; release consumers and CI need a matching toolchain. `go.mod:6-11` also depends on major-version terminal UI packages from `charm.land` plus Cobra and `golang.org/x/term`.
- `go.mod:15-35` includes many terminal/ANSI indirect dependencies; without CI or dependency scanning artifacts, the release has no recorded vulnerability review.
- `internal/cli/service.go:209-214` accepts an edit target and constructs an absolute path under home without validating it is managed; an absolute or traversal-style target could be passed directly to `chezmoi edit`.
- `internal/process/runner.go:26` and `internal/process/runner.go:42` correctly avoid shell invocation via `exec.Command`, which reduces injection risk for target arguments.
- `internal/cli/service.go:83-86` includes the source directory and stderr in git status errors; useful for debugging but may disclose local filesystem paths in shared logs.
- `internal/syncdiff/diff.go:52-56` emits full diffs of managed configuration; docs should make clear that secrets in dotfiles can be displayed in terminal scrollback.
- `internal/chezmoi/content.go:12-17` diffs arbitrary target paths by passing them to `chezmoi cat` and `os.Open`; safety relies on the caller using managed targets from `chezmoi status`, except direct `cm diff <target>` accepts arbitrary args at `internal/cli/cli.go:54-60`.

## Performance and memory risks

- `internal/cli/diff.go:18-21` builds a full target list before diffing all dirty files; `internal/cli/diff.go:24-31` then writes each full diff sequentially, so `cm diff` can produce large terminal output for many targets.
- `internal/syncdiff/diff.go:11-13` caps each side at 1 MiB plus sentinel, which bounds common diff memory but can still allocate two large buffers per target plus diff output.
- `internal/chezmoi/content.go:30-35` allocates `limit+1` bytes for every local read; repeated diff refreshes in TUI can reallocate for the same large target.
- `internal/ui/diff.go:23-24` caches diff state by target for the session; repeated browsing of many files retains all loaded diff lines in memory.
- `internal/ui/view.go:119-123` allocates a visible slice and rendered strings every view render; acceptable for small terminal views, but large diff panes and rapid key repeats may be noticeable.
- `internal/cli/service.go:124-160` batches add/apply slices in memory before execution; low risk, but very large status sets are not streamed.
- `internal/ui/view.go:203-214` byte-slices long strings and can split multibyte UTF-8; this is both a rendering correctness and potential terminal garbage-output risk for non-ASCII paths or diff content.

## Maintainability gaps

- `internal/cli/service.go:66-230` centralizes status, source git, diff, execute, direct target commands, edit, and subprocess adaptation in one service type; release fixes in this file can easily affect unrelated commands.
- `internal/reconcile/service.go:5-52` defines action kinds and service contracts, but the original recommendation/description logic from `docs/superpowers/specs/2026-06-13-cm-mvp-design.md:119-132` is absent; current behavior relies on simplified UI labeling instead of an explicit domain interpretation layer.
- `internal/ui` state transitions span `model.go`, `update.go`, `confirm.go`, `diff.go`, and `tui.go`; lifecycle invariants such as `modeExecuting` key handling and diff load freshness are implicit and not centralized.
- `internal/cli/status.go:23-39` collapses every chezmoi status code to `! differs from chezmoi`; this is simple but makes future source/local/both distinctions harder to add without changing user-facing output.
- `internal/cli/service.go:163-168` uses coarse target summaries for multi-target failures; maintainers diagnosing release reports may need more detail about which target failed.
- `internal/ui/path.go:17-30` path display assumes home-relative shortening with OS separators; tests for Windows paths or non-home absolute paths are not visible even though terminal libraries include Windows-related dependencies.
- `internal/cli/cli.go:29-152` wires every command inline; command growth will make cross-command behavior such as IO, completions, and argument validation harder to audit before releases.
- `internal/chezmoi/client.go:15-21` exposes binary, stdio, dir, and runner knobs directly on the client struct; this is convenient for tests but makes construction invariants implicit.

## Highest-priority release checks

- Confirm whether missing `LICENSE`, `CHANGELOG`, CI, and release packaging are required before tagging v0.1.0; these are repository-level release blockers rather than code defects.
- Manually exercise `cm sync` in a real terminal because `internal/ui/tui.go:40-59` intentionally changes behavior for non-terminal outputs, which is what command tests observe.
- Verify external-command failure UX for missing `chezmoi`, missing `git`, missing `lazygit`, command non-zero exits, and a hung merge tool; all pass through `internal/process/runner.go:25-52`.
- Decide whether `cm git` should be documented as mutating-capable (`README.md:52`) or non-mutating from `cm`'s direct perspective; current copy says mutating because lazygit can mutate the source repository.
- Decide whether direct wrappers `cm add`, `cm apply`, and `cm merge` intentionally allow zero targets (`internal/cli/cli.go:72-93`), which forwards whole-repository operations to chezmoi.
