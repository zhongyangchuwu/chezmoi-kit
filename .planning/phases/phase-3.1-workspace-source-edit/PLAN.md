# Plan: Workspace Source Edit

1. Add an app-owned source-edit eligibility rule and a workspace terminal command that delegates to chezmoi with `--apply=false --watch=false`.
2. Reuse the existing terminal command abstraction but add workspace-specific handoff messages; do not share sync execution failure semantics that quit the TUI.
3. On editor return, inventory the original scopes. Restore surviving directory/filter/focus/preview state, clear preview/reveal/match/scroll caches, and reload the selected authoritative preview.
4. Bind `e` only outside search/help/handoff states. Show actionable unavailable reasons and contextual footer/help entries.
5. Add behavior-focused app/TUI tests, actual editor PTY round-trip, Go checks, docs/planning records, review, PR, CI, and squash merge.

## Non-Goals

- No local destination editor (`E`), merge (`m`), lazygit (`g`), workspace apply, source Git mutation, modal action menu, or automatic mutation after editor exit.
