# Capture: 7-semantic-reports

## Durable Docs Updated

- `docs/design.md` documents `internal/report`, semantic documents, ANSI/NO_COLOR behavior, Markdown as renderer only, and shared diff classification.
- `docs/development.md` documents the new `internal/report` package, report responsibilities, and dependency changes.
- `.planning/codebase/ARCHITECTURE.md`, `CONVENTIONS.md`, `STACK.md`, `CONCERNS.md`, and `MAP.md` now reflect semantic report ownership and remaining UI risks.

## Planning Records Updated

- `.planning/phases/phase-7-semantic-reports/SUMMARY.md` records implementation changes, files changed, deviations, evidence, and risks.
- `.planning/phases/phase-7-semantic-reports/VERIFICATION.md` records success-criteria evidence and requirement coverage.
- Root planning artifacts were updated to mark Phase 7 complete and project workflow complete.

## Learnings

- A small semantic model was sufficient for current outputs; adding table/layout abstractions now would add maintenance cost without an immediate consumer.
- Markdown can be a deterministic renderer without becoming the internal model.
- Sharing diff classification between report rendering and TUI styling removes duplicated prefix rules while leaving TUI layout unchanged.
- Removing `fatih/color` simplified output policy: all command color now flows through one renderer with `NO_COLOR` handling.

## Ship Inputs

- Existing external release gates remain: observed remote CI, explicit tag approval, and real-terminal smoke where available.
- Future output flags can reuse `report.Options` for `--color=auto|always|never` and `report.Markdown` for an explicit Markdown mode.
