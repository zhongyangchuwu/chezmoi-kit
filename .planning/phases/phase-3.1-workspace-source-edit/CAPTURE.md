# Capture: Workspace Source Edit

**Merged:** 2026-09-10  
**Pull request:** [#6](https://github.com/zhongyangchuwu/chezmoi-kit/pull/6)  
**Squash commit:** `116504866785be37744b940be79cdefebacac5e7` — `feat(tui): edit workspace source safely`  
**Merged-main CI:** [run 34396345865](https://github.com/zhongyangchuwu/chezmoi-kit/actions/runs/34396345865) — pass

## Delivered

- `cm ui` offers contextual `e` source edit for eligible managed file/symlink entries.
- Chezmoi invocation explicitly disables inherited `edit.apply` and `edit.watch` behavior.
- Bubble Tea yields the terminal, then refreshes original inventory scopes after success, cancellation, or nonzero editor return.
- Refresh restores viable workspace context, clears sensitive/preview/search/scroll state, and rejects old async preview completions with a refresh epoch.
- Template and encrypted source edits remain chezmoi-owned and were covered with real isolated integration state; no destination is auto-applied.
- Documentation defines `cm ui` as inspection-first rather than absolutely read-only.

## Review

Copilot review was requested automatically but unavailable because the requester had exhausted the review quota. The implementation received a local code-quality review recorded in `REVIEW.md`; CI and merged-main CI both passed.

## Follow-on

Phase 3 remains in progress. The next plan covers destination-local edit, merge, and lazygit handoffs, reusing the source-edit terminal ownership and refresh boundary without adding workspace apply.
