# Capture: Workspace Destination Edit and Merge

**Merged:** 2026-10-23  
**Pull request:** [#9](https://github.com/zhongyangchuwu/chezmoi-kit/pull/9)  
**Squash commit:** `ee51c5e6967dc22ca9237c415ee262062ea685e6` — `feat(tui): hand off destination edit and merge`  
**Merged-main CI:** [run 39823893504](https://github.com/zhongyangchuwu/chezmoi-kit/actions/runs/39823893504) — pass

## Delivered

- `E` opens an eligible managed regular destination file through `VISUAL` or `EDITOR`, supports quoted editor arguments, rejects symlink referents, and never changes source or invokes chezmoi apply/watch behavior.
- `m` opens configured `chezmoi merge` for eligible managed regular files. Workspace help and status identify destination, source, and rendered target; merge can change source or destination, but cm never applies after it returns.
- `e`, `E`, and `m` now share terminal ownership, error recovery, original-scope refresh, context restoration, and stale-preview invalidation.
- Real isolated integrations prove destination edit source preservation and configured merge argument roles without implicit apply.

## Review

Copilot review was requested automatically but unavailable because the requester had exhausted the review quota. Local review is recorded in `REVIEW.md`; branch and merged-main CI passed.

## Follow-on

Phase 3 remains in progress. The remaining handoff is source Git/lazygit with source Git refresh after return. Workspace apply remains out of scope.
