# Verification: Phase 2 Unified File Workspace

## Claims Checked

1. `cm ui [path...]` persists clean managed entries without changing bare `cm` or focused `cm sync` behavior.
2. Inventory combines managed clean/dirty/script, source-ignored, and explicitly scoped recursive unmanaged entries through chezmoi authority.
3. Tree/flat/filter/search/collapse preserve selection or use deterministic fallback.
4. Bounded labeled diff, destination, rendered target, and source previews support full view, horizontal/vertical movement, hunk jumps, and text-match traversal.
5. Secret-skipped templates and encrypted state are uninspected/withheld until explicit reveal.
6. Workspace has no mutation path; quitting changes neither source nor destination.
7. Tests, static analysis, module state, vulnerability checks, release checks, diagnostics, docs, and planning records are coherent.

## Evidence Observed

### Real chezmoi and PTY scenarios

- Isolated real chezmoi integration passed clean, dirty, template, encrypted, symlink, managed directory, pending script, source-ignored entry, explicit scoped unmanaged file, nested unmanaged file, and all preview contracts.
- Secret-backed template integration proved inventory and withheld diff did not invoke its configured provider; explicit diff reveal invoked the provider and returned an inspected dirty result.
- Actual normal PTY workspace showed clean, dirty, uninspected encrypted, ignored, symlink, directory, script, and scoped unmanaged rows. It exercised tree/filter/path search, authoritative dirty diff, long-line horizontal tail, both diff hunks, full preview, preview search, template source, target withholding/reveal, and encrypted source withholding.
- Actual narrow/full-preview rendering was covered by TUI state tests and manual PTY navigation.
- The same isolated source/destination archive hash before and after `q` was `f4c92711ff0f5bf1e2e50969f52bd5d083854777d24697d841ec8a91f81bb804`.

### Final automated gates

| Check | Observed result |
|---|---|
| `gofmt -w internal/app/*.go internal/chezmoi/*.go internal/cli/*.go internal/tui/*.go` | completed with no output |
| `go test ./...` | passed all packages |
| `go test -race ./...` | passed all packages |
| `go vet ./...` | passed with no output |
| `go mod tidy -diff && go mod verify` | passed; `all modules verified` |
| `govulncheck ./...` | `No vulnerabilities found.` |
| GoReleaser v2 `check` | one configuration validated |
| GoReleaser snapshot | passed; built six Linux/macOS/Windows amd64/arm64 artifacts |
| gopls workspace diagnostics | no issues found |
| CI workflow YAML diagnostics | OK |
| `git diff --check` | passed with no output |

## Coverage

- WORKSPACE-01: CLI wiring, persistent clean-row TUI test, real PTY workspace, and existing bare/sync CLI regressions.
- WORKSPACE-02: app integration inventory matrix and scoped recursive unmanaged real fixture.
- WORKSPACE-03: TUI projection/collapse/filter/path-search selection tests.
- WORKSPACE-04: TUI navigation/full/narrow tests plus actual PTY long-tail/hunk/search interactions.
- TEMPLATE-01: app integration target/source/reveal scenarios, secret-provider test, and PTY template target reveal.

## Gaps

- Windows runtime chezmoi integration is not exercised; snapshot cross-build validates compilation only.
- Remote CI, PR review, merge, and release publication were not requested.

## Result

passed
