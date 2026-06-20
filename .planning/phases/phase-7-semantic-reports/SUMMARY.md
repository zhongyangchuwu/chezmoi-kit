# Summary: 7-semantic-reports

## Completed Changes

- Added `internal/report` with semantic `Document`, block, inline role, diff line, plain renderer, ANSI renderer, Markdown renderer, and color policy.
- Added shared `report.ClassifyDiffLine` and moved TUI diff styling to use the same classifier as report rendering.
- Converted app status output to semantic `StatusReport(targets) (report.Document, error)`.
- Converted app diff output to semantic `DiffReport(targets) (report.Document, error)`.
- Converted version output to `Info.Report(name) report.Document`.
- Updated CLI status, diff, and version commands to render reports through `internal/cli/render.go`.
- Implemented ANSI auto-color policy: TTY stdout enables color, non-TTY stays plain, and non-empty `NO_COLOR` disables ANSI.
- Implemented deterministic Markdown rendering in `internal/report` without using Markdown as the internal model.
- Removed the unused `github.com/fatih/color` dependency and tidy-cleaned `go.mod`/`go.sum`.
- Updated public and planning docs to describe semantic reports, color policy, and shared diff classification.

## Files Changed

- `internal/report/report.go`
- `internal/report/diff.go`
- `internal/report/render.go`
- `internal/report/report_test.go`
- `internal/app/services.go`
- `internal/app/services_test.go`
- `internal/app/version.go`
- `internal/app/version_test.go`
- `internal/cli/render.go`
- `internal/cli/status.go`
- `internal/cli/diff_cmd.go`
- `internal/cli/root.go`
- `internal/cli/cli_test.go`
- `internal/tui/diff_view.go`
- `go.mod`
- `go.sum`
- `docs/design.md`
- `docs/development.md`
- `.planning/codebase/ARCHITECTURE.md`
- `.planning/codebase/CONVENTIONS.md`
- `.planning/codebase/STACK.md`
- `.planning/codebase/CONCERNS.md`
- `.planning/codebase/MAP.md`

## Deviations

- Markdown renderer was implemented now because it was small and deterministic; no CLI `--output markdown` flag was added because output-mode selection was out of scope.
- The report model intentionally remains minimal: headings, paragraphs, blank lines, code blocks, diff blocks, inline roles, and diff line roles only.

## Evidence

- `go test ./internal/report ./internal/app ./internal/cli ./internal/tui` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go mod tidy -diff` produced no diff after applying tidy.
- Go workspace diagnostics reported no issues.
- `go build -o /tmp/cm-report-check ./cmd/cm && /tmp/cm-report-check version` built and printed version report output.
- `NO_COLOR=1 /tmp/cm-report-check version` printed version output without ANSI escapes.
- `/tmp/cm-report-check status` and `/tmp/cm-report-check diff` ran successfully in the local environment.
- Searches found no `StatusOutput`, `FormatDetailed`, `github.com/fatih/color`, stale raw-byte output docs, or stale Phase 7 TODO docs.

## Unresolved Risks

- No user-facing `--output` or `--color` flags exist yet; report renderers have internal support but CLI policy remains auto ANSI only.
- TUI file/diff layout still truncates by byte length rather than display width; this predates Phase 7 and remains a separate UI rendering concern.
- Existing external release gates remain: observed remote CI, explicit tag approval, and real-terminal smoke where available.
