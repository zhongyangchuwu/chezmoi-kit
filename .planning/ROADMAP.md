# Roadmap: cm

## Overview

`v0.1.1` is planned as a safe usability-polish release after the published `v0.1.0`. Its scope is limited to additions and fixes that do not change default command usage, the README mental model, existing command semantics, or mutation behavior.

## Version Policy

- Patch releases (`v0.1.Z`) may add optional commands or flags, improve diagnostics, fix display defects, and update non-runtime project hygiene when defaults and the user mental model stay stable.
- Minor releases (`v0.Y.0`) are reserved for behavior, safety-model, status-model, output-contract, or command-semantics changes that users must consciously absorb.
- Major release planning (`v1.0.0`) is deferred until the public command contracts and safety model are stable enough to commit to long-term compatibility.

## Phases

- [ ] Phase 1: Output Controls
- [ ] Phase 2: Doctor Diagnostics
- [ ] Phase 3: TUI Display Width Polish

## Phase Details

### Phase 1: Output Controls

- Goal: Expose existing semantic report renderers through explicit user flags while preserving current default rendering.
- Depends on: `v0.1.0` semantic report documents and renderers.
- Requirements: OUT-FLAGS-01, OUT-FLAGS-02.
- Success Criteria:
  - `cm status`, `cm diff`, and `cm version` support explicit output selection where their report model already supports it.
  - Color policy can be selected explicitly without changing default TTY and `NO_COLOR` behavior.
  - Invalid output or color values fail with clear CLI errors.
  - README documents optional flags only after implementation is verified.
- Plans: TBD

### Phase 2: Doctor Diagnostics

- Goal: Add a read-only `cm doctor` command that explains environment readiness for existing `cm` workflows.
- Depends on: Phase 1 if doctor output uses the same user-facing output/color controls.
- Requirements: DOCTOR-01.
- Success Criteria:
  - Required tools and state such as `chezmoi`, `git`, and chezmoi source path are checked without mutation.
  - Optional tools such as `lazygit` are reported as warnings when absent, not as universal failure.
  - Diagnostic output distinguishes pass, warning, and failure in a script-friendly way.
  - Existing commands keep their current behavior and error paths unless explicitly invoked through `cm doctor`.
- Plans: TBD

### Phase 3: TUI Display Width Polish

- Goal: Fix display-width truncation defects in the sync TUI without changing the sync workflow or action model.
- Depends on: None; may run after Phase 1 and Phase 2 for release ordering.
- Requirements: TUI-WIDTH-01.
- Success Criteria:
  - File list, status text, and diff panes avoid byte-based truncation that can split UTF-8 runes.
  - Wide-character paths and diff lines render without corrupting layout in normal and narrow terminal sizes.
  - Existing `cm sync` review, pending-action, confirmation, and preflight behavior remains unchanged.
  - Focused TUI tests cover Unicode and narrow-width cases.
- Plans: TBD

## Progress

| Phase | Plans Complete | Status | Completed |
|---|---:|---|---|
| Phase 1: Output Controls | 0 | planned | - |
| Phase 2: Doctor Diagnostics | 0 | planned | - |
| Phase 3: TUI Display Width Polish | 0 | planned | - |
