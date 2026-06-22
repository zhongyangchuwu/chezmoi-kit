# Plan: 2-release-docs

## Objective

Clean public release artifacts and documentation after Phase 1 runtime stability fixes.

## Scope

In scope:

- Delete `docs/superpowers/`.
- Add `LICENSE`, `.gitignore`, and `CHANGELOG.md`.
- Add version-aware install/build recipes in `justfile`.
- Update `README.md`, `docs/usage.md`, `docs/design.md`, and `docs/development.md`.
- Update planning state after verification.

Out of scope:

- GitHub Actions CI/release workflows.
- Runtime code changes beyond build helper needs.
- GoReleaser.
- New command features.

## Tasks

1. Delete stale `docs/superpowers/` directory.
2. Add MIT `LICENSE`, `.gitignore`, and initial `CHANGELOG.md`.
3. Update `justfile` with `install` using `VERSION` and `build-release` writing `dist/cm`.
4. Rewrite README with requirements, install, quick start, complete command table, safety model, and license.
5. Rewrite usage docs with prerequisites, sync behavior, direct commands including edit, completion, and version behavior.
6. Rewrite design docs with current package boundaries and execution semantics.
7. Rewrite development docs with test/release commands and phase-relevant project layout.
8. Verify docs/metadata and write summary/verification artifacts.

## Acceptance Criteria

- `docs/superpowers/` is removed.
- `LICENSE` contains MIT license text.
- `.gitignore` covers local binary, test/coverage outputs, and `dist/`.
- `CHANGELOG.md` contains a v0.1.0 entry.
- README and usage docs document all v0.1.0 commands including `cm edit`.
- Development docs contain release gate and tag flow.
- `justfile` supports versioned release builds.
- Go tests still pass.

## Verification

```bash
go test ./...
go mod tidy -diff
just --list
VERSION=v0.1.0 just build-release
./dist/cm version
```

If `just` is unavailable, verify equivalent `go build` command manually.

## Risks

- License copyright holder may need later adjustment if the repository owner wants a different legal name.
- `VERSION=v0.1.0 just install` writes to the user's Go bin and zsh completion path; avoid running install in automated verification unless explicitly needed.
