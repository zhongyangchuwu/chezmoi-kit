# Codebase Concerns

**Mapped:** 2026-09-07

## Safety

- Review preflight is optimistic, not an atomic filesystem lock. A very small race remains between final fingerprint comparison and chezmoi reading the target.
- Direct `cm apply` preserves chezmoi semantics and may execute scripts; docs must keep this distinct from `cm sync`.
- New/unknown chezmoi target types are read-only in sync until evidence supports an action matrix.
- A configured merge tool can block, fail, or leave a target dirty; postflight retains it, but terminal behavior remains environment-specific.

## Secrets

- `--skip-secrets` intentionally leaves secret-backed templates uninspected; those entries must never be labeled clean until an explicit authoritative diff succeeds.
- Explicit workspace reveal can place rendered/decrypted content in terminal scrollback. Content remains bounded and is not written to timing logs.
- Source-ignored entries expose only chezmoi's returned fact; workspace does not infer ignore-rule provenance or destination-only ignored state.
- Markdown diff output should not be pasted into public issues without review.

## Subprocesses

- Active mutating subprocesses are intentionally non-cancellable until a context-aware runner/client/app contract is designed.
- Successful buffered output is retained for one target; a stricter output-size contract may be needed if real hooks produce excessive output.
- CI/release real-tool tests depend on downloading a pinned public chezmoi release asset.

## TUI

- Workspace starts several bounded chezmoi metadata commands before lazy preview loading; very large inventories can feel slower, but eager content loading is intentionally avoided.
- Scoped unmanaged recursion caps discovered entries and queries; users must narrow a pathological scope rather than trigger an unbounded home scan.
- Long preview lines have horizontal scrolling and narrow terminals switch to one focused pane; extremely tiny terminals may still have limited ergonomics.
- Cached selected previews are bounded but cumulative for a long browsing session.
- Stale asynchronous preview completions cache their own result but do not alter the active view message, matches, or scroll position.
- Semantic colors are intentionally non-configurable in the current slice. Text markers and `NO_COLOR=1` provide the accessibility fallback.
- Very narrow terminals shorten help/legend prose; the quick help and footer preserve operation routes instead of introducing a paged help state.

## Compatibility

- Status and dump parsers intentionally fail on unexpected non-empty formats; future chezmoi format changes should fail visibly and update pinned integration evidence.
- Source git porcelain parsing is newline-based rather than `-z`; unusual filenames may be quoted or unsupported.
- The module path is `github.com/zhongyangchuwu/cm` while the current repository is named `chezmoi-kit`; public `go install` identity remains unresolved.

## Workflow

- GitHub-hosted CI/release jobs are configured but require remote observation after push/tag.
- Package-manager publishing, signing, notarization, and public repository visibility remain deferred decisions.
