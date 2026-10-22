# Review: Workspace Destination Edit and Merge

**Date:** 2026-10-23  
**Result:** Approved for PR.

## Reviewed Surface

- Destination editor selection/parsing, eligibility, merge command construction, and real chezmoi integration in `internal/app`.
- Shared Bubble Tea handoff lifecycle, operation-specific outcome text, key routing, contextual action footer, help, and recovery tests in `internal/tui`.
- User documentation and Phase 3 requirement/roadmap state.

## Findings Resolved

1. **Merge initially accepted symlinks.** The established sync action matrix permits merge only for regular files. Workspace merge now rejects symlinks and its test covers the gate.
2. **Workspace merge command duplicated sync merge construction.** Both routes now use one app-owned `mergeTerminalCommand` builder.
3. **Full action footer could hide individual handoffs at medium widths.** The footer now tries progressive action sets before falling back to navigation-only guidance, while preserving help and quit routes.

## Residual Risks

- `E` executes the user-configured `VISUAL` or `EDITOR` program directly. Shell syntax, environment assignments, and expansion are intentionally unsupported; quoted executable arguments are supported and target paths are never shell-interpreted.
- A configured merge tool may modify source or destination. cm states all three merge roles, does not apply after return, and refreshes authoritative inventory before rendering a new preview.
