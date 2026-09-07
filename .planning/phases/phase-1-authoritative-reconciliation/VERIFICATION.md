# Verification: Phase 1 Authoritative Reconciliation

## Claims Checked

1. Diff and sync previews use authoritative chezmoi target semantics for content, metadata, symlinks, directories, removes, and templates.
2. Reconciliation work follows the second status column and scripts are separated from ordinary sync.
3. Only actions valid for the reviewed target type/template state are offered and accepted.
4. Pending actions are invalidated when reviewed state changes before execution.
5. Successful commands remove targets only after a clean postflight result.
6. Successful buffered output/warnings remain visible.
7. Relative edit targets honor custom chezmoi destination configuration.
8. Encrypted re-add preserves encrypted source identity and updated content.
9. Obsolete custom diff code/dependency is fully removed.
10. Tests, static analysis, module state, vulnerability scan, release config, docs, and planning state are coherent.

## Evidence Observed

### Real CLI/TUI scenarios

- Isolated `cm status --output plain` showed ordinary local targets, a separate pending-script `automation:` block, and source git status.
- Isolated targetless `cm diff --output plain` showed:
  - symlink link-target changes without referent contents;
  - directory permission changes;
  - regular content changes in target-to-destination direction;
  - no pending script diff.
- Script-target `cm sync` exited read-only with `no file changes to reconcile; 1 script pending` and directed the user to chezmoi.
- Real PTY apply flow: review loaded, `p` selected, confirmation opened, `y` executed, process exited 0, completion reported `executed 1, skipped 0, deferred 0`, and destination content matched target.
- Real PTY stale flow: destination changed after confirmation; execution reported `deferred 1`, refreshed the new diff, and left changed destination content untouched.
- Real PTY no-op flow: a wrapper returned successful `re-add` without changing source; postflight retained the target and displayed both the simulated warning and `target still differs after execution`.

### Real chezmoi integration tests

`go test ./internal/app` passed isolated cases for:

- target-to-destination text direction;
- executable-bit-only drift;
- symlink target comparison and `TargetSymlink` metadata;
- directory mode drift and `TargetDirectory` metadata;
- remove entry and apply-only action matrix;
- pending script separation;
- rendered template content and template membership;
- custom destination edit with a noninteractive editor;
- builtin-age encrypted re-add preserving the `encrypted_` source entry and applied updated plaintext.

### Final automated gates

| Check | Observed result |
|---|---|
| `gofmt` on changed Go packages | completed with no output |
| `go test ./...` | passed all packages |
| `go test -race ./...` | passed all packages |
| `go vet ./...` | passed with no output |
| `go mod tidy -diff` | passed with no output |
| `go mod verify` | `all modules verified` |
| `govulncheck ./...` | `No vulnerabilities found.` |
| GoReleaser v2 check | one configuration validated |
| gopls workspace diagnostics | no issues |
| CI/release YAML diagnostics | OK |
| official chezmoi release asset | downloaded; reported v2.72.1; app tests passed with it first on PATH |
| obsolete code search | no runtime/docs references to `internal/diff`, `ContentLoader`, `DiffBytes`, or `go-internal` |

## Coverage

- AUTH-DIFF-01: authoritative forced builtin diff and real target-type scenarios.
- AUTH-STATUS-01: second-column partitioning and docs.
- AUTH-TYPE-01: app-owned types, integration metadata, action matrix, dynamic help.
- AUTH-SCRIPT-01: status report, script-only sync, targetless diff exclusion.
- AUTH-REVIEW-01: unit and real PTY stale-review deferral.
- AUTH-VERIFY-01: unit and real PTY successful-no-op postflight retention.
- AUTH-OUTPUT-01: service/TUI tests and real warning visibility.
- EDIT-DEST-01: service tests and real custom-destination editor integration.
- INTEGRATION-01: isolated real chezmoi suite plus downloaded pinned binary run.
- DOC-02: public docs, changelog, workflows, phase records, and codebase maps updated.

## Failed Checks

- Initial `managed --recursive=false` integration failed because the flag is unsupported; fixed by removing it and avoiding template scans for non-file types.
- Initial age smoke used the wrong key-generation syntax; fixed to `age-keygen --output`.
- Initial `go install github.com/twpayne/chezmoi/v2@v2.72.1` attempts failed: proxy EOF, then module exclude-directive rejection. Replaced with official release-asset download and validated the resulting binary.
- Initial asynchronous CLI sync test timed out because input arrived before review load. Replaced timing dependence with a condition-gated reader.
- Initial remove metadata test found that dump returns no entry for deletion effects. Review now derives remove type from status and the test passes.
- All corresponding final checks passed.

## Skipped Checks

- Remote GitHub Actions execution was not available before push.
- Tagged release publication was not requested.
- Real interactive merge tool execution was not run because it depends on user merge configuration and can mutate source state; terminal handoff and postflight paths remain covered by state-machine tests.

## Untested Claims

- Windows runtime integration; release cross-build remains covered by GoReleaser configuration rather than real chezmoi execution.
- Atomic protection against a filesystem change occurring after final preflight and before chezmoi reads the target; current contract is explicitly optimistic.

## Gaps

- No phase-blocking gaps. Remote CI observation is a ship input, not an implementation-verification blocker.

## Result

passed
