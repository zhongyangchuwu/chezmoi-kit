package report

import (
	"strings"
	"testing"
)

func TestRenderPlainPreservesSemanticText(t *testing.T) {
	doc := Document{Blocks: []Block{
		Heading(1, Strong("local:")),
		Paragraph(Warning("!"), Text(" "), Path("/home/me/.zshrc"), Text("  differs from chezmoi")),
		Paragraph(Muted("  run cm sync")),
	}}

	got := string(Plain(doc))
	want := "local:\n! /home/me/.zshrc  differs from chezmoi\n  run cm sync\n"
	if got != want {
		t.Fatalf("Plain() = %q, want %q", got, want)
	}
}

func TestANSIUsesPaletteWhenTTY(t *testing.T) {
	doc := Document{Blocks: []Block{Paragraph(Warning("!"), Text(" "), Path("/tmp/a"))}}

	got := string(ANSI(doc, Options{Color: ColorAuto, IsTTY: true, Env: map[string]string{}}))
	for _, want := range []string{"\x1b[33;1m!\x1b[0m", "\x1b[36;1m/tmp/a\x1b[0m"} {
		if !strings.Contains(got, want) {
			t.Fatalf("ANSI() = %q, missing %q", got, want)
		}
	}
}

func TestNoColorDisablesAutoANSI(t *testing.T) {
	doc := Document{Blocks: []Block{Paragraph(Warning("!"), Text(" "), Path("/tmp/a"))}}

	got := string(ANSI(doc, Options{Color: ColorAuto, IsTTY: true, Env: map[string]string{"NO_COLOR": "1"}}))
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("ANSI() = %q, want no escape sequences", got)
	}
	if got != "! /tmp/a\n" {
		t.Fatalf("ANSI() = %q, want plain text", got)
	}
}

func TestColorAlwaysOverridesNoColor(t *testing.T) {
	doc := Document{Blocks: []Block{Paragraph(Warning("!"), Text(" "), Path("/tmp/a"))}}

	got := string(ANSI(doc, Options{Color: ColorAlways, IsTTY: false, Env: map[string]string{"NO_COLOR": "1"}}))
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("ANSI() = %q, want ANSI escapes", got)
	}
}

func TestMarkdownRendererIsDeterministic(t *testing.T) {
	doc := Document{Blocks: []Block{
		Heading(2, Text("Status")),
		Paragraph(Strong("cm:"), Text(" "), Code("v0.1.0")),
	}}

	got := string(Markdown(doc))
	want := "## Status\n**cm:** `v0.1.0`\n"
	if got != want {
		t.Fatalf("Markdown() = %q, want %q", got, want)
	}
}

func TestDiffLineClassification(t *testing.T) {
	tests := []struct {
		line string
		kind DiffLineKind
	}{
		{line: "diff --git a b", kind: DiffLineHeader},
		{line: "--- a", kind: DiffLineHeader},
		{line: "+++ b", kind: DiffLineHeader},
		{line: "@@ -1 +1 @@", kind: DiffLineHunk},
		{line: "+new", kind: DiffLineAdd},
		{line: "-old", kind: DiffLineRemove},
		{line: `\ No newline at end of file`, kind: DiffLineMeta},
		{line: " context", kind: DiffLineContext},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := ClassifyDiffLine(tt.line); got.Kind != tt.kind || got.Text != tt.line {
				t.Fatalf("ClassifyDiffLine(%q) = %#v, want kind %v", tt.line, got, tt.kind)
			}
		})
	}
}

func TestDiffBlockPreservesFinalNewline(t *testing.T) {
	doc := Document{Blocks: []Block{DiffBlock("--- a\n+++ b\n")}}

	got := string(Plain(doc))
	want := "--- a\n+++ b\n"
	if got != want {
		t.Fatalf("Plain(diff) = %q, want %q", got, want)
	}
}
