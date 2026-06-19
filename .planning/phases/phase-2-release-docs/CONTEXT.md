# Phase 2 Context: Release Metadata and Documentation Cleanup

## Goal

Make repository metadata and public documentation clean, current, and suitable for a v0.1.0 GitHub release.

## Constraints

- Use MIT license text; no external download is required.
- Delete `docs/superpowers/` as requested.
- Keep documentation small and current: README, usage, design, development, changelog, license, gitignore.
- Do not add CI in this phase; GitHub workflows belong to Phase 3.
- Do not add new runtime features.

## Decisions

- Write canonical MIT license text directly into `LICENSE` with the current year and project owner from module path.
- Add `.gitignore` for local binaries, test/coverage outputs, and release dist artifacts.
- Add `CHANGELOG.md` with an initial `v0.1.0` unreleased entry.
- Update `justfile` with version-aware install and release build helpers.
- Rewrite docs to include prerequisites, `cm edit`, executing-mode semantics, release commands, and current architecture.

## Open Questions

- Exact GitHub release date remains unknown until Phase 4.

## Verification Expectations

- `docs/superpowers/` no longer exists.
- Docs mention `cm edit`, prerequisites, release build, and MIT license.
- `just --list` should show install/build helpers if `just` is available.
- `VERSION=v0.1.0 just build-release` should build a versioned binary if `just` is available.
- `go test ./...` should still pass after metadata/docs changes.
