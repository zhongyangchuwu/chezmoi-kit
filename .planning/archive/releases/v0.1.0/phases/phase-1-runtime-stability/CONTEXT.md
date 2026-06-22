# Phase 1 Context: Runtime Stability

## Goal

Make the existing `cm` runtime safe and tidy enough to serve as the base for v0.1.0 release documentation and CI.

## Constraints

- Keep this phase focused on code/runtime fixes only.
- Do not delete `docs/superpowers/` in this phase; documentation cleanup belongs to Phase 2.
- Do not add new user-facing features beyond required fixes.
- Do not implement broad subprocess cancellation for v0.1.0 unless needed to preserve correctness.
- Prefer conservative executing-mode semantics: no fake quit/cancel once confirmed execution starts.
- Preserve existing command surface and user workflows.

## Decisions

- Remove production use of `tea.WithoutSignals()` so Bubble Tea can clean up terminal state on Ctrl+C.
- During `modeExecuting`, ignore keypresses and do not advertise quit; execution completes or returns an error.
- Make source git status parsing strict enough to surface malformed non-empty porcelain output.
- Add command wiring coverage for `cm edit <target>` by extending the CLI test fake.
- Fix `go.mod`/`go.sum` with `go mod tidy` instead of hand-maintaining inconsistent module metadata.

## Open Questions

- Whether Ctrl+C terminal cleanup can be manually smoke-tested in this non-interactive harness. If not, record it as an environment-limited manual check for Phase 4.

## Verification Expectations

- `go test ./internal/ui`
- `go test ./internal/cli`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go mod tidy -diff`
- `go mod verify`
- Build `./cmd/cm` and inspect `--help`/`version` if runtime wiring changed.
