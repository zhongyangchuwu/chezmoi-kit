package testutil

import (
	"io"
	"strings"
	"testing"
)

// TerminalCommand is a test double for terminal-owning commands.
type TerminalCommand struct {
	RunFunc func() error
}

func (c TerminalCommand) Run() error {
	if c.RunFunc == nil {
		return nil
	}
	return c.RunFunc()
}

func (TerminalCommand) SetStdin(io.Reader)  {}
func (TerminalCommand) SetStdout(io.Writer) {}
func (TerminalCommand) SetStderr(io.Writer) {}

func DebugLogPathFromLine(t *testing.T, output, prefix string) string {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	t.Fatalf("output %q missing prefix %q", output, prefix)
	return ""
}
