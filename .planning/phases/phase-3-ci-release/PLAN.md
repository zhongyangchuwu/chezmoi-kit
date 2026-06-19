# Plan: 3-ci-release

## Objective

Automate validation and tag-based release publishing through GitHub Actions and GoReleaser.

## Scope

In scope:

- `.goreleaser.yaml` for v0.1.0 artifacts.
- `.github/workflows/ci.yml` for PR/push checks.
- `.github/workflows/release.yml` for `v*` tag releases.
- Documentation updates for GoReleaser release flow.
- Local verification of GoReleaser config and snapshot build.

Out of scope:

- Signing/notarization.
- Homebrew/Scoop/Winget packages.
- Docker images.
- GoReleaser Pro features.
- Actual tag creation or GitHub release publishing.

## Tasks

1. Add `.goreleaser.yaml` with GoReleaser v2 config for `cm`.
2. Add CI workflow that runs Go release gate and snapshot release validation.
3. Add tag release workflow using `goreleaser/goreleaser-action@v7`.
4. Update docs/development.md with GoReleaser local and GitHub release flow.
5. Verify GoReleaser config and snapshot build locally.
6. Write phase summary and verification artifacts.

## Acceptance Criteria

- `.goreleaser.yaml` builds Linux, macOS, and Windows artifacts for amd64 and arm64.
- GoReleaser injects `internal/build.Version={{ .Version }}`.
- CI workflow runs on push and pull request.
- Release workflow runs on `v*` tags and has `contents: write` permission.
- Development docs explain GoReleaser snapshot check and tag release flow.
- Local GoReleaser check and snapshot build pass.

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
```

## Risks

- Running GoReleaser via `go run` may add module cache downloads during verification but should not modify `go.mod`.
- GoReleaser config syntax can drift by version; docs checked were for current v2 action/config.
- Snapshot release creates `dist/`, which is ignored by `.gitignore`.
