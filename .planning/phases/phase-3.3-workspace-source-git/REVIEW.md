# Review: Workspace Source Git Handoff

**Date:** 2026-10-23  
**Result:** Approved for PR.

## Reviewed Surface

- Shared lazygit terminal command construction and direct `cm git` compatibility in `internal/app/services.go`.
- Nonfatal source Git summary in workspace inventory.
- Workspace `g` routing, handoff outcome, footer/help, and status-line presentation.
- Isolated source-directory, Git-summary, and TUI refresh coverage.

## Findings Resolved

1. **Direct `cm git` and workspace `g` could diverge.** Both now use the same source-directory command builder; direct CLI injects its streams, while Bubble Tea injects terminal streams for workspace handoff.
2. **A source Git status failure could make file browsing unavailable.** Inventory captures the failure in `SourceStatusError`; the workspace status line exposes it without discarding the authoritative file inventory.
3. **Footer allocation capacity lagged the fourth workspace action.** Action-footer capacity covers `e`, `E`, `m`, and `g` without growing its slice.

## Residual Risks

- Workspace performs one `git status --porcelain=v1` per inventory refresh. This is bounded, necessary for the visible summary, and source-repository failures remain nonfatal.
- Lazygit can change membership, branches, or source paths. The post-handoff inventory and existing deterministic selection fallback, rather than Git heuristics, determine the resulting workspace state.
