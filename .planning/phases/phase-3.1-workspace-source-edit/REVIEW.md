# Review: Workspace Source Edit

**Date:** 2026-09-10  
**Result:** Approved for PR.

## Reviewed Surface

- App eligibility and terminal command construction in `internal/app/workspace_edit.go`.
- Workspace terminal lifecycle, original-scope refresh, context restoration, and cache invalidation in `internal/tui/workspace_edit.go`.
- Key routing, input lock, footer/help, status rendering, preview staleness handling, service fakes, and integration coverage.
- User/developer docs and planning traceability.

## Findings Resolved

1. **Stale async preview could overwrite a post-editor refresh.** Preview cache clearing alone was insufficient when the selected target/view survived and an old response arrived after refresh. Added `previewEpoch` to preview requests/responses and reject older epochs; regression coverage proves the stale response is ignored.
2. **Workspace service interface update broke the CLI test fake.** Added its `SourceEditCommand` implementation and reran full test/vet/race gates.
3. **Availability explanations prioritized source-map absence over an actual directory/script type.** Eligibility now gives directory/script reasons first while retaining unmapped unmanaged/ignored explanations.

## Residual Risks

- External editor configuration and behavior remains chezmoi-owned by design. `cm` explicitly disables apply/watch, refreshes after every returned terminal command, and remains usable after construction, editor, or refresh errors.
- The handoff is intentionally a suspend/run/restore path rather than an embedded terminal; resize/cancellation behavior remains Bubble Tea/chezmoi terminal behavior for future handoffs.
