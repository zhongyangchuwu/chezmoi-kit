# Plan: 5-package-architecture

## Objective

Normalize package and file boundaries so later app-service and report-rendering changes have clear homes.

## Scope

In scope:

- Rename `internal/ui` package and imports to `internal/tui`.
- Move version metadata from `internal/build` into `internal/app` and update release/build ldflags references.
- Remove `internal/testutil` by moving tiny test helpers to local tests or eliminating them through ownership-specific tests.
- Split `internal/tui/diff.go` into state/loading and presentation files.
- Rename CLI files to clarify that CLI is command wiring only after app extraction begins.
- Update docs and planning/codebase references that mention old package names.

Out of scope:

- Introducing semantic report documents; Phase 7 owns that.
- Moving all `chezmoiService` behavior into `app`; Phase 6 owns that.
- Replacing current CLI output format.
- Adding community test libraries.

## Tasks

1. Rename `internal/ui` directory to `internal/tui`; update package declarations, imports, docs, and tests.
2. Move `internal/build/info.go` to `internal/app/version.go`; update imports and ldflags references in release/build config and docs.
3. Delete `internal/testutil`; replace `TerminalCommand` and debug path helpers with local test helpers only where still needed.
4. Split TUI diff responsibilities:
   - `diff_state.go`: `diffState`, load command, apply, scroll, line splitting.
   - `diff_view.go`: diff line classification-to-style rendering for the TUI.
5. Rename CLI entry file to `root.go` if still appropriate; keep Cobra command construction and completion in CLI.
6. Update `.planning/codebase/*`, `docs/design.md`, and `docs/development.md` to reflect package names.
7. Run targeted and full verification.

## Acceptance Criteria

- Package names express architecture layers: `cli`, `tui`, `app`, `chezmoi`, `diff`, `process`.
- `internal/build` and `internal/testutil` no longer exist.
- TUI diff state/loading and diff rendering are in separate files.
- Build/version ldflags point at `internal/app.Version`.
- Tests no longer rely on shared test helper packages.
- Documentation and codebase maps no longer mention stale `internal/ui`, `internal/build`, or `internal/testutil` package boundaries.

## Verification

Run after implementation:

```bash
go test ./internal/app ./internal/cli ./internal/tui ./internal/chezmoi ./internal/diff ./internal/process
go test ./...
go vet ./...
go mod tidy -diff
```

Search checks:

```text
No `internal/ui` imports.
No `internal/build` imports or ldflags.
No `internal/testutil` imports.
```

## Risks

- Renaming packages can leave stale docs, workflow config, or ldflags references.
- Moving version metadata changes release ldflags; release build smoke must confirm injected versions still work.
- Removing shared test helpers may briefly duplicate tiny local fakes; that is acceptable until Phase 6 moves ownership-specific tests into app.
