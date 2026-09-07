# Capture: Phase 1 Authoritative Reconciliation

## Durable Docs Updated

- `README.md`
- `docs/usage.md`
- `docs/design.md`
- `docs/development.md`
- `CHANGELOG.md`
- `.planning/codebase/MAP.md`
- `.planning/codebase/STRUCTURE.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/STACK.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/CONCERNS.md`

## Planning Records Updated

- `.planning/PROJECT.md`
- `.planning/REQUIREMENTS.md`
- `.planning/ROADMAP.md`
- `.planning/STATE.md`
- Phase `CONTEXT.md`, `RESEARCH.md`, `PLAN.md`, `SUMMARY.md`, `REVIEW.md`, and `VERIFICATION.md`
- Superseded template-discussion phase files removed; template workspace remains future Phase 2.

## Learnings

- A wrapper around chezmoi must delegate diff semantics, not merely rendered content retrieval; target state includes modes, links, directories, removes, and scripts.
- The second status column is the destination/target reconciliation fact. The first column alone does not justify a sync action.
- Script execution cannot share the ordinary clean-after-file-sync model, especially for always-run scripts.
- Review confirmation needs an optimistic concurrency token and an observed postcondition.
- `chezmoi dump` supplies target types but not remove entries; deletion effect comes from status.
- `managed --include=templates` is a better template fact than parsing source filenames.
- Real-tool integration tests caught unsupported flags and command-shape assumptions that handwritten fakes could not.
- ChezMoi v2.72.1 cannot be installed from its module root with `go install` because of exclude directives; use the official release asset in CI.
- Async terminal tests need condition-based input gates rather than fixed or immediate input streams.

## Ship Inputs

- Branch: `feat/authoritative-reconciliation`
- Phase verification: passed.
- Required before merge/release: review changed paths, create a conventional commit/PR if desired, observe remote CI, then choose `v0.2.0` release timing.
- Release notes source: `CHANGELOG.md` Unreleased section.
- No tag, push, PR, merge, or release was performed in this phase.
