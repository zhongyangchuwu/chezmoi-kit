# Summary: Phase 2 Release Metadata and Documentation Cleanup

## Completed Changes

- Deleted stale `docs/superpowers/` planning/spec artifacts from public docs.
- Added `LICENSE` with MIT license text.
- Added `.gitignore` for local binary, `dist/`, test binaries, and coverage output.
- Added `CHANGELOG.md` with initial `v0.1.0` unreleased notes.
- Updated `justfile` with version-aware `install` and `build-release` recipes.
- Rewrote `README.md` with requirements, install, command table, sync safety model, and license.
- Rewrote `docs/usage.md` with prerequisites, status/diff/sync/direct command workflows, `cm edit`, completion, and version behavior.
- Rewrote `docs/design.md` to match current package boundaries, sync execution semantics, and strict git status parsing.
- Rewrote `docs/development.md` with test commands, release gate, versioned build flow, tag flow, and updated package layout.
- Fixed `internal/build.Current` so explicit `-ldflags ...Version=v0.1.0` overrides Go module VCS version metadata.
- Added a test proving explicit build version precedence.

## Files Changed

- `README.md`
- `docs/usage.md`
- `docs/design.md`
- `docs/development.md`
- `docs/superpowers/` deleted
- `LICENSE`
- `.gitignore`
- `CHANGELOG.md`
- `justfile`
- `internal/build/info.go`
- `internal/build/info_test.go`

## Deviations

- Phase 2 included a small build metadata code fix because release recipe verification showed `VERSION=v0.1.0 just build-release` still printed the Go module pseudo-version. This was necessary to satisfy the phase acceptance criterion for versioned release builds.
- GitHub Actions were not added; Phase 3 owns CI and release automation.

## Evidence

- `docs/superpowers/` is absent from `find` output.
- `LICENSE`, `.gitignore`, and `CHANGELOG.md` exist.
- `just --list` shows `build-release` and `install`.
- `go test ./...` passed.
- `go mod tidy -diff` passed with no output.
- Go workspace diagnostics reported no issues.
- `rm -rf dist && VERSION=v0.1.0 just build-release && ./dist/cm version` passed and printed `cm: v0.1.0`.

## Unresolved Risks

- Phase 3 still needs GitHub CI and tag release workflows.
- Phase 4 still needs manual terminal smoke for `cm sync` Ctrl+C cleanup.
- `CHANGELOG.md` date remains `Unreleased` until the final tag date is known.
