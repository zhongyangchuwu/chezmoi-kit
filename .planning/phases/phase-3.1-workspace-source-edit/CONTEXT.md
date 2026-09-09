# Context: Workspace Source Edit

## Goal

Make `cm ui` the primary interactive entrypoint for safe source edits while retaining `cm edit <target>` for direct shell use.

## Accepted Decisions

- `e` opens the selected managed file or symlink through `chezmoi edit`; Enter continues to mean directory navigation only.
- Workspace source edit forces `--apply=false --watch=false`, overriding inherited chezmoi settings so source edits return to reviewed diff rather than implicitly mutating destination.
- Supported: managed regular files, symlinks, templates, encrypted files, and uninspected managed files with a source mapping.
- Unsupported: directories/virtual nodes, unmanaged, ignored, scripts, remove/external/unknown entries, or entries without a source mapping. The UI explains why.
- Editor return always re-inventories the workspace from the original scopes. Preserve viable directory/filter/focus/preview-kind/fullscreen/selection state; clear stale preview, reveal, search-match, and scroll state.
- Nonzero editor exit still refreshes. Refresh failure leaves the prior inventory visible but clears stale previews and reports both outcomes.
- Build a reusable terminal-handoff/refresh path, but expose only `e` in this slice. No local editor, merge, lazygit, workspace apply, or mutation menu.

## Interaction States

- Idle: eligible file footer exposes `e edit source`.
- Preparing: input locked, status says `opening source editor...`.
- Handoff: configured editor owns terminal via Bubble Tea exec.
- Refreshing: input locked, status says `refreshing workspace...`.
- Returned: refreshed authoritative preview loads, with success/nonzero editor outcome message.
- Error: unavailable target, command construction, editor failure, or inventory refresh failure keeps TUI usable.

## Verification

- Service tests prove absolute target, `--apply=false`, and `--watch=false` command construction plus eligibility rejection.
- TUI tests prove allowed/disallowed `e`, terminal return refresh, selection fallback, stale sensitive-cache clearing, and nonzero/refresh-failure behavior.
- A disposable PTY fixture changes a managed source and confirms refresh plus no destination auto-apply. Real service integration covers template and encrypted source edits without destination changes.
