# Plan: Phase 1 Output Controls

## Objective

Add explicit output and color controls for existing report-backed non-interactive commands while keeping defaults unchanged.

## Scope

In scope:

- `cm status [target...]`
- `cm diff [target...]`
- `cm version`
- CLI rendering policy and tests.
- README documentation for optional flags after implementation.

Out of scope:

- `cm sync` TUI rendering controls.
- `cm completion` output format.
- New report document types beyond what the phase needs.
- Any change to default output, status model, sync flow, or mutation behavior.

## Tasks

1. Inspect existing CLI render wiring and report renderer APIs.
2. Add typed CLI options for `--output plain|ansi|markdown` and `--color auto|always|never`.
3. Thread options through `status`, `diff`, and `version` command rendering without changing app service contracts.
4. Add command tests for defaults, explicit output formats, explicit color behavior, and invalid values.
5. Update README with verified optional flag usage.
6. Run targeted verification and record evidence.

## Acceptance Criteria

- `cm status`, `cm diff`, and `cm version` accept output selection for supported renderers.
- Explicit color policy can force ANSI, disable ANSI, or keep existing auto behavior.
- Existing default behavior remains unchanged for non-TTY output and `NO_COLOR` defaults.
- Invalid output or color values fail with clear CLI errors.
- README documents optional flags accurately.

## Verification

- Run focused CLI/report tests for changed behavior.
- Run broader affected package tests if implementation touches shared report code.
- Smoke command output where practical with `go run ./cmd/cm` for flag parsing and user-visible behavior.

## Risks

- Global flags could accidentally affect Cobra completion output if attached too broadly.
- `NO_COLOR` precedence can surprise users if explicit color flags are not tested.
- Markdown output for command reports must use existing report semantics rather than ad hoc string rendering.
