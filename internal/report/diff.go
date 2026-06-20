package report

import "strings"

// DiffLineKind identifies semantic unified diff line roles.
type DiffLineKind int

const (
	DiffLineContext DiffLineKind = iota
	DiffLineHeader
	DiffLineHunk
	DiffLineAdd
	DiffLineRemove
	DiffLineMeta
)

// DiffLine is a unified diff line plus its semantic role.
type DiffLine struct {
	Text string
	Kind DiffLineKind
}

func ClassifyDiffLine(line string) DiffLine {
	switch {
	case strings.HasPrefix(line, "@@"):
		return DiffLine{Text: line, Kind: DiffLineHunk}
	case strings.HasPrefix(line, "diff "), strings.HasPrefix(line, "---"), strings.HasPrefix(line, "+++"):
		return DiffLine{Text: line, Kind: DiffLineHeader}
	case strings.HasPrefix(line, "+"):
		return DiffLine{Text: line, Kind: DiffLineAdd}
	case strings.HasPrefix(line, "-"):
		return DiffLine{Text: line, Kind: DiffLineRemove}
	case strings.HasPrefix(line, `\ No newline`):
		return DiffLine{Text: line, Kind: DiffLineMeta}
	default:
		return DiffLine{Text: line, Kind: DiffLineContext}
	}
}

func DiffLineRole(kind DiffLineKind) Role {
	switch kind {
	case DiffLineHeader:
		return RoleDiffHeader
	case DiffLineHunk:
		return RoleDiffHunk
	case DiffLineAdd:
		return RoleDiffAdd
	case DiffLineRemove:
		return RoleDiffRemove
	case DiffLineMeta:
		return RoleDiffMeta
	default:
		return RoleText
	}
}
