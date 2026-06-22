# Verification: Phase 1 Output Controls

## Claims Checked

- `cm status`, `cm diff`, and `cm version` accept explicit report output selection.
- Explicit color policy can force ANSI, disable ANSI, or keep auto behavior.
- Default non-TTY output remains plain text through the existing ANSI renderer with `ColorAuto` and non-TTY detection.
- Invalid output or color values fail with clear CLI errors.
- README documents only implemented optional behavior.

## Evidence Observed

- `go test ./internal/cli ./internal/report` passed.
- `go test ./...` passed.
- `go run ./cmd/cm version --output markdown` printed Markdown:
  - `**cm:** **dev**`
  - `commit: `unknown``
- `go run ./cmd/cm version --color always` completed successfully.
- `go run ./cmd/cm status --output json` returned exit status 1 with `unsupported output "json"`.
- Tests added:
  - `TestRunReportOutputFlagSelectsMarkdown`
  - `TestRunReportColorFlagRequestsANSIOnNonTTY`
  - `TestRunReportFlagsRejectInvalidValues`
  - `TestNoColorDisablesAutoANSI`
  - `TestColorAlwaysOverridesNoColor`

## Coverage

- CLI flag parsing and command wiring: covered by `internal/cli` tests.
- Report renderer policy: covered by `internal/report` tests.
- User-visible smoke: covered for Markdown output and invalid output error.
- Regression: `go test ./...` passed across the repository.

## Gaps

- No interactive TTY smoke was run for visual ANSI rendering.
- `cm sync` TUI is intentionally out of scope for output flags.
- `cm completion` inherits Cobra persistent flags but its generated shell script is not reformatted by report renderers.

## Result

Passed with noted TTY visual-smoke gap.
