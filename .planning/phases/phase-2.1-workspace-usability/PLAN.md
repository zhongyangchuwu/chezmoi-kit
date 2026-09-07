# Plan: Workspace Usability

## Objective

Make `cm ui` easy to scan and operate without external documentation by adding semantic visual hierarchy, legend, responsive contextual feedback, and a keyboard-isolated help overlay.

## Scope

1. Move visual styling from key bindings into a focused semantic palette.
2. Render file state/type/attribute tokens with semantic color and retain textual markers.
3. Add selected-row and active-pane hierarchy without hiding focus from non-color users.
4. Add a compact in-pane legend and responsive footer.
5. Add a workspace-only `?` help overlay with quick start, keys by focus, preview guidance, and complete state/type legend.
6. Style preview metadata, loading, error, withheld, and clean states distinctly.
7. Add behavior-focused TUI coverage and real PTY checks including `NO_COLOR=1` fallback.

## Non-goals

- Theme selection, persistent preferences, mouse controls, icon fonts, localization, command palette, external-tool handoffs, or workspace mutations.

## Acceptance Criteria

- `C/D/U/I/R/?`, `f/d/l/s/x/e`, `[T]`, and `[E]` remain readable without color and gain distinct semantic visual treatment when color is available.
- Current selection and active focus are visually distinguishable at normal and narrow sizes.
- `?` shows complete quick-start, key, and legend information; `Esc`, `?`, and `q` close it; other keys do not change workspace state behind the overlay.
- Footer remains usable at 24, 60, 80, and 120 columns, retaining a help and quit route.
- Loading, error, withheld, and clean preview feedback are visibly distinct.
- Existing read-only, sensitive reveal, search, projection, preview, and sync contracts remain unchanged.

## Verification

- `go test ./internal/tui` and `go test ./...`.
- `go test -race ./...`, `go vet ./...`, and `git diff --check` before completion.
- Real PTY checks at normal/narrow widths and with `NO_COLOR=1`; quit still leaves isolated source/destination unchanged.

## Risks

- ANSI styling can invalidate width assumptions; all renderers continue using display-width-aware truncation and narrow-width tests.
- Extra style branches can obscure semantics; centralize palette decisions and preserve token text.
- Overlay handling can accidentally trigger actions; route it before normal workspace key dispatch.
