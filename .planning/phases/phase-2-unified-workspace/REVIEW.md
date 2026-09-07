# Review: Phase 2 Unified File Workspace

## Scope Reviewed

- Workspace CLI wiring, app domain/service, chezmoi inventory/content adapter, shared Bubble Tea model, projection and preview state, tests, user docs, and codebase records.

## Findings

1. **High — secret-skipped templates could be mistaken for clean.** `--skip-secrets` omits secret-backed template status. Fixed by representing missing template status as `uninspected` and withholding its authoritative diff until explicit reveal.
2. **High — native unmanaged output can hide nested unmanaged files behind a directory.** Fixed by bounded recursive chezmoi unmanaged queries over immediate child paths.
3. **Medium — stale asynchronous preview results could reset active preview state.** Fixed by retaining stale cache data while leaving the active message, matches, and scroll position unchanged.
4. **Medium — path filters were not visible after committing search.** Fixed by retaining the committed path filter in the workspace status line.

## Fixes Applied

- Added app-owned workspace contracts and a single concrete service implementation backed by the existing chezmoi client.
- Reused the sync model/rendering path with explicit workspace branches instead of introducing a parallel TUI.
- Kept content lazy and bounded; required explicit reveal for sensitive target/source/diff cases.
- Added behavior-focused TUI tests for projection fallback, withheld reveal, stale preview isolation, navigation, search, and narrow/full layouts.

## Waivers

- Workspace has no editor, merge, Git, watcher, ignore parser, or mutation behavior. These are deliberate Phase 3/4 scope exclusions.
- Startup uses several sequential metadata commands. Correct chezmoi classification and secret safety take precedence over speculative parallelism.

## Remaining Risks

- Explicit reveal can expose sensitive content in terminal scrollback.
- Large scoped unmanaged trees can hit the documented query/entry limits.
- Remote CI observation, review/PR, and release publication remain separate ship actions.
