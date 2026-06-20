package tui

import (
	"fmt"
	"io"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
	"golang.org/x/term"
)

type syncDiffMsg struct {
	target string
	diff   string
	err    error
}

type executeMsg struct {
	target   string
	executed []app.Action
	skipped  int
	err      error
}

type terminalRequestMsg struct {
	action app.Action
	cmd    app.TerminalCommand
	err    error
}

type terminalExecuteMsg struct {
	action app.Action
	err    error
}

func RunSyncTUI(service app.SyncService, targets []string, input io.Reader, output io.Writer, options *app.Options) error {
	options = app.NormalizeOptions(options)
	timing, err := newSyncTimingLogger(options.Debug)
	if err != nil {
		return fmt.Errorf("create debug log: %w", err)
	}
	if timing.Enabled() {
		_, _ = fmt.Fprintf(options.Stderr, "debug log: %s\n", timing.Path())
		defer fmt.Fprintf(options.Stderr, "debug log kept at: %s\n", timing.Path())
	}
	defer timing.Close()

	start := time.Now()
	entries, err := service.Status(targets)
	timing.Info("sync initial status", "targets", len(targets), "entries", len(entries), "duration", elapsed(start), "err", err)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		_, err = fmt.Fprintln(output, "clean")
		return err
	}

	programOptions := []tea.ProgramOption{
		tea.WithInput(input),
		tea.WithOutput(output),
	}
	renderFinal := !isTerminalWriter(output)
	if renderFinal {
		programOptions = append(programOptions, tea.WithoutRenderer())
	}

	program := tea.NewProgram(newSyncTUIModel(service, entries, timing), programOptions...)
	model, err := program.Run()
	if err != nil {
		return err
	}
	final, ok := model.(syncTUIModel)
	if !ok {
		return nil
	}
	if final.err != nil {
		return final.err
	}
	if renderFinal {
		_, err = fmt.Fprint(output, final.viewString())
		return err
	}
	return nil
}

func isTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}

func (m syncTUIModel) Init() tea.Cmd {
	if m.service == nil || m.currentTarget() == "" {
		return nil
	}
	return loadDiffCmd(m.service, m.currentTarget(), m.timing)
}

func (m syncTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		if m.mode == modeExecuting {
			return m.updateExecuting(msg)
		}
		if m.mode == modeConfirm {
			return m.updateConfirm(msg)
		}
		return m.updateReview(msg)
	case syncDiffMsg:
		m.applyDiff(msg)
		if msg.err != nil {
			return m, nil
		}
		return m, nil
	case executeMsg:
		return m.applyExecuteMsg(msg)
	case terminalRequestMsg:
		return m.applyTerminalRequestMsg(msg)
	case terminalExecuteMsg:
		return m.applyTerminalExecuteMsg(msg)
	}
	return m, nil
}
