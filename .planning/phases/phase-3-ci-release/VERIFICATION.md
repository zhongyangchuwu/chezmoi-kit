# Verification: Phase 3 CI and Release Automation

## Claims Checked

- GoReleaser config exists and is valid.
- GoReleaser builds release artifacts for Linux, macOS, and Windows on amd64 and arm64.
- Release builds inject the GoReleaser version into `internal/build.Version`.
- CI workflow runs the expected release gate on push and pull request.
- Release workflow runs GoReleaser on `v*` tags with the required GitHub token permission.
- Development docs explain the GoReleaser release flow.

## Evidence Observed

| Claim | Evidence |
|---|---|
| GoReleaser config valid | `go run github.com/goreleaser/goreleaser/v2@latest check` validated `.goreleaser.yaml`. |
| Cross-platform artifacts build | `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean` built linux/darwin/windows archives for amd64/arm64. |
| Version injection configured | `.goreleaser.yaml` ldflags set `github.com/zhongyangchuwu/cm/internal/build.Version={{ .Version }}`. |
| CI command valid locally | `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean --skip=publish` passed, matching the CI snapshot command. |
| Go tests pass | `go test ./...` passed. |
| Race tests pass | `go test -race ./...` passed. |
| Vet passes | `go vet ./...` passed with no output. |
| Module tidy clean | `go mod tidy -diff` passed with no output. |
| Module verification passes | `go mod verify` reported all modules verified. |
| Vulnerability scan clean | `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` reported no vulnerabilities. |
| Workflow docs updated | `docs/development.md` now documents GoReleaser local checks, CI workflow, release workflow, and tag flow. |
| Go diagnostics clean | Go workspace diagnostics reported no issues. |

## Coverage

- Covered Phase 3 success criteria 1-4 from `.planning/ROADMAP.md`.
- Covered requirements REL-03 and REL-04 at repository-file and local-command level.
- Verified workflow YAML exists and matches GoReleaser's documented GitHub Action usage: `fetch-depth: 0`, `goreleaser/goreleaser-action@v7`, `version: "~> v2"`, and `GITHUB_TOKEN` for release publishing.

## Gaps

- Remote GitHub Actions have not run yet because the branch/tag has not been pushed.
- Release publishing has not been exercised against a real `v*` tag.
- Manual terminal smoke remains for Phase 4.

## Result

Phase 3 is verified locally and ready to transition to Phase 4: Release Verification.
