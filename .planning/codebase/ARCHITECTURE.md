# Architecture

**Mapped:** 2026-09-07

## System Boundary

`cm` is a review/orchestration layer around the `chezmoi` executable. It does not own source naming, target calculation, templates, encryption, ignore rules, diff mechanics, or mutations.

```text
Cobra CLI / Bubble Tea TUI
            |
            v
      internal/app
      /     |      \
chezmoi   report   process
adapter   model    runner
            |
            v
         output
```

## Dependency Direction

- `cmd/cm` depends on `internal/cli`.
- `internal/cli` depends on app contracts, report renderers, and TUI startup.
- `internal/tui` depends only on app-owned sync/domain values plus Bubble Tea presentation libraries.
- `internal/app` depends on chezmoi, process, and report capabilities.
- `internal/chezmoi` depends on `internal/process`.
- `internal/report` and `internal/process` are leaf capabilities.

No runtime package depends on planning or docs.

## App Domain

`internal/app/reconcile.go` and `internal/app/workspace.go` own the values
crossing the TUI boundary:

- `SyncStatus` separates ordinary destination/target entries from scripts.
- `ReconcileEntry` carries raw status code and absolute target path.
- `Review` carries target type, template state, authoritative diff, dirty state, and SHA-256 fingerprint.
- `Action` binds target/action kind to the reviewed fingerprint.
- `ActionResult` preserves successful buffered stdout/stderr.
- `WorkspaceSnapshot` holds a sorted absolute-path inventory and explicit discovery scopes.
- `WorkspaceEntry` carries state, target type, source mapping, template/encryption flags, and status code.
- `WorkspacePreview` separates bounded content, notices, withheld state, and an inspected review.

This keeps `chezmoi.StatusEntry` inside infrastructure/app composition instead of leaking it into the TUI.

## Status Flow

1. App calls `chezmoi status --include=all --exclude=none --path-style=absolute`.
2. Chezmoi adapter strictly parses `XY path` lines.
3. App interprets the second status column:
   - space: destination already matches target; no reconciliation entry;
   - `R`: pending script in `SyncStatus.Scripts`;
   - other effect: ordinary `SyncStatus.Entries` item.
4. Status report independently loads source repository git status.
5. Semantic report renderers produce plain, ANSI, or Markdown output.

## Review Flow

1. TUI selects an absolute target from `SyncStatus.Entries`.
2. App rechecks exact target status.
3. App captures bounded authoritative diff from:

   ```text
   chezmoi --color=false --no-pager --use-builtin-diff \
     diff --include=all --exclude=none --reverse --script-contents=true <target>
   ```

4. A second-column delete effect becomes `TargetRemove` without requiring a dump entry.
5. Other targets load bounded single-entry metadata through `chezmoi dump`.
6. File/symlink targets query `managed --include=templates` for template membership.
7. App fingerprints path, status code, target type, template flag, and diff.
8. TUI caches the review and advertises only allowed actions.

## Action Matrix

- regular file: add, apply, merge;
- template file: apply, merge;
- symlink: apply;
- directory: apply;
- remove: apply;
- script/unknown: no ordinary sync action.

The matrix is conservative because `re-add` ignores non-files and refuses to overwrite templates.

## Execution Flow

1. User selects actions; each stores the current review fingerprint.
2. Confirmation starts sequential execution.
3. App recomputes review immediately before mutation.
4. Clean target: skip and remove.
5. Fingerprint/type mismatch: defer without mutation and cache refreshed review.
6. Non-interactive add/apply captures output; merge receives terminal control.
7. After successful command, app recomputes target status/review.
8. Clean target: remove.
9. Still dirty: retain, clear pending action, show output/reason, and return to review.
10. Completion occurs only when ordinary entries are empty; pending script count remains explicit.

## Workspace Flow

1. CLI calls `WorkspaceService.Inventory` from `cm ui [path...]`.
2. App loads the destination root, normalized scopes, managed path mappings, secret-skipping status, typed managed membership, and source-ignored entries from chezmoi.
3. Explicit scopes trigger bounded recursive candidate discovery: filesystem enumeration only provides child paths; chezmoi `unmanaged` remains the membership authority.
4. TUI derives direct-child browser views and virtual path ancestors from the inventory. State filters retain a directory when any descendant matches; global search locates a selected node's parent directory.
5. Wide layouts render Parent/Current/Preview, medium layouts omit Parent, and narrow layouts switch Current/Preview by focus. Current-focus `h`/`l` traverses directories; preview-focus `h`/`l` scrolls horizontally.
6. File previews load lazily as authoritative diff, destination, rendered target, or source views. Directory previews are local matching-child summaries; loading, error, clean, and withheld feedback have distinct presentation.
7. Template/encrypted rendered targets, encrypted source, and secret-skipped diffs remain withheld until per-target explicit reveal.
8. `?` opens a keyboard-isolated quick-start/key/legend overlay; `Esc`, `?`, and `q` close it without changing workspace data.
9. `e` is available only for eligible managed files/symlinks with a source mapping. It yields the terminal to `chezmoi edit --apply=false --watch=false <absolute-target>`, then re-inventories original scopes regardless of editor exit. Successful refresh preserves viable workspace context, invalidates stale preview epochs/caches, and reloads an authoritative preview; failed refresh keeps the old inventory browsable without a stale preview.

## Edit Flow

- Completion comes from NUL-delimited `chezmoi managed` output.
- Absolute targets pass through.
- Relative targets join to `chezmoi target-path`, not OS home.
- Chezmoi owns editor, decryption/re-encryption, template syntax checking, and source replacement.

## Output Architecture

App services construct semantic `report.Document` values. Roles remain separate from presentation. CLI selects plain/ANSI/Markdown. TUI reuses shared unified-diff classification.

## Process Architecture

`process.Runner` has captured `Output` and streaming/buffered `Run` paths. It normalizes missing executables. Sync add/apply use nil stdin and captured stdout/stderr; merge is terminal-bound.

## Verification Architecture

- Unit/service/TUI/CLI tests use handwritten fakes.
- `services_integration_test.go` uses a real chezmoi executable behind an isolated wrapper with temporary source, destination, cache, config, and persistent state.
- CI and release workflows download pinned chezmoi v2.72.1 before GoReleaser/test execution.
