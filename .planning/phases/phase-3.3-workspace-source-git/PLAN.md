# Plan: Workspace Source Git Handoff

1. Expose an app terminal command for lazygit in the configured source directory and reuse it from direct `cm git`.
2. Add a nonfatal source Git working-tree summary to authoritative workspace inventory.
3. Add `g` to workspace handoff kinds, key routing, footer, and help; reuse existing handoff refresh behavior.
4. Cover source-directory launch, summary/error presentation, return refresh, and context preservation with app/TUI/real-tool integration tests.
5. Update docs/planning, run quality gates, obtain CI, and merge.
