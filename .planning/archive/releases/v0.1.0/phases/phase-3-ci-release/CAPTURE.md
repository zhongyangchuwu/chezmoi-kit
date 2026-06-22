# Capture: Phase 3 CI and Release Automation

## Durable Docs Updated

- `docs/development.md` now documents GoReleaser checks, CI workflow, release workflow, and tag flow.
- `.goreleaser.yaml` is the release artifact configuration.
- `.github/workflows/ci.yml` is the push/PR validation workflow.
- `.github/workflows/release.yml` is the tag release workflow.

## Planning Records Updated

- `.planning/phases/phase-3-ci-release/CONTEXT.md`
- `.planning/phases/phase-3-ci-release/PLAN.md`
- `.planning/phases/phase-3-ci-release/SUMMARY.md`
- `.planning/phases/phase-3-ci-release/VERIFICATION.md`
- `.planning/phases/phase-3-ci-release/CAPTURE.md`

## Learnings

- GoReleaser v2 config validates and snapshot-builds successfully using `go run github.com/goreleaser/goreleaser/v2@latest`.
- The CI snapshot command `release --snapshot --clean --skip=publish` is valid locally.
- GoReleaser snapshot builds tolerate the current repository having no tags when `--snapshot` is used.

## Ship Inputs

- Phase 4 must confirm GitHub CI passes remotely.
- Phase 4 must tag `v0.1.0` only after final local release gate and manual smoke checks.
- Phase 4 should update `CHANGELOG.md` from `Unreleased` to the release date before tagging.
