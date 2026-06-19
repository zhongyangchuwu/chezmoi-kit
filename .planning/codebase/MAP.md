# Codebase Map

**Mapped:** 2026-06-19

## Documents

| Document | Lines | Summary |
|---|---:|---|
| STRUCTURE.md | 62 | `cmd/cm` is a thin entrypoint; all behavior lives under `internal/` packages. |
| ARCHITECTURE.md | 78 | CLI composition, chezmoi adapter, process runner, syncdiff, and TUI are cleanly separated. |
| STACK.md | 88 | Go module CLI using Cobra, Bubble Tea, Lip Gloss, fatih/color, and rogpeppe diff. |
| CONVENTIONS.md | 85 | Standard Go tests with handwritten fakes, injected readers/writers, and narrow interfaces. |
| CONCERNS.md | 77 | v0.1.0 risks center on TUI lifecycle, subprocess behavior, stale docs, and missing release artifacts. |

## Key Takeaways

- The package architecture is already suitable for v0.1.0: `internal/cli` composes services, `internal/ui` consumes `reconcile.ReviewService`, and `internal/process` isolates external commands.
- Release blockers are mostly lifecycle and repository hygiene, not broad architecture rewrites.
- `cm sync` is the highest-risk behavior because terminal signal handling and execution-state semantics directly affect user files and terminal state.
- Documentation must be cleaned before release because `docs/superpowers/` contains stale MVP assumptions that contradict current TUI behavior.
- CI should enforce `go test`, race tests, `go vet`, `go mod tidy -diff`, module verification, vulnerability scanning, and release builds before tagging.
