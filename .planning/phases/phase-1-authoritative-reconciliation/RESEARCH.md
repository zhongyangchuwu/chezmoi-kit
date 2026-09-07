# Research: Phase 1 Authoritative Reconciliation

## Question

How should `cm` preview, classify, execute, and verify chezmoi changes without duplicating target-state semantics or mutating state the user did not review?

## Sources Checked

- Current code: `internal/app/services.go`, `internal/chezmoi/*`, `internal/diff/*`, `internal/tui/*`, `internal/cli/*`.
- Existing planning and release records under `.planning/`.
- Official chezmoi command docs for `status`, `diff`, `dump`, `managed`, `target-path`, `edit`, `merge`, `verify`, `re-add`, templates, scripts, and common entry types.
- Isolated runtime probes against installed chezmoi v2.72.1 with temporary source, destination, cache, and persistent state.

## Findings

- `chezmoi status` uses its first column for destination changes since last write and its second column for apply effects; `R` means a script will run.
- Current `cm` treats every status row as a regular file and offers the same actions for every row.
- Current local reads use `os.Open`, which dereferences symlinks and cannot represent directory or mode-only state.
- Isolated current-code probes produced `no diff` for executable-bit drift, compared symlink referent contents instead of link targets, failed on a directory target, and rendered a pending script as file deletion even though apply executes it.
- Forced chezmoi builtin diff represented content, file mode, directory mode, symlink target, and script changes correctly. `--reverse` preserved cm's existing target-to-local direction.
- `chezmoi dump --format=json --recursive=false <target>` returns typed target state. `managed --include=templates` identifies template-backed targets without source filename parsing.
- `chezmoi re-add` preserves encrypted attributes, ignores non-files, and refuses to overwrite templates; action gating and postflight verification are therefore required.
- Current preflight checks only whether a target remains dirty. It does not prove that the diff is unchanged since review.
- Current execution removes a target after command success without proving the target is clean, and successful buffered stdout/stderr is discarded.
- `cm edit` joins relative completion values to OS home even when chezmoi uses a different destination directory.

## Tradeoffs

- Keep custom byte diff and add special cases: smaller immediate change but continues duplicating target semantics and invites more type-specific defects. Rejected.
- Use configured `chezmoi diff`: correct semantics but may launch user diff commands or pagers inside the TUI. Rejected.
- Force chezmoi builtin diff and parse/render its text: preserves authority, avoids external tools, and keeps the current semantic renderer. Selected.
- Eagerly dump all dirty target contents: fewer subprocesses but higher memory and secret exposure. Rejected.
- Lazily load bounded selected-target metadata and cache it for the session: more subprocesses but safer and aligned with current lazy diff behavior. Selected.
- Execute scripts as an apply action: convenient but semantically different, potentially non-idempotent, and incompatible with a clean-after-sync invariant. Rejected for this phase.

## Confidence

High. Findings are supported by current source, official command contracts, and isolated real-tool reproduction of each load-bearing defect.
