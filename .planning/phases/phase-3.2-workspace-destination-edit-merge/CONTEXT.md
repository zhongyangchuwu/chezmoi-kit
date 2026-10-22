# Context: Workspace Destination Edit and Merge

## Goal

Extend `cm ui` with explicit destination-local editing and chezmoi three-way merge while retaining the source edit safety boundary from Phase 3.1.

## Accepted Decisions

- `E` edits the selected managed regular destination file. It prefers `VISUAL`, falls back to `EDITOR`, parses quoted arguments without invoking a shell, and passes the absolute destination path directly. It never edits source or invokes chezmoi apply/watch behavior.
- `E` rejects symlinks rather than following their referent, plus directories, scripts, unmanaged/ignored, remove/external/unknown, and unmapped entries.
- `m` opens `chezmoi merge` only for eligible managed regular files. It rejects symlinks, matching the established sync merge action matrix.
- Help and handoff status explain that merge receives destination, source, and rendered target. A merge tool may change source or destination; cm never applies after it returns.
- `e`, `E`, and `m` share one terminal handoff lifecycle. Every launched tool return re-inventories original scopes, preserves viable workspace context, invalidates previews/reveals/search/scroll state, and rejects stale async preview responses.
- Missing editor, merge preparation failure, tool cancellation/nonzero exit, and refresh failure leave the workspace usable. Only a launched tool return triggers refresh.

## Non-Goals

- No lazygit workspace handoff, workspace apply, action menu, source Git refresh, tabs, or embedded terminal.
