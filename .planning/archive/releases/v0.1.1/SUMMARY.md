# v0.1.1 Summary

## Released
2026-06-22

## What Shipped
`cm` v0.1.1 shipped safe usability polish after the initial public release. The release added optional report output/color controls, a read-only environment doctor command, and display-width-aware sync TUI truncation while preserving existing defaults, command semantics, and mutation behavior.

## Included Phases
- `phase-1-output-controls` — Output Controls.
- `phase-2-doctor-diagnostics` — Doctor Diagnostics.
- `phase-3-tui-display-width-polish` — TUI Display Width Polish.

## Completed Scope

### Roadmap Items
- Phase 1: Exposed existing semantic report renderers through explicit user flags while preserving current default rendering.
- Phase 2: Added `cm doctor` as a read-only prerequisite diagnostic for required and optional external tools.
- Phase 3: Replaced byte-length sync TUI truncation with display-width-aware clipping for wide-character paths and diff lines.

### Requirements
- OUT-FLAGS-01 — Users can explicitly select plain, ANSI, or Markdown report output for existing non-interactive report-backed commands without changing default output behavior.
- OUT-FLAGS-02 — Users can explicitly select color policy (`auto`, `always`, `never`) while preserving current TTY auto-detection and `NO_COLOR` defaults.
- DOCTOR-01 — `cm doctor` provides read-only prerequisite and environment diagnostics for required and optional external tools without mutating local files or chezmoi source state.
- TUI-WIDTH-01 — `cm sync` TUI truncates file names, status text, and diff lines by display width rather than byte length, preserving valid UTF-8 and improving wide-character rendering.

## Notable Decisions
- Use `v0.1.Z` for safe additions/fixes and reserve `v0.Y.0` for behavior, safety-model, status-model, output-contract, or command-semantics changes; outcome: `v0.1.1` remained a low-surprise usability polish release.
- Implement output and color controls as CLI render options over the existing semantic report model; outcome: app services continue returning `report.Document` values.
- Treat `cm doctor` as diagnostic-only; outcome: it reports required failures and optional warnings without remediation or mutation.
- Keep TUI width handling local to rendering; outcome: sync workflow state, pending actions, confirmation, and execution contracts were unchanged.

## Follow-ups
- Template operation model planning: decide how `cm sync` should show destination, rendered target, source template, and merge/edit behavior for template-backed targets.
- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and app sync services.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.
- Issue templates, PR templates, and security policy files.
