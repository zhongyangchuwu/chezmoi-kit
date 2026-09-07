# Research: Phase 2 Unified File Workspace

## Question

Which chezmoi commands and current cm seams can support a complete read-only inventory and preview workspace without reimplementing chezmoi state rules or duplicating the sync TUI?

## Sources Checked

- Current app/chezmoi/TUI/CLI code and Phase 1 verification artifacts.
- Official chezmoi docs for `managed`, `unmanaged`, `ignored`, `status`, `cat`, `decrypt`, `source-path`, common type filters, and building frontends.
- Isolated chezmoi v2.72.1 probes for path-style-all JSON, template/type filters, scoped unmanaged recursion, ignored output, and secret-skip behavior.

## Findings

- `managed --path-style=all --format=json` provides relative, absolute, source-relative, and source-absolute mappings without parsing source filenames.
- Type-filtered `managed` queries expose dirs, symlinks, scripts, removes, externals, templates, and encrypted membership without loading file contents.
- `unmanaged <scope>` is authoritative but can return a directory as one candidate; recursively querying returned directories is required to expose nested file candidates while retaining chezmoi filtering.
- `ignored -0` reports source entries excluded from target state as relative target paths. It does not explain arbitrary destination-only ignored files or provide source mappings; cm must label only the fact actually returned.
- `cat` returns target contents for files/scripts and link targets for symlinks. `decrypt` can explicitly reveal encrypted source content.
- `--skip-secrets` skips secret-backed templates entirely rather than redacting values, so skipped output cannot be treated as clean or complete content.
- Current TUI is a dirty-entry model that exits when entries are clean. It can be generalized, but persistent workspace behavior must be mode-specific.
- Current diff rendering truncates long lines and provides only vertical scrolling. The ANSI package already supports display-width-aware left truncation needed for horizontal cropping.

## Tradeoffs

- Dump all target state once: fewer subprocesses, but loads rendered contents/secrets and can be unbounded. Rejected.
- Multiple path/type metadata commands: more subprocesses, but bounded and does not retain all contents. Selected.
- Scan the filesystem and reproduce ignore rules: faster display control but duplicates chezmoi semantics. Rejected.
- Show only unmanaged directories: simpler, but fails the accepted nested new-file scenario. Selected bounded recursive chezmoi queries instead.
- Make `cm ui` mutating immediately: reduces switching but couples inventory, projection, and execution before their state contract is stable. Rejected for this phase.
- Separate new TUI: initially simpler but creates duplicate focus/render/search behavior. Rejected; share the workbench model with explicit sync/workspace modes.

## Confidence

High for command shapes and current seams: official docs plus isolated real-tool probes. Medium-high for exact terminal ergonomics until PTY verification exercises normal and narrow layouts.
