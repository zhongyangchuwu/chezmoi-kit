# Verification: Phase 4 Release Gate

## Claims Checked

- Local release gate passes.
- GoReleaser config validates.
- GoReleaser snapshot release builds expected artifacts.
- Built binary supports safe non-mutating CLI smoke checks.
- Environment-limited smoke checks are recorded honestly.
- Release notes are dated for v0.1.0.

## Evidence Observed

| Claim | Evidence |
|---|---|
| Go tests pass | `go test ./...` passed. |
| Race tests pass | `go test -race ./...` passed. |
| Vet passes | `go vet ./...` passed with no output. |
| Module tidy clean | `go mod tidy -diff` passed with no output. |
| Modules verified | `go mod verify` reported all modules verified. |
| Vulnerability scan clean | `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` reported no vulnerabilities. |
| GoReleaser config valid | `go run github.com/goreleaser/goreleaser/v2@latest check` validated `.goreleaser.yaml`. |
| Snapshot release builds | `go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean` succeeded. |
| Expected artifacts exist | `dist/` contains six archives for linux/darwin/windows on amd64/arm64 plus `checksums.txt`, metadata, and extracted build directories. |
| Version smoke passes | `./dist/cm_linux_amd64_v1/cm version` printed snapshot version `0.0.1-next` with commit metadata. |
| Help smoke passes | `./dist/cm_linux_amd64_v1/cm --help` listed the v0.1.0 command surface including `edit`, `sync`, `git`, and `completion`. |
| Completion smoke passes | `./dist/cm_linux_amd64_v1/cm completion zsh` generated zsh completion script. |
| Status smoke passes | `./dist/cm_linux_amd64_v1/cm status` ran against the current chezmoi setup and returned local mismatch status. |
| Diff smoke passes | `./dist/cm_linux_amd64_v1/cm diff` ran against current dirty managed files and produced internal diffs. |
| Edit/sync help smoke passes | `./dist/cm_linux_amd64_v1/cm edit --help` and `./dist/cm_linux_amd64_v1/cm sync --help` returned command help successfully. |
| Changelog dated | `CHANGELOG.md` now has `## v0.1.0 - 2026-06-19`. |

## Coverage

- Local release gate is covered.
- GoReleaser release artifact generation is covered locally in snapshot mode.
- Safe CLI command smoke is covered for version, help, completion, status, diff, edit help, and sync help.
- Real mutating commands were not executed.

## Gaps

- `cm git` attempted to launch `lazygit` but failed because this harness has no `/dev/tty`; this verifies command wiring reaches lazygit but does not verify interactive lazygit usability.
- `cm sync` Ctrl+C terminal cleanup requires an interactive terminal and remains manual.
- `cm edit <known-managed-file>` was not run because it may open an editor and mutate files; only `cm edit --help` was smoke-tested here.
- Remote GitHub CI has not been observed until the branch is pushed and the workflow run completes.
- A real tag release has not been published; tag creation should wait for remote CI and user approval.

## Result

Local release verification passed. The remaining gates are remote GitHub CI, interactive terminal smoke, and final tag publication approval.
