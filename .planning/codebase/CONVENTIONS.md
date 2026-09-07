# Conventions

**Mapped:** 2026-09-07

## Go Style

- Use standard Go package layout under `cmd/` and `internal/`.
- Keep `main` thin; behavior belongs in app/adapter packages.
- Run `gofmt`; avoid project-specific formatting layers.
- Prefer narrow interfaces at adapter boundaries and concrete unexported service implementations.
- Return errors with operation/target context where recovery is not local.
- Reserve terminal ownership for commands that genuinely require interaction.

## Architecture

- CLI/TUI are adapters, not orchestration owners.
- `internal/app` owns use cases and user-facing reconciliation contracts.
- `internal/chezmoi` owns command construction and strict chezmoi parsing.
- `internal/process` is the sole runtime `os/exec` boundary.
- `internal/report` owns semantic output and rendering.
- One authoritative diff path: forced chezmoi builtin reverse diff.
- Clean cutovers remove obsolete packages, dependencies, tests, and docs.

## Reconciliation

- Interpret destination/target work from chezmoi's second status column.
- Separate scripts before ordinary reconciliation.
- Load target metadata lazily for the selected item.
- Gate actions by reviewed type/template state.
- Store the review fingerprint in pending actions.
- Recompute review before mutation and status/review after success.
- Command exit success does not imply a clean postcondition.
- Keep scripts outside ordinary sync execution.

## Paths and Parsers

- Request absolute status paths from chezmoi.
- Use NUL-delimited path output when available.
- Reject malformed non-empty status/git lines rather than skipping them.
- Resolve relative edit targets from `chezmoi target-path`.
- Do not infer chezmoi source attributes from filename prefixes/suffixes when a command/filter supplies the fact.
- Build workspace inventory from chezmoi path/type/status/ignore/unmanaged outputs; filesystem traversal may enumerate scoped child candidates only and never replaces chezmoi classification.
- Treat `--skip-secrets` omissions as uninspected, not clean. Keep preview content bounded and require explicit reveal for sensitive target/source/diff content.
- Project tree/flat/filter/search views from absolute target identity, retaining the same selected target where visible and otherwise using deterministic index fallback.
- Workspace color augments rather than replaces `C/D/U/I/R/?`, target-type letters, and `[T]/[E]` badges; `NO_COLOR=1` retains complete text semantics.
- TUI overlays intercept their own keys before normal workspace dispatch; compact footers retain a help and quit route at narrow widths.

## Errors and Output

- CLI writes command errors once and exits 1 without duplicate Cobra usage.
- Benign clean/no-file states render as semantic data.
- Oversized authoritative diff/metadata is an error; do not permit mutation without complete review.
- Successful buffered stdout/stderr is preserved for the TUI.
- Debug timing logs include targets, action/type flags, durations, and errors, never rendered contents or subprocess output.
- Public reports use `report.Document`; Markdown is a renderer, not the internal model.

## Tests

- Tests live beside implementation and normally use the same package.
- Handwritten fakes copy mutable slices before retaining/returning them.
- Assert user-observable behavior, action/state transitions, parser rejection, and exact external argv contracts where those commands are the adapter contract.
- Avoid tests that pin private helper names or source text.
- Real chezmoi integration tests use temporary source/destination/cache/config/state and a wrapper binary; they never touch user state.
- Integration tests skip only when chezmoi is absent or the POSIX wrapper cannot run.
- Async TUI tests use condition gates rather than arbitrary sleeps.

## Verification

Focused packages:

```text
go test ./internal/app
go test ./internal/chezmoi
go test ./internal/cli
go test ./internal/tui
```

Full gates:

```text
go test ./...
go test -race ./...
go vet ./...
go mod tidy -diff
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Manual/isolated TUI smoke remains required for review, confirmation, stale-state deferral, terminal restoration, and postflight messages.
