# Research: Semantic Reports and Terminal Rendering

## Question

How should `cm` represent user-facing command information so status, diff, version, and future TUI/agent output can share content semantics while rendering to plain text, ANSI, Markdown, or TUI layouts?

## Sources Checked

- `https://clig.dev/` — Command Line Interface Guidelines.
- `https://no-color.org/` — NO_COLOR convention.
- `https://github.com/charmbracelet/glamour` — Markdown-to-ANSI renderer for Go CLI apps.
- `https://github.com/charmbracelet/lipgloss` — Terminal styling/layout library already used by the project.
- Current code discussion and codebase inspection: `internal/cli/status.go`, `internal/cli/diff.go`, `internal/tui/diff.go`, `internal/diff/diff.go`.

## Findings

- CLI Guidelines recommend primary command output on stdout and diagnostic/log/error messaging on stderr. This supports app-owned command output with CLI as stream plumbing.
- CLI Guidelines frame modern CLIs as human-first but composable, which supports stable plain/Markdown output and optional structured rendering rather than TUI-only formatting.
- NO_COLOR states that software adding ANSI color by default should disable color when `NO_COLOR` is present and non-empty; command flags/config may override it.
- Glamour provides stylesheet-based Markdown rendering for Go CLI apps and is used by larger CLIs, but using Markdown as the internal model would lose app-specific semantics like path, command, status severity, and diff line kind.
- Lip Gloss already supports semantic terminal styling, ANSI colors, text attributes, and terminal-aware downsampling. It remains suitable for TUI and custom ANSI renderers.

## Tradeoffs

### Markdown as internal model

Pros:

- Simple text format.
- Human and agent friendly.
- Can be rendered with Glamour later.

Cons:

- Loses domain semantics needed by TUI and future formats.
- Requires parsing Markdown if another renderer needs structured data.
- Couples content construction to presentation markup.

### Semantic document model as internal model

Pros:

- Preserves meaning: heading, paragraph, path, command, warning, strong, diff line kind.
- Supports plain, Markdown, ANSI, JSON-like future formats, and TUI reuse.
- Lets style palette map semantics to output-specific rendering.

Cons:

- Requires a small custom document/rendering layer.
- Can become over-engineered if too many block types are added early.

### Add Glamour now

Pros:

- Mature Markdown-to-ANSI rendering.
- Custom styles and wrapping already implemented.

Cons:

- New dependency.
- Does not remove need for a semantic app document model.
- Diff-specific coloring and TUI reuse are more direct with a semantic renderer.

## Confidence

High confidence for the core direction: app should produce semantic report documents and renderers should map them to plain/ANSI/Markdown/TUI output. Medium confidence on exact package placement (`internal/report` versus `internal/app/report`); choose during implementation based on import direction after Phase 6.
