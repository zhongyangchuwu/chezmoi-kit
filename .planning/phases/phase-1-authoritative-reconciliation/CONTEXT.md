# Context: Phase 1 Authoritative Reconciliation

## Goal

Make every `cm diff` and `cm sync` decision reflect chezmoi's actual target-state semantics and the exact state reviewed by the user.

## Constraints

- Keep chezmoi authoritative; force its builtin diff and disable pagers rather than invoking configured external diff tools.
- Preserve the established diff direction: chezmoi target is the old side and local destination is the new side.
- Keep review data bounded and never log rendered contents or subprocess output.
- Scripts must not enter ordinary file reconciliation execution.
- Use conservative action gating: regular files may add/apply/merge; template files may apply/merge; symlinks, directories, and removes may apply; unknown types are read-only.
- Do not add template source views, template data commands, horizontal scrolling, JSON output, package managers, or cancellation in this phase.
- Remove obsolete custom diff code after all callers migrate.

## Decisions

- Add app-owned reconciliation entry, review, target-type, action-result, and sync-status values.
- Derive reconciliation work from the second chezmoi status column; first-column-only history drift is already target-matched.
- Partition raw chezmoi status rows into file reconciliation entries and scripts.
- Load selected-target metadata from bounded `chezmoi dump --format=json --recursive=false` output and detect templates with `chezmoi managed --include=templates`.
- Generate previews with `chezmoi --color=false --no-pager --use-builtin-diff diff --include=all --exclude=none --reverse`.
- Hash the reviewed status/type/template/diff state and include the fingerprint in pending actions.
- Recompute review state before mutation; stale actions are deferred and refreshed.
- Re-check status after successful execution; unresolved targets remain visible.
- Resolve relative edit targets with `chezmoi target-path`.

## Open Questions

- None blocking. Unsupported or ambiguous target types remain read-only until real chezmoi evidence justifies additional actions.

## Verification Expectations

- Focused unit regressions must cover parsing, action gating, stale review, postflight retention, warnings, script separation, and destination resolution.
- Real chezmoi integration checks must exercise content changes, metadata-only changes, symlinks, directories, scripts, templates, encrypted re-add, and custom destinations where practical.
- User-visible CLI/TUI smoke must prove script messaging and authoritative preview behavior.
- Final gates: gofmt, focused tests, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go mod tidy -diff`, and `go mod verify`.
