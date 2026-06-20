# Plan: 7-semantic-reports

## Objective

Add a semantic report abstraction and renderers so app-owned outputs keep meaning separate from presentation.

## Scope

In scope:

- Define a small semantic document model for command outputs.
- Implement plain text and ANSI rendering with a semantic palette.
- Implement Markdown rendering if it is useful for deterministic docs/agent output.
- Move status, diff, and version outputs to semantic reports.
- Share diff line classification between report rendering and TUI diff rendering.
- Add color behavior based on TTY detection and `NO_COLOR`.

Out of scope:

- Full Markdown parser integration.
- Glamour dependency unless the implementation proves custom rendering is too costly.
- JSON output mode.
- Large table/layout framework.
- Changing TUI layout or sync interaction semantics.

## Tasks

1. Add a report package or app-local report module with minimal semantic types:
   - `Document`
   - `Heading`
   - `Paragraph`
   - `List` if needed
   - `CodeBlock`
   - `DiffBlock`
   - inline tokens: `Text`, `Strong`, `Muted`, `Code`, `Path`, `Command`, `Status`.
2. Add diff line semantics:
   - `DiffLineKind`
   - `DiffLine`
   - `ClassifyDiffLine`.
3. Add renderers:
   - plain text renderer.
   - ANSI renderer using a palette.
   - Markdown renderer if accepted for command/agent output.
4. Add color policy:
   - auto color only for TTY stdout.
   - `NO_COLOR` disables ANSI color when non-empty.
   - leave room for future `--color=auto|always|never`.
5. Convert app status output to `StatusReport(targets) (Document, error)`.
6. Convert app diff output to `DiffReport(targets) (Document, error)` or a report-backed byte writer that preserves current output.
7. Convert version output to a report-backed representation.
8. Update CLI to render reports through the selected renderer and write only the rendered output.
9. Update TUI diff rendering to use shared diff line classification.
10. Add tests for document construction, plain/ANSI/Markdown rendering, NO_COLOR behavior, and diff line classification.
11. Update docs describing output architecture and color behavior.

## Acceptance Criteria

- App no longer hand-builds user-facing status/version strings without semantic tokens.
- Status output content and styling are generated through the report model.
- Diff line classification is defined once and reused by TUI/report rendering.
- Plain output remains readable and suitable for non-TTY capture.
- ANSI output uses semantic palette rules rather than hard-coded styling at call sites.
- Color is disabled when `NO_COLOR` is set and not empty.
- Markdown output exists or is explicitly deferred with rationale.
- Existing command behavior is preserved unless documented.

## Verification

Run after implementation:

```bash
go test ./internal/report ./internal/app ./internal/cli ./internal/tui
go test ./...
go vet ./...
go mod tidy -diff
```

Behavior checks:

```bash
go build -o /tmp/cm-report-check ./cmd/cm
/tmp/cm-report-check status
NO_COLOR=1 /tmp/cm-report-check status
/tmp/cm-report-check diff
/tmp/cm-report-check version
```

## Risks

- Overbuilding the document model can slow delivery. Keep the first model minimal and only add block types with immediate consumers.
- ANSI golden tests can be brittle. Prefer semantic renderer tests and ANSI/no-ANSI contains checks for key tokens.
- Markdown renderer may tempt content authors to embed raw markdown; keep app output construction semantic.
