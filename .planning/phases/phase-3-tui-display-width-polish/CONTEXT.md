# Context: Phase 3 TUI Display Width Polish

## Goal

Fix sync TUI display truncation so file names, status text, and diff lines are clipped by terminal display width without splitting UTF-8 runes.

## Constraints

- Do not change sync workflow, key bindings, pending-action model, confirmation, or preflight behavior.
- Keep the change local to TUI rendering helpers where possible.
- Preserve existing ASCII rendering output unless display-width limits require truncation.
- Add focused tests for Unicode and narrow widths.

## Decisions

- This is a `v0.1.1` patch fix because it preserves behavior and improves display correctness.
- Use a small local helper rather than introducing broad layout abstractions.

## Open Questions

- Existing dependencies may already include display-width helpers through TUI libraries; inspect before adding dependencies.

## Verification Expectations

- TUI tests cover Unicode paths/diff lines and very narrow widths.
- Existing sync TUI tests remain passing.
- `go test ./internal/tui` and `go test ./...` pass.
