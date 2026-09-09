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

type syncReviewMsg struct {
	target string
	review app.Review
	err    error
}

type workspacePreviewMsg struct {
	key     string
	epoch   uint64
	preview app.WorkspacePreview
	err     error
}

type workspaceHandoffRequestMsg struct {
	entry   app.WorkspaceEntry
	command app.TerminalCommand
	err     error
}

type workspaceHandoffDoneMsg struct {
	target string
	err    error
}

type workspaceRefreshMsg struct {
	snapshot  app.WorkspaceSnapshot
	target    string
	editorErr error
	err       error
}

type executeMsg struct {
	target   string
	action   app.Action
	result   app.ActionResult
	executed bool
	skipped  bool
	resolved bool
	deferred bool
	review   app.Review
	reason   string
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
	status, err := service.Status(targets)
	timing.Info("sync initial status", "targets", len(targets), "entries", len(status.Entries), "scripts", len(status.Scripts), "duration", elapsed(start), "err", err)
	if err != nil {
		return err
	}
	if len(status.Entries) == 0 {
		if len(status.Scripts) > 0 {
			_, err = fmt.Fprintf(output, "no file changes to reconcile; %s pending\nuse chezmoi diff and chezmoi apply to review and run scripts\n", scriptCount(len(status.Scripts)))
			return err
		}
		_, err = fmt.Fprintln(output, "clean")
		return err
	}
	return runTUIProgram(newSyncTUIModel(service, status, timing), input, output)
}

func RunWorkspaceTUI(service app.WorkspaceService, scopes []string, input io.Reader, output io.Writer, options *app.Options) error {
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
	snapshot, err := service.Inventory(scopes)
	timing.Info("workspace inventory", "scopes", len(scopes), "entries", len(snapshot.Entries), "duration", elapsed(start), "err", err)
	if err != nil {
		return err
	}
	return runTUIProgram(newWorkspaceModel(service, snapshot, timing), input, output)
}

func runTUIProgram(initial workspaceModel, input io.Reader, output io.Writer) error {
	programOptions := []tea.ProgramOption{
		tea.WithInput(input),
		tea.WithOutput(output),
	}
	renderFinal := !isTerminalWriter(output)
	if renderFinal {
		programOptions = append(programOptions, tea.WithoutRenderer())
	}

	program := tea.NewProgram(initial, programOptions...)
	model, err := program.Run()
	if err != nil {
		return err
	}
	final, ok := model.(workspaceModel)
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

func (m workspaceModel) Init() tea.Cmd {
	if m.service == nil || m.currentTarget() == "" {
		return nil
	}
	_, cmd := m.startDiffLoad(false)
	return cmd
}

func (m workspaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		if m.isWorkspace() && m.workspaceBusy {
			return m, nil
		}
		if m.isWorkspace() && m.helpVisible {
			return m.updateWorkspaceHelp(msg)
		}
		if m.search != searchNone {
			return m.updateSearch(msg)
		}
		if m.mode == modeExecuting {
			return m.updateExecuting(msg)
		}
		if m.mode == modeConfirm {
			return m.updateConfirm(msg)
		}
		return m.updateReview(msg)
	case syncReviewMsg:
		return m.applyReview(msg)
	case workspacePreviewMsg:
		return m.applyWorkspacePreview(msg)
	case workspaceHandoffRequestMsg:
		return m.applyWorkspaceHandoffRequest(msg)
	case workspaceHandoffDoneMsg:
		return m.applyWorkspaceHandoffDone(msg)
	case workspaceRefreshMsg:
		return m.applyWorkspaceRefresh(msg)
	case executeMsg:
		return m.applyExecuteMsg(msg)
	case terminalRequestMsg:
		return m.applyTerminalRequestMsg(msg)
	case terminalExecuteMsg:
		return m.applyTerminalExecuteMsg(msg)
	}
	return m, nil
}
