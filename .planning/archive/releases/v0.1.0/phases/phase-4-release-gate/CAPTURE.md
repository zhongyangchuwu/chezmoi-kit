# Capture: Phase 4 Release Gate

## Durable Docs Updated

- `CHANGELOG.md` carries the `v0.1.0` release date.
- `.goreleaser.yaml` targets the actual GitHub repository, `zhongyangchuwu/chezmoi-kit`.

## Planning Records Updated

- Phase 4 verification evidence records the local release gate, GoReleaser snapshot artifacts, safe CLI smoke checks, and environment-limited manual checks.
- Roadmap, requirements, project state, and workflow state mark Phase 4 complete and keep release publication approval as a ship-time gate.

## Learnings

- The harness can verify Go checks, GoReleaser config, snapshot artifacts, and non-mutating CLI behavior.
- Interactive lazygit and Ctrl+C TUI cleanup require a real terminal; these are documented environment-limited smoke gaps, not local gate failures.
- Remote GitHub CI cannot be proven locally; tag publication remains gated on observed CI and user approval.

## Ship Inputs

- Before publishing `v0.1.0`, observe remote GitHub CI on the pushed branch.
- Before publishing `v0.1.0`, run or explicitly waive real-terminal smoke for `cm sync` Ctrl+C cleanup and `cm git` lazygit startup.
- Tag creation and release publication require explicit user approval.
