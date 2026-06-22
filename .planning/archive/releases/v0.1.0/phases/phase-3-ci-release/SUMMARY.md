# Summary: Phase 3 CI and Release Automation

## Completed Changes

- Added `.goreleaser.yaml` using GoReleaser v2 configuration.
- Configured GoReleaser to build `cm` from `./cmd/cm` for Linux, macOS, and Windows on amd64 and arm64.
- Configured GoReleaser release ldflags to inject `github.com/zhongyangchuwu/cm/internal/build.Version={{ .Version }}`.
- Configured archives with `tar.gz` for Unix targets and `zip` for Windows targets.
- Configured release checksums as `checksums.txt`.
- Added `.github/workflows/ci.yml` for push/PR validation.
- Added `.github/workflows/release.yml` for `v*` tag publishing through `goreleaser/goreleaser-action@v7`.
- Updated `docs/development.md` with GoReleaser local checks, snapshot release validation, GitHub workflows, and tag release flow.

## Files Changed

- `.goreleaser.yaml`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `docs/development.md`
- `.planning/phases/phase-3-ci-release/CONTEXT.md`
- `.planning/phases/phase-3-ci-release/PLAN.md`
- `.planning/phases/phase-3-ci-release/SUMMARY.md`

## Deviations

- The earlier project decision to defer GoReleaser was superseded by the user's explicit request to use GoReleaser.
- Local verification used `go run github.com/goreleaser/goreleaser/v2@latest ...` instead of requiring a preinstalled `goreleaser` binary.
- Snapshot release validation produced `dist/` artifacts, which are ignored by `.gitignore`.

## Evidence

- `go test ./...` passed.
- `go test -race ./...` passed.
- `go vet ./...` passed with no output.
- `go mod tidy -diff` passed with no output.
- `go mod verify` reported all modules verified.
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` reported no vulnerabilities.
- `go run github.com/goreleaser/goreleaser/v2@latest check` validated `.goreleaser.yaml`.
- `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean` built six snapshot archives and checksums.
- `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean --skip=publish` matched the CI workflow command and passed.
- Go workspace diagnostics reported no issues.

## Unresolved Risks

- GitHub Actions workflows have not run remotely yet; Phase 4 must confirm CI passes after push.
- The tag release workflow has not been exercised against a real `v*` tag yet; Phase 4 owns final tag verification.
- Manual terminal smoke for `cm sync` Ctrl+C cleanup remains in Phase 4.
