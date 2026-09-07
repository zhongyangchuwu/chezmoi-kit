# Capture: Directory Browser

## Durable Records Updated

- `CHANGELOG.md`
- `docs/usage.md`
- `docs/design.md`
- `docs/development.md`
- `.planning/ROADMAP.md`
- `.planning/REQUIREMENTS.md`
- `.planning/STATE.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/CONCERNS.md`
- `.planning/codebase/STRUCTURE.md`
- This phase's `CONTEXT.md`, `PLAN.md`, `SUMMARY.md`, `REVIEW.md`, `VERIFICATION.md`, and `CAPTURE.md`

## Durable Decisions

- Direct-child browsing is the only workspace directory-navigation model; tree/flat/collapse modes are removed rather than maintained in parallel.
- Virtual directories represent inventory-path context only. They never reconstruct chezmoi authority or request sensitive file content.
- Contextual keyboard behavior is explicit: Current `h/l` means parent/child traversal; Preview `h/l` means horizontal scroll.

## Ship Input

- Branch: `feat/tui-usability`
- Local verification: passed.
- Before merge: review staged diff, commit, push, PR, and remote CI.
