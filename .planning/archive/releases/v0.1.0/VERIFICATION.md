# v0.1.0 Verification

## Result
Passed with known gaps

## Claims Verified
- All seven included phases have `SUMMARY.md` and `VERIFICATION.md` evidence.
- Current requirements SAFE-01, SAFE-02, SAFE-03, CLI-01, CLI-02, REL-01, REL-02, REL-03, REL-04, DOC-01, DOC-02, DOC-03, VER-01, VER-02, ARCH-01, ARCH-02, ARCH-03, TUI-01, OUT-01, OUT-02, and TEST-01 are complete and mapped.
- Local release gate passed before publication.
- Remote GitHub CI passed for commit `f8e80afb408b93a53e743a58be52d31f9cf842f4` on `main`.
- Tag `v0.1.0` points to commit `f8e80afb408b93a53e743a58be52d31f9cf842f4`.
- GitHub Release `v0.1.0` is published, non-draft, and non-prerelease.
- Release workflow ran GoReleaser successfully for tag `v0.1.0`.

## Evidence
| Phase | Verification |
| --- | --- |
| phase-1-runtime-stability | phases/phase-1-runtime-stability/VERIFICATION.md |
| phase-2-release-docs | phases/phase-2-release-docs/VERIFICATION.md |
| phase-3-ci-release | phases/phase-3-ci-release/VERIFICATION.md |
| phase-4-release-gate | phases/phase-4-release-gate/VERIFICATION.md |
| phase-5-package-architecture | phases/phase-5-package-architecture/VERIFICATION.md |
| phase-6-app-services | phases/phase-6-app-services/VERIFICATION.md |
| phase-7-semantic-reports | phases/phase-7-semantic-reports/VERIFICATION.md |

Additional observed release evidence:

- `git status --short --branch` reported `main...origin/main`, with no staged, unstaged, or untracked files before release archival work.
- `git rev-parse HEAD` reported `f8e80afb408b93a53e743a58be52d31f9cf842f4`.
- `git describe --tags --exact-match HEAD` reported `v0.1.0`.
- `git rev-list -n 1 v0.1.0` reported `f8e80afb408b93a53e743a58be52d31f9cf842f4`.
- `git ls-remote --tags origin v0.1.0` reported the remote `v0.1.0` tag.
- GitHub Actions run `27878479411` (`CI`) completed successfully for commit `f8e80afb408b` on `main`.
- GitHub Actions run `27878480851` (`Release`) completed successfully for tag `v0.1.0`.
- `gh release view v0.1.0 --json tagName,name,isDraft,isPrerelease,publishedAt,url,targetCommitish` reported `isDraft:false`, `isPrerelease:false`, `publishedAt:2026-06-20T17:27:10Z`, and URL `https://github.com/zhongyangchuwu/chezmoi-kit/releases/tag/v0.1.0`.

## Coverage
- Requirements covered: SAFE-01, SAFE-02, SAFE-03, CLI-01, CLI-02, REL-01, REL-02, REL-03, REL-04, DOC-01, DOC-02, DOC-03, VER-01, VER-02, ARCH-01, ARCH-02, ARCH-03, TUI-01, OUT-01, OUT-02, TEST-01.
- Automated checks covered across phase evidence: targeted Go tests, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go mod tidy -diff`, `go mod verify`, `govulncheck`, GoReleaser config check, GoReleaser snapshot release builds, stale-reference searches, Go diagnostics, version build smoke, `NO_COLOR` smoke, status smoke, and diff smoke.
- Remote checks covered: GitHub CI run on `main` and GitHub Release workflow run for `v0.1.0` tag.
- Manual/smoke checks covered: safe non-mutating binary checks for version, help, completion, status, diff, edit help, and sync help.

## Known Gaps
- `cm sync` Ctrl+C terminal cleanup requires a real interactive terminal; production signal handling was fixed and automated/TUI state tests passed, but this harness did not provide final manual terminal evidence.
- `cm git` lazygit startup requires `/dev/tty`; harness evidence reached lazygit command wiring but could not prove full interactive lazygit usability.
- `cm edit <known-managed-file>` was not executed because it may open an editor and mutate files; command wiring and `cm edit --help` were verified.
- No user-facing `--output` or `--color` flags exist yet; report renderers support plain, ANSI, and Markdown internally.
- TUI file/diff layout still truncates by byte length rather than display width.
