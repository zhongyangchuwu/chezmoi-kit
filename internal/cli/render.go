package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/zhongyangchuwu/cm/internal/report"
	"golang.org/x/term"
)

type outputFormat string

const (
	outputPlain    outputFormat = "plain"
	outputANSI     outputFormat = "ansi"
	outputMarkdown outputFormat = "markdown"
)

type colorPolicy string

const (
	colorAuto   colorPolicy = "auto"
	colorAlways colorPolicy = "always"
	colorNever  colorPolicy = "never"
)

type renderOptions struct {
	Output string
	Color  string
}

func defaultRenderOptions() renderOptions {
	return renderOptions{Output: string(outputANSI), Color: string(colorAuto)}
}

func renderReport(w io.Writer, doc report.Document, opts renderOptions) error {
	out, err := renderBytes(w, doc, opts)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func renderBytes(w io.Writer, doc report.Document, opts renderOptions) ([]byte, error) {
	color, err := parseColorPolicy(opts.Color)
	if err != nil {
		return nil, err
	}
	switch outputFormat(opts.Output) {
	case outputPlain:
		return report.Plain(doc), nil
	case outputANSI:
		return report.ANSI(doc, report.Options{Color: color, IsTTY: isTerminalWriter(w)}), nil
	case outputMarkdown:
		return report.Markdown(doc), nil
	default:
		return nil, fmt.Errorf("unsupported output %q", opts.Output)
	}
}

func parseColorPolicy(value string) (report.ColorMode, error) {
	switch colorPolicy(value) {
	case colorAuto:
		return report.ColorAuto, nil
	case colorAlways:
		return report.ColorAlways, nil
	case colorNever:
		return report.ColorNever, nil
	default:
		return report.ColorAuto, fmt.Errorf("unsupported color %q", value)
	}
}

func isTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}
