# Capture: Workspace Source Git Handoff

**Merged:** 2026-10-23  
**Pull request:** [#10](https://github.com/zhongyangchuwu/chezmoi-kit/pull/10)  
**Squash commit:** `f3c1538e1da5c9db4510e1b32fe51d74329815bb` — `feat(tui): hand off source git workspace`  
**Merged-main CI:** [run 39878431912](https://github.com/zhongyangchuwu/chezmoi-kit/actions/runs/39878431912) — pass

## Delivered

- `g` launches lazygit inside `cm ui` using the same configured source-directory command path as direct `cm git`.
- Source Git porcelain working-tree changes are visible in the workspace status line; status errors remain visible without blocking file inventory.
- Lazygit return re-inventories original workspace scopes and source Git state through the shared terminal-handoff lifecycle.
- Isolated integrations verify lazygit working directory and source Git summary; workspace tests verify key routing, refresh, footer/help, and context preservation.

## Review

Copilot review was requested automatically but unavailable because the requester had exhausted the review quota. Local review is recorded in `REVIEW.md`; branch and merged-main CI passed.

## Follow-on

Phase 3 mature tool handoffs is complete. Phase 4 begins management and diagnostics: explicit lifecycle operations, attribute changes, contextual menus, and diagnostics with the same inventory and confirmation safety boundaries.
