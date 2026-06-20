# Verification: 7-semantic-reports

## Claims Checked

1. App no longer hand-builds user-facing status/version strings without semantic tokens.
2. Status output content and styling are generated through the report model.
3. Diff line classification is defined once and reused by TUI/report rendering.
4. Plain output remains readable and suitable for non-TTY capture.
5. ANSI output uses semantic palette rules rather than hard-coded styling at call sites.
6. Color is disabled when `NO_COLOR` is set and not empty.
7. Markdown output exists and is deterministic.
8. Existing command behavior is preserved.

## Evidence Observed

- `internal/report` defines `Document`, semantic block/inline roles, diff line roles, `Plain`, `ANSI`, `Markdown`, and `ClassifyDiffLine`.
- `internal/app/services.go` exposes `StatusReport` and `DiffReport` returning `report.Document`.
- `internal/app/version.go` exposes `Info.Report` returning `report.Document`.
- `internal/cli/render.go` renders reports with `report.ANSI` using auto TTY detection.
- `internal/tui/diff_view.go` styles diff lines using `report.ClassifyDiffLine`.
- `internal/report/report_test.go` covers plain rendering, ANSI palette output, `NO_COLOR`, Markdown output, diff line classification, and diff block newline preservation.
- `internal/app/services_test.go` covers semantic status report construction and diff report target expansion.
- `internal/app/version_test.go` covers semantic version report output.
- Searches found no `StatusOutput`, `FormatDetailed`, `github.com/fatih/color`, stale raw-byte output docs, or stale Phase 7 TODO docs.
- `go test ./internal/report ./internal/app ./internal/cli ./internal/tui` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go mod tidy -diff` produced no diff after applying tidy.
- Go workspace diagnostics reported no issues.
- `go build -o /tmp/cm-report-check ./cmd/cm && /tmp/cm-report-check version` built and printed version report output.
- `NO_COLOR=1 /tmp/cm-report-check version` printed version output without ANSI escapes.
- `/tmp/cm-report-check status` and `/tmp/cm-report-check diff` ran successfully in the local environment.

## Coverage

- OUT-01 covered by semantic documents for status, diff, and version outputs plus plain/ANSI/Markdown renderers.
- OUT-02 covered by role-based ANSI palette, TTY auto-color policy, and `NO_COLOR` tests.
- Phase 7 success criteria 1-5 all have observed evidence.

## Gaps

- No explicit `--output markdown` or `--color=auto|always|never` CLI flags were added; internal renderer support exists for future command flags.
- Status and diff smoke output depends on local chezmoi state; tests provide deterministic behavior coverage.
- Real-terminal release gates from earlier phases remain external.

## Result

Pass. Phase 7 Semantic Reports and Palette Rendering is verified against its success criteria and mapped requirements.
