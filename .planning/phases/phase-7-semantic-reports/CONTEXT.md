# Context: Phase 7 Semantic Reports and Palette Rendering

## Goal

Introduce a semantic information model so status, diff, version, and future outputs separate content meaning from plain, ANSI, Markdown, and TUI presentation.

## Constraints

- Do not use Markdown as the internal model. Markdown is an output format only.
- Keep the first report model small and domain-driven; add tables or rich constructs only when needed.
- Preserve current default human output unless an acceptance criterion explicitly changes it.
- Respect stdout/stderr separation: command data goes to stdout; debug/log paths and errors go to stderr.
- Respect `NO_COLOR` and TTY detection for ANSI color behavior.
- Do not add Glamour unless the internal renderer becomes too costly or Markdown-to-ANSI quality is a concrete requirement.

## Decisions

- App should produce semantic report documents for non-interactive outputs.
- Renderers own output formats: plain text, Markdown, ANSI/palette, and TUI-specific rendering helpers.
- Style is semantic: warning, strong, muted, path, command, diff add/remove/header/hunk/meta.
- `internal/diff` may continue producing raw unified diff bytes, but report rendering should classify diff lines once and share that classification with TUI.
- Initial implementation can use a small internal `report` package rather than a Markdown parser.

## Open Questions

- Whether the report package should live at `internal/report` or inside `internal/app/report`. Default: `internal/report`, because it is a cross-interface formatting capability.
- Whether CLI should expose `--output markdown` immediately or only implement internal Markdown rendering for future use.
- Whether ANSI rendering should use existing Lip Gloss styles directly or a renderer-local palette abstraction.

## Verification Expectations

- Status report rendering is covered in plain and ANSI/no-color modes.
- Diff line classification is shared by CLI/report and TUI, with tests for headers, hunks, add/remove, and metadata lines.
- Markdown renderer output is deterministic if implemented.
- `NO_COLOR` disables color when color is otherwise automatic.
