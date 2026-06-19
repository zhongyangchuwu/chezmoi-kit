# Phase 4 Context: Release Verification

## Goal

Prove the repository is ready to tag `v0.1.0` and publish through GoReleaser.

## Constraints

- Do not create or push the tag until local verification, remote CI, and manual smoke requirements are satisfied or explicitly documented as environment-limited.
- Keep release notes honest: update `CHANGELOG.md` date only when preparing the actual release tag.
- Do not add new features in this phase.
- Do not bypass failing checks.

## Decisions

- Run the documented local release gate.
- Run non-mutating command smoke checks in this environment.
- Document environment-limited smoke gaps for commands that require a real chezmoi-managed setup, lazygit, or an interactive terminal.
- Commit release verification artifacts before any tag is created.

## Open Questions

- Remote GitHub CI can only be confirmed after pushing the branch/commit.
- Manual `cm sync` Ctrl+C terminal cleanup requires a real interactive terminal.

## Verification Expectations

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go mod tidy -diff`
- `go mod verify`
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`
- `go run github.com/goreleaser/goreleaser/v2@latest check`
- `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean`
- `./dist` artifact inspection
- CLI smoke: `version`, `--help`, `completion zsh`, and any safe read-only command available in this environment
