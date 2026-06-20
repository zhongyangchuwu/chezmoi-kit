package ui

import (
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
	"golang.org/x/term"
)

type syncDiffMsg struct {
	target string
	diff   string
	err    error
}

type executeMsg struct {
	target   string
	executed []reconcile.Action
	skipped  int
	err      error
}

func RunSyncTUI(service reconcile.ReviewService, targets []string, input io.Reader, output io.Writer) error {
	entries, err := service.Status(targets)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		_, err = fmt.Fprintln(output, "clean")
		return err
	}

	options := []tea.ProgramOption{
		tea.WithInput(input),
		tea.WithOutput(output),
	}
	renderFinal := !isTerminalWriter(output)
	if renderFinal {
		options = append(options, tea.WithoutRenderer())
	}

	program := tea.NewProgram(newSyncTUIModel(service, entries), options...)
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
	return loadDiffCmd(m.service, m.currentTarget())
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
	}
	return m, nil
}
