# Summary: Directory Browser

## Completed Changes

- Replaced workspace recursive tree/flat/collapse projections with a Yazi-inspired direct-child browser.
- Added virtual directory ancestors derived from authoritative inventory paths, aggregate state for virtual nodes, per-directory cursor restoration, initial single-scope entry, and global search locating.
- Added Parent/Current/Preview wide layout, Current/Preview medium layout, and focus-switched narrow layout.
- Removed disclosure glyphs and fixed all rows to `selection + state:type + name`; directories use a trailing `/`.
- Added focus-sensitive `h`/`l`/Enter navigation, local directory summary previews, ancestor-preserving filters, and updated help/footer guidance.

## Preserved Behavior

- Workspace is read-only.
- Chezmoi remains authoritative for inventory, file classification, source paths, attributes, rendered targets, and file previews.
- Sensitive reveal, bounded previews, no-color semantics, full preview, preview search, hunk navigation, and sync behavior are unchanged.

## Evidence

- Focused and full Go tests, race tests, vet, tidy diff, whitespace, and TUI LSP diagnostics passed.
- PTY smoke verified wide Parent/Current/Preview, entering/leaving a directory, keyboard-isolated help, no-color markers, and narrow Current-only fallback with preview focus available through Tab.
- Source/destination archive hash before and after workspace exits: `f4c92711ff0f5bf1e2e50969f52bd5d083854777d24697d841ec8a91f81bb804`.

## Risks

- Virtual nodes are rebuilt from the full bounded inventory per projection; exceptionally large inventories may make projection changes noticeable.
- Very narrow terminals intentionally omit context panes and shorten guidance.
