# Plan: Phase 3 TUI Display Width Polish

## Objective

Replace byte-count truncation in sync TUI views with display-width-aware clipping while preserving existing sync behavior.

## Scope

In scope:

- TUI file list rendering.
- TUI status/help line truncation if it uses the same helper.
- TUI diff pane line truncation.
- Focused tests for Unicode, wide characters, and narrow terminal widths.

Out of scope:

- Sync action behavior.
- Diff generation.
- Report renderer behavior outside TUI.
- Terminal resize feature changes.

## Tasks

1. Inspect current TUI view truncation helpers and dependency support for display width.
2. Implement display-width-aware truncation that preserves valid UTF-8.
3. Apply the helper to file names, status text, and diff lines.
4. Add focused tests for ASCII parity, wide-character truncation, and narrow widths.
5. Run TUI and repository tests; record verification.

## Acceptance Criteria

- Byte-based truncation no longer splits UTF-8 runes.
- Wide-character paths and diff lines fit within requested display widths.
- Existing sync review, pending-action, confirmation, and preflight tests still pass.
- No new broad layout abstraction is added.

## Verification

- Run `go test ./internal/tui`.
- Run `go test ./...`.

## Risks

- Display width for combining marks and emoji sequences can be subtle; tests should cover representative wide CJK paths at minimum.
- Narrow widths can make ellipsis behavior ambiguous; choose predictable clipping over elaborate layout.
