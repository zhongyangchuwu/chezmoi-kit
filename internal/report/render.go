package report

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ColorMode controls ANSI styling.
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

// Options configures report rendering.
type Options struct {
	Color ColorMode
	IsTTY bool
	Env   map[string]string
}

func Plain(doc Document) []byte {
	return render(doc, Options{Color: ColorNever}, false)
}

func ANSI(doc Document, opts Options) []byte {
	return render(doc, opts, useColor(opts))
}

func Markdown(doc Document) []byte {
	var out bytes.Buffer
	for _, block := range doc.Blocks {
		switch block.Kind {
		case BlockHeading:
			level := block.Level
			if level <= 0 {
				level = 1
			}
			out.WriteString(strings.Repeat("#", level))
			out.WriteByte(' ')
			writeMarkdownInlines(&out, block.Inlines)
			out.WriteByte('\n')
		case BlockParagraph:
			writeMarkdownInlines(&out, block.Inlines)
			out.WriteByte('\n')
		case BlockBlank:
			out.WriteByte('\n')
		case BlockCode:
			out.WriteString("```")
			out.WriteString(block.Language)
			out.WriteByte('\n')
			out.WriteString(block.Code)
			if !strings.HasSuffix(block.Code, "\n") {
				out.WriteByte('\n')
			}
			out.WriteString("```\n")
		case BlockDiff:
			out.WriteString("```diff\n")
			writeDiffLines(&out, block, false)
			if len(block.DiffLines) > 0 && !block.FinalNewline {
				out.WriteByte('\n')
			}
			out.WriteString("```\n")
		}
	}
	return out.Bytes()
}

func useColor(opts Options) bool {
	if noColor(opts) {
		return false
	}
	switch opts.Color {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	default:
		return opts.IsTTY
	}
}

func noColor(opts Options) bool {
	if opts.Env != nil {
		return opts.Env["NO_COLOR"] != ""
	}
	return os.Getenv("NO_COLOR") != ""
}

func render(doc Document, opts Options, color bool) []byte {
	var out bytes.Buffer
	for _, block := range doc.Blocks {
		switch block.Kind {
		case BlockHeading, BlockParagraph:
			writeStyledInlines(&out, block.Inlines, color)
			out.WriteByte('\n')
		case BlockBlank:
			out.WriteByte('\n')
		case BlockCode:
			out.WriteString(block.Code)
			if !strings.HasSuffix(block.Code, "\n") {
				out.WriteByte('\n')
			}
		case BlockDiff:
			writeDiffLines(&out, block, color)
		}
	}
	_ = opts
	return out.Bytes()
}

func writeStyledInlines(w io.Writer, inlines []Inline, color bool) {
	for _, inline := range inlines {
		text := inline.Text
		if color {
			text = style(inline.Role, text)
		}
		_, _ = io.WriteString(w, text)
	}
}

func writeMarkdownInlines(w io.Writer, inlines []Inline) {
	for _, inline := range inlines {
		text := inline.Text
		switch inline.Role {
		case RoleStrong, RoleWarning, RoleStatus:
			_, _ = fmt.Fprintf(w, "**%s**", text)
		case RoleCode, RolePath, RoleCommand:
			_, _ = fmt.Fprintf(w, "`%s`", text)
		default:
			_, _ = io.WriteString(w, text)
		}
	}
}

func writeDiffLines(w io.Writer, block Block, color bool) {
	for i, line := range block.DiffLines {
		text := line.Text
		if color {
			text = style(DiffLineRole(line.Kind), text)
		}
		_, _ = io.WriteString(w, text)
		if i < len(block.DiffLines)-1 || block.FinalNewline {
			_, _ = io.WriteString(w, "\n")
		}
	}
}

func style(role Role, text string) string {
	code := ansiCode(role)
	if code == "" || text == "" {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}

func ansiCode(role Role) string {
	switch role {
	case RoleStrong:
		return "1"
	case RoleMuted:
		return "90"
	case RoleCode, RolePath, RoleCommand:
		return "36;1"
	case RoleWarning, RoleStatus:
		return "33;1"
	case RoleDiffHeader:
		return "36"
	case RoleDiffHunk:
		return "36;1"
	case RoleDiffAdd:
		return "32"
	case RoleDiffRemove:
		return "31"
	case RoleDiffMeta:
		return "90"
	default:
		return ""
	}
}
