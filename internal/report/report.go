package report

import "strings"

// Document is a semantic command-output document.
type Document struct {
	Blocks []Block
}

// BlockKind identifies a top-level document block.
type BlockKind int

const (
	BlockParagraph BlockKind = iota
	BlockHeading
	BlockBlank
	BlockCode
	BlockDiff
)

// Role describes the semantic meaning of inline text.
type Role int

const (
	RoleText Role = iota
	RoleStrong
	RoleMuted
	RoleCode
	RolePath
	RoleCommand
	RoleWarning
	RoleStatus
	RoleDiffHeader
	RoleDiffHunk
	RoleDiffAdd
	RoleDiffRemove
	RoleDiffMeta
)

// Inline is semantic inline text.
type Inline struct {
	Text string
	Role Role
}

// Block is one semantic document block.
type Block struct {
	Kind         BlockKind
	Level        int
	Inlines      []Inline
	Language     string
	Code         string
	DiffLines    []DiffLine
	FinalNewline bool
}

func Text(text string) Inline    { return Inline{Text: text, Role: RoleText} }
func Strong(text string) Inline  { return Inline{Text: text, Role: RoleStrong} }
func Muted(text string) Inline   { return Inline{Text: text, Role: RoleMuted} }
func Code(text string) Inline    { return Inline{Text: text, Role: RoleCode} }
func Path(text string) Inline    { return Inline{Text: text, Role: RolePath} }
func Command(text string) Inline { return Inline{Text: text, Role: RoleCommand} }
func Warning(text string) Inline { return Inline{Text: text, Role: RoleWarning} }
func Status(text string) Inline  { return Inline{Text: text, Role: RoleStatus} }

func Paragraph(inlines ...Inline) Block {
	return Block{Kind: BlockParagraph, Inlines: inlines}
}

func Heading(level int, inlines ...Inline) Block {
	if level <= 0 {
		level = 1
	}
	return Block{Kind: BlockHeading, Level: level, Inlines: inlines}
}

func Blank() Block { return Block{Kind: BlockBlank} }

func CodeBlock(language, code string) Block {
	return Block{Kind: BlockCode, Language: language, Code: code, FinalNewline: strings.HasSuffix(code, "\n")}
}

func DiffBlock(diff string) Block {
	block := Block{Kind: BlockDiff, FinalNewline: strings.HasSuffix(diff, "\n")}
	if block.FinalNewline {
		diff = strings.TrimSuffix(diff, "\n")
	}
	if diff == "" {
		return block
	}
	lines := strings.Split(diff, "\n")
	block.DiffLines = make([]DiffLine, 0, len(lines))
	for _, line := range lines {
		block.DiffLines = append(block.DiffLines, ClassifyDiffLine(line))
	}
	return block
}
