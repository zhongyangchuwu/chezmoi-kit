# Research: Phase 1 Template Operation Model Discussion

## Question
For template-backed targets, what should `cm sync` show and which operations can safely be delegated to chezmoi without `cm` writing or rendering templates itself?

## Sources Checked
- Chezmoi templating docs: `https://www.chezmoi.io/user-guide/templating/`
- Chezmoi machine differences docs: `https://www.chezmoi.io/user-guide/manage-machine-to-machine-differences/`
- Chezmoi command docs: `cat`, `diff`, `status`, `source-path`, `edit`, `merge`, `execute-template`, `add`
- Current `cm` architecture and seams: `.planning/codebase/ARCHITECTURE.md`, `internal/chezmoi/client.go`, `internal/chezmoi/content.go`, `internal/app/services.go`, `internal/tui/*`

## Findings
- Chezmoi treats a source file as a template when its source filename has a `.tmpl` suffix or it is under `.chezmoitemplates/`.
- Template data comes from chezmoi built-ins, `.chezmoidata.$FORMAT`, and the config `data` section; `cm` should not reproduce this loading order.
- `chezmoi cat <target>` writes target contents. For template-backed files, this is the rendered target content. Current `cm diff` already delegates target content loading to `chezmoi cat` and compares that rendered output to the destination file.
- `chezmoi source-path [target...]` prints each target's source state path; with no targets it prints the source directory. This is the correct source path lookup for TUI metadata.
- `chezmoi edit <target>` edits source state. For templates, it preserves `.tmpl` extension so editors can highlight templates, and chezmoi checks template syntax when the editor exits. This is the preferred path for source template edits.
- `chezmoi merge <target>` performs a three-way merge between destination state, source state, and target state. `merge.args` defaults to `{{ .Destination }}`, `{{ .Source }}`, and `{{ .Target }}`. If target state cannot be computed, chezmoi performs a two-way merge instead.
- For template-backed targets, source state can be a `.tmpl` file while destination and target are rendered/non-template contents. TUI copy must not imply that merge is a normal rendered-file-only merge.
- `chezmoi add --template <target>` creates/updates source state as a template. `chezmoi add --autotemplate` exists but chezmoi warns its greedy substitutions can generate unwanted templates.

## Tradeoffs
- Directly edit template files in `cm`: gives control but duplicates chezmoi source path, encryption, syntax-check, and editor behavior. Reject for now.
- Delegate edit/merge to chezmoi from TUI: preserves chezmoi semantics and keeps `cm` as review/orchestration layer. Preferred.
- Show both rendered diff and template source in TUI: improves understanding but risks clutter. Make template/source detail opt-in, not default.
- Add CLI template commands immediately: useful for scripts but can distract from TUI-first UX. Keep CLI compatibility minimal until TUI semantics are chosen.

## Confidence
High for chezmoi command semantics: official docs cover `cat`, `source-path`, `edit`, and `merge`. Medium for TUI layout and keybindings: this needs product choice, not more source research.
