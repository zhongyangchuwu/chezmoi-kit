# Codebase Map

**Mapped:** 2026-06-20

## Documents

| Document | Lines | Summary |
|---|---:|---|
| STRUCTURE.md | 87 | `cmd/cm` is a thin entrypoint; app, CLI, TUI, chezmoi, diff, and process packages have explicit boundaries. |
| ARCHITECTURE.md | 111 | CLI composition, app support, chezmoi adapter, process runner, diff engine, and TUI are separated. |
| CONVENTIONS.md | 116 | Standard Go tests with handwritten fakes, app-owned service contracts, injected readers/writers, and narrow interfaces. |
| CONCERNS.md | 38 | Remaining risks center on ship-time terminal/CI gates, small report model growth, and byte-width TUI truncation. |

## Key Takeaways

- Phase 5 package naming is now the expected architecture language: `internal/cli`, `internal/tui`, `internal/app`, `internal/chezmoi`, `internal/diff`, and `internal/process`.
- The previous build, reconcile, syncdiff, and shared test helper package boundaries are intentionally retired.
- `cm sync` remains the highest-risk behavior because terminal signal handling and execution-state semantics directly affect user files and terminal state.
- Phase 6 moved broad application service ownership and sync contracts into `internal/app`.
- Phase 7 moved user-facing status, diff, and version output into semantic report documents and renderers.
