# Context: Workspace Usability

## Goal

Make the read-only `cm ui` workspace visually scannable and self-explanatory with semantic status colors, a persistent legend, contextual feedback, and an in-TUI `?` help overlay.

## Constraints

- Keep `cm ui` read-only and preserve all Phase 1/2 safety, preview, reveal, and navigation behavior.
- Existing status/type letters remain the non-color semantic source; color augments rather than replaces them.
- Keep one shared Bubble Tea model/rendering path with explicit workspace-only branches.
- Do not add themes, config, fonts/icons, mouse behavior, persistent onboarding state, or lifecycle actions.
- Help must be accessible at any size and must not let keyboard input reach workspace actions while visible.

## Decisions

- Use a small semantic style palette: clean, dirty, unmanaged, ignored, script, uninspected, template, encrypted, directory, symlink, error, loading, withheld, and selected row.
- Color only meaningful tokens/badges and selection chrome; paths retain a readable default color.
- Add `?` to toggle a workspace help overlay. `Esc`, `?`, and `q` close it; no other workspace action runs while it is open.
- Render status/type legend in the file pane when room permits; help always includes the complete legend and quick-start sequence.
- Adapt the footer by terminal width, retaining `? help` and `q quit` at narrow sizes.

## Open Questions

- None blocking. Use the existing fixed palette and Lip Gloss terminal profile behavior; do not introduce user configuration.

## Verification Expectations

- TUI tests cover semantic labels/styles without relying only on ANSI bytes, help isolation/closing, complete legend, responsive footer, and narrow display widths.
- Actual PTY smoke verifies normal and narrow layouts, focus/selection contrast, help overlay, state visibility, and `NO_COLOR=1` text fallback.
- Existing workspace, sync, and full Go regressions remain green.
