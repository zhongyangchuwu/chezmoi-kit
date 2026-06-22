# Summary: Phase 5 Package Architecture Cleanup

## Completed Changes

- Renamed the interactive terminal package from `internal/ui` to `internal/tui` and updated Go imports, package declarations, docs, and codebase maps.
- Merged version metadata from `internal/build` into `internal/app/version.go` and updated `justfile` and GoReleaser ldflags to target `internal/app.Version`.
- Removed `internal/testutil`; CLI and TUI tests now keep tiny terminal-command and debug-log helpers local to their owning packages.
- Split TUI diff responsibilities into `internal/tui/diff_state.go` for cache/loading/scrolling/splitting and `internal/tui/diff_view.go` for diff line styling.
- Renamed CLI command wiring files to `root.go` and `diff_cmd.go` to make command ownership clearer.
- Refreshed docs and `.planning/codebase/*` maps to match the new package names and file boundaries.

## Files Changed

- `.goreleaser.yaml`
- `justfile`
- `docs/design.md`
- `docs/development.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/CONCERNS.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/MAP.md`
- `.planning/codebase/STACK.md`
- `.planning/codebase/STRUCTURE.md`
- `internal/app/version.go`
- `internal/app/version_test.go`
- `internal/cli/root.go`
- `internal/cli/diff_cmd.go`
- `internal/cli/cli_test.go`
- `internal/cli/service.go`
- `internal/cli/service_test.go`
- `internal/tui/*`
- Removed `internal/build/*`
- Removed `internal/ui/*`
- Removed `internal/testutil/helpers.go`

## Deviations

- `internal/tui` still owns the sync action contract for this phase; Phase 6 is planned to move that contract into `internal/app` with the service graph.
- The codebase maps were refreshed directly rather than preserving historical stale package descriptions because Phase 5 explicitly changes those package boundaries.

## Evidence

- `go test ./internal/app ./internal/cli ./internal/tui ./internal/chezmoi ./internal/diff ./internal/process` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go mod tidy -diff` passed with no output.
- Search found no code/build references to `github.com/zhongyangchuwu/cm/internal/ui`, `internal/build`, or `internal/testutil`.
- `VERSION=v9.9.9 just build-release && ./dist/cm version` printed `cm: v9.9.9`, proving the new `internal/app.Version` ldflag path works.

## Unresolved Risks

- Phase 6 still needs to move application service ownership and the sync contract into `internal/app`.
- Phase 7 still needs to separate semantic output content from renderers.
