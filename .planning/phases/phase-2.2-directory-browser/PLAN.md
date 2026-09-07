# Plan: Directory Browser

## Navigation Model

1. Derive direct-child directory views from authoritative workspace inventory; synthesize missing path ancestors as virtual directories.
2. Track current directory, per-directory cursor, and search origin. Remove flat/tree/collapse fields and keys.
3. Give virtual directories aggregate semantic state and local summary previews; never request a chezmoi content preview for virtual navigation nodes.
4. Apply filters to files while retaining ancestor directories with matching descendants.

## Interaction

1. Map `h`/left to parent and `l`/right/enter to selected-directory entry while file focus is active.
2. Keep preview focus scrolling and existing preview/reveal/full-screen behavior intact.
3. Implement global path-search results; Enter locates the selected result, Escape restores origin.
4. Update contextual help and footer to accurately expose focus-sensitive navigation.

## Rendering

1. Replace the recursive Files pane with Current pane and optional Parent pane.
2. Render aligned fixed columns without disclosure glyphs; directory names end in `/`.
3. Implement three-, two-, and one-pane responsive layouts.
4. Render local directory summaries in Preview.

## Verification

1. Replace tree/collapse coverage with directory navigation/filter/search tests.
2. Run focused/full/race Go tests, vet, tidy diff, whitespace, and TUI diagnostics.
3. Exercise the real PTY fixture, including no-color and no-mutation hash checks.
4. Capture planning, design, usage, development, changelog, and codebase records.
