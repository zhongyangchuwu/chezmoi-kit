# Phase 3 Context: CI and Release Automation

## Goal

Add GitHub CI and GoReleaser-based tag release automation for v0.1.0.

## Constraints

- Use GoReleaser as requested.
- Keep workflows minimal and auditable.
- Do not add signing, Homebrew, Docker, package managers, or GoReleaser Pro features in v0.1.0.
- CI must match the documented release gate.
- Release workflow must use `goreleaser/goreleaser-action@v7`, checkout with `fetch-depth: 0`, and `GITHUB_TOKEN` with `contents: write`.

## Decisions

- Add `.goreleaser.yaml` using GoReleaser v2 config.
- Build `./cmd/cm` with binary name `cm` for Linux, macOS, and Windows on amd64/arm64.
- Use `CGO_ENABLED=0` and `-trimpath` for portable static-ish Go binaries.
- Inject `internal/build.Version={{ .Version }}` through GoReleaser ldflags.
- Use GoReleaser archives and checksum generation.
- CI workflow runs tests, race tests, vet, tidy diff, module verify, govulncheck, GoReleaser check, and snapshot build.
- Release workflow runs only on `v*` tags and publishes GitHub release artifacts through GoReleaser.

## Open Questions

- Whether release artifacts should be draft releases. Default for v0.1.0: publish non-draft on tag push.

## Verification Expectations

- `goreleaser check` passes, or equivalent `go run github.com/goreleaser/goreleaser/v2@latest check` passes if binary is not installed.
- `goreleaser release --snapshot --clean` or equivalent `go run ... release --snapshot --clean` builds snapshot artifacts locally.
- `go test ./...`, `go test -race ./...`, `go vet ./...`, `go mod tidy -diff`, and `go mod verify` pass.
