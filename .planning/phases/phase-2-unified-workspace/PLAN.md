# Plan: Phase 2 Unified File Workspace

## Objective

Implement and verify a persistent, read-only `cm ui [path...]` workbench with authoritative file inventory, stable tree/flat filtering, complete diff navigation, and labeled opt-in content views, while preserving all Phase 1 CLI and sync contracts.

## Scope

- App workspace domain and service boundary.
- Chezmoi inventory/path/type/content adapters and strict parsers.
- Explicit CLI `ui` command and shared TUI startup/model.
- Managed-clean/dirty/script, scoped unmanaged, and ignored entries.
- Tree/flat projections, collapse, category filters, file search, and selection retention.
- Diff/content preview modes, sensitive reveal, full-screen preview, horizontal/vertical scroll, preview search, and hunk navigation.
- Focused unit, real-tool integration, CLI, PTY, regression, docs, review, verification, and capture.

## Tasks

1. Add strict chezmoi managed-all JSON and path-list adapters, scoped recursive unmanaged discovery, ignored inventory, target/source/destination content loaders, and bounded errors.
2. Add app inventory merging and preview construction, including dirty status mapping, type/template/encrypted attributes, script/ignored/unmanaged state, scope normalization, and sensitive withholding.
3. Add `Services.Workspace`, `cm ui [path...]`, and a persistent workspace startup path without changing root/status/sync defaults.
4. Generalize the TUI entry model and startup mode while preserving focused sync mutation semantics.
5. Implement flat/tree projections, collapse, filters, file search, and stable selection restoration.
6. Implement labeled preview kinds, explicit reveal, full-screen mode, horizontal/vertical navigation, text search, repeated match traversal, and diff hunk jumps.
7. Add behavior-focused unit and isolated real-chezmoi integration coverage; keep existing tests meaningful and remove no longer valid implementation assertions.
8. Update public docs, changelog, codebase maps, phase summary/review/verification/capture, and root planning state.
9. Run actual PTY scenarios and final Go/release quality gates; revise until every Phase 2 requirement has observed evidence.

## Acceptance Criteria

- `cm ui` stays open and shows a clean managed file when no destination/target drift exists.
- Bare `cm` still renders status and `cm sync` still exits clean or uses the verified focused reconciliation flow.
- Explicit directory scopes include managed entries and recursively discovered unmanaged files under that scope; no scope means no unmanaged HOME scan.
- Ignored entries returned by chezmoi are visibly labeled and never represented as inspected/clean destination state.
- Tree/flat, category filter, collapse, and path search rebuilds retain the selected absolute target when visible and otherwise choose the nearest previous index deterministically.
- Default preview is authoritative diff. Destination, rendered target, and source views are distinctly labeled.
- Encrypted source and template/encrypted rendered target content remain withheld until an explicit per-target reveal command.
- Long-line tails, every diff hunk, and every preview search match are reachable; full-screen and narrow layouts remain usable.
- `cm ui` cannot add/apply/merge or mutate source/destination in this phase.
- Real-tool and PTY scenarios prove clean persistence, inventory categories, navigation, sensitive reveal, and no-mutation quit behavior.

## Verification

- Focused tests: `go test ./internal/chezmoi`, `go test ./internal/app`, `go test ./internal/cli`, `go test ./internal/tui`.
- Isolated real chezmoi tests using the existing wrapper harness and temporary source/destination/cache/config/state.
- Actual PTY workbench at normal and narrow terminal sizes; check tree/flat/filter/search/collapse, horizontal scrolling, hunk/search jumps, preview switching, reveal, and exit.
- Regression smoke for bare `cm`, `cm status`, `cm diff`, and focused `cm sync`.
- `gofmt`, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go mod tidy -diff`, `go mod verify`, govulncheck, GoReleaser config check, LSP diagnostics, YAML diagnostics, and `git diff --check`.

## Risks

- Inventory may require several chezmoi metadata subprocesses; avoid content-wide eager loading and measure real startup before adding concurrency.
- `ignored` cannot explain rule provenance; label returned facts only.
- Recursive unmanaged discovery can expand large trees; enforce explicit scopes, deduplicate paths, and cap entries/queries with an actionable error.
- Content previews can expose secrets. Keep diff baseline unchanged, but require explicit view selection and an additional reveal for known sensitive target/source content.
- Generalizing the TUI can regress focused sync. Keep mode branches small and run the existing sync state-machine and PTY scenarios unchanged.
- The working tree already contains verified Phase 1 changes. A temporary baseline archive exists outside the repository for recovery; do not discard unrelated or pre-existing work.
