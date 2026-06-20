# Codebase Map

**Mapped:** 2026-06-20

## Documents

| Document | Lines | Summary |
|---|---:|---|
| STRUCTURE.md | 87 | `cmd/cm` is a thin entrypoint; app, CLI, TUI, chezmoi, diff, and process packages have explicit boundaries. |
| ARCHITECTURE.md | 111 | CLI composition, app support, chezmoi adapter, process runner, diff engine, and TUI are separated. |
| STACK.md | 103 | Go module CLI using Cobra, Bubble Tea, Lip Gloss, fatih/color, and rogpeppe diff. |
| CONVENTIONS.md | 116 | Standard Go tests with handwritten fakes, injected readers/writers, and narrow interfaces. |
| CONCERNS.md | 38 | Remaining risks center on ship-time terminal/CI gates and future Phase 6/7 ownership moves. |

## Key Takeaways

- Phase 5 package naming is now the expected architecture language: `internal/cli`, `internal/tui`, `internal/app`, `internal/chezmoi`, `internal/diff`, and `internal/process`.
- The previous build, reconcile, syncdiff, and shared test helper package boundaries are intentionally retired.
- `cm sync` remains the highest-risk behavior because terminal signal handling and execution-state semantics directly affect user files and terminal state.
- Phase 6 should move broad application service ownership and sync contracts into `internal/app`.
- Phase 7 should move user-facing output content into semantic report documents and renderers.
