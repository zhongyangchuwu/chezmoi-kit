# Plan: 4-release-gate

## Objective

Run final release verification and prepare the repository for a v0.1.0 tag.

## Scope

In scope:

- Local release gate commands.
- GoReleaser snapshot validation and artifact inspection.
- Safe CLI smoke checks.
- Changelog release-date preparation if tag is ready.
- Phase verification records and final planning state updates.
- Commit release verification artifacts.

Out of scope:

- Adding features or changing release scope.
- Publishing the GitHub release before CI passes remotely.
- Creating the tag without user approval after verification.

## Tasks

1. Run local release gate.
2. Inspect generated GoReleaser snapshot artifacts.
3. Run available non-mutating CLI smoke checks.
4. Record environment-limited smoke checks that require real terminal or user chezmoi setup.
5. Update changelog release date if proceeding to tag preparation.
6. Write `SUMMARY.md` and `VERIFICATION.md`.
7. Update root planning artifacts and commit phase artifacts.

## Acceptance Criteria

- Local release gate passes.
- GoReleaser snapshot build succeeds and produces expected archive/checksum artifacts.
- CLI smoke checks pass or limitations are documented.
- Remote CI requirement is explicitly marked pending until branch is pushed and observed.
- Release notes are ready for v0.1.0 or the remaining blocker is named.

## Verification

```bash
go test ./...
go test -race ./...
go vet ./...
go mod tidy -diff
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go run github.com/goreleaser/goreleaser/v2@latest check
go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean
./dist/cm_*/cm version
```

## Risks

- Commands that need a configured chezmoi repository or lazygit may fail in the harness for environment reasons; do not claim them as behavior-verified unless actually observed.
- Remote CI cannot be proven locally.
