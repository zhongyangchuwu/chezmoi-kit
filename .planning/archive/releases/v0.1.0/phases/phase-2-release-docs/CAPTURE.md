# Capture: Phase 2 Release Metadata and Documentation Cleanup

## Durable Docs Updated

- `README.md`
- `docs/usage.md`
- `docs/design.md`
- `docs/development.md`
- `CHANGELOG.md`
- `LICENSE`
- `.gitignore`
- `justfile`

## Planning Records Updated

- `.planning/phases/phase-2-release-docs/CONTEXT.md`
- `.planning/phases/phase-2-release-docs/PLAN.md`
- `.planning/phases/phase-2-release-docs/SUMMARY.md`
- `.planning/phases/phase-2-release-docs/VERIFICATION.md`
- `.planning/phases/phase-2-release-docs/CAPTURE.md`

## Learnings

- Explicit ldflag version was previously overridden by Go module VCS version metadata; release builds need `internal/build.Current` to prefer explicit `Version` when it is not `dev`.
- Public docs are now small and current; historical planning details live only under `.planning/`.
- The release flow is now documented locally, but not yet automated on GitHub.

## Ship Inputs

- `v0.1.0` release notes already have a changelog draft.
- Phase 3 should add CI and tag-release workflows that run the documented release gate.
- Phase 4 should replace `Unreleased` with the tag date before publishing.
