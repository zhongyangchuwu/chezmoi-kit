# Context: Phase 1 Output Controls

## Goal

Expose existing semantic report renderers through explicit CLI flags while preserving current default rendering and the README mental model.

## Constraints

- Do not change default command usage or default output behavior.
- Do not change existing command semantics, exit behavior, or mutation behavior.
- Limit scope to existing report-backed non-interactive commands: `cm status`, `cm diff`, and `cm version`.
- Keep README changes factual and limited to verified optional behavior.
- Preserve `NO_COLOR` behavior unless an explicit color flag intentionally overrides it.

## Decisions

- `v0.1.1` is a patch release because this phase only adds optional controls.
- Reuse the existing semantic report model and renderers; do not add a second output pipeline.
- Add flags at the CLI boundary; app services continue returning semantic report documents.

## Open Questions

- Exact flag placement and inheritance should follow Cobra conventions discovered in the existing CLI wiring.
- Explicit color precedence over `NO_COLOR` must be confirmed from implementation and tests.

## Verification Expectations

- Focused CLI tests cover explicit output selection, explicit color policy, default behavior preservation, and invalid values.
- Targeted Go tests run for affected packages.
- README is updated only after behavior is implemented and verified.
