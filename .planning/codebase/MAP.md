# Codebase Map

**Mapped:** 2026-09-07

## Documents

| Document | Summary |
|---|---|
| STRUCTURE.md | Current package/file layout after authoritative reconciliation cutover. |
| ARCHITECTURE.md | CLI, app domain, chezmoi adapter, report, process, and TUI data flows. |
| CONVENTIONS.md | Go, error, parser, testing, integration, output, and workflow conventions. |
| STACK.md | Toolchain, direct dependencies, external commands, and release integration. |
| CONCERNS.md | Remaining safety, terminal, secrecy, subprocess, compatibility, and distribution risks. |

## Key Takeaways

- `cm` owns reconciliation review and inspection-first workspace browsing; chezmoi owns target-state and mutation semantics.
- Runtime packages are `internal/cli`, `internal/app`, `internal/chezmoi`, `internal/process`, `internal/report`, and `internal/tui`.
- `cm ui [path...]` is a persistent workbench: unscoped inventory avoids unmanaged destination traversal, explicit scopes receive recursively discovered chezmoi unmanaged candidates, and `e` explicitly edits eligible source state without applying/watching destination changes.
- `internal/app` owns reconciliation entries/reviews plus workspace snapshots, previews, bounded inventory merging, explicit sensitive reveal policy, and source-edit command construction.
- `internal/chezmoi` forces builtin reverse diffs and supplies strict managed path mappings, typed membership, source-ignored paths, scoped unmanaged lists, and content adapters.
- `internal/tui` shares one model/rendering path for sync and workspace, with explicit mode behavior, stable projections, cached previews, stale-completion isolation, and external editor handoff refresh.
- Real chezmoi integration tests run against isolated temporary state and a pinned CI/release binary.
