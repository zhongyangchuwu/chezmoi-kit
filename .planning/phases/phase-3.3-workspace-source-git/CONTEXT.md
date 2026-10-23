# Context: Workspace Source Git Handoff

## Goal

Make the existing `cm git` lazygit workflow available as `g` inside `cm ui`, returning to the same workspace instead of ending the command.

## Accepted Decisions

- `g` launches lazygit in the configured chezmoi source repository through the same app-owned source-directory resolution as direct `cm git`.
- `g` is available even when the current workspace directory has no selected file. It is an explicit terminal handoff, not a Git reimplementation.
- Workspace inventory records the source repository porcelain working-tree change count. Git status failures stay visible as a workspace status message but do not block file inventory.
- Every lazygit return, including nonzero exit, re-inventories original scopes, refreshes source Git summary, preserves viable workspace context, and invalidates stale previews/caches.
- Lazygit can change source membership, branch, or paths. Refresh uses chezmoi inventory as authority and applies existing selection fallback; cm makes no branch, commit, pull, or push automatically.

## Non-Goals

- No embedded lazygit terminal, Git action menu, source Git staging/commit implementation, workspace apply, or shell-directory change-on-exit integration.
