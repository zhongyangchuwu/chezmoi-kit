package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
	"golang.org/x/term"
)

type syncMode int

const (
	modeReview syncMode = iota
	modeConfirm
)

type syncActionMsg struct {
	target string
	diff   string
	err    error
}

type executeMsg struct {
	executed []reconcile.Action
	skipped  int
	err      error
}

type syncTUIModel struct {
	service    reconcile.ReviewService
	entries    []chezmoi.StatusEntry
	cursor     int
	pending    map[string]reconcile.ActionKind
	diffTarget string
	diff       string
	mode       syncMode
	help       help.Model
	message    string
	err        error
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
		tea.WithoutSignals(),
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

func newSyncTUIModel(service reconcile.ReviewService, entries []chezmoi.StatusEntry) syncTUIModel {
	return syncTUIModel{
		service: service,
		entries: append([]chezmoi.StatusEntry(nil), entries...),
		pending: make(map[string]reconcile.ActionKind),
		help:    help.New(),
	}
}

func (m syncTUIModel) Init() tea.Cmd {
	return nil
}

func (m syncTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.mode == modeConfirm {
			return m.updateConfirm(msg)
		}
		return m.updateReview(msg)
	case syncActionMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.message = ""
		m.diffTarget = msg.target
		m.diff = strings.TrimRight(msg.diff, "\n")
		return m, nil
	case executeMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		if len(msg.executed) == 0 {
			m.message = "no pending dirty targets"
		} else {
			m.message = fmt.Sprintf("executed %d action(s)", len(msg.executed))
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m syncTUIModel) updateReview(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Key().Code == tea.KeyEnter || msg.Key().Code == tea.KeyReturn:
		if len(m.pending) == 0 {
			m.message = "no pending actions"
			return m, nil
		}
		m.mode = modeConfirm
		m.message = "confirm pending actions"
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, defaultSyncKeys.Up):
		return m.moveUp(), nil
	case key.Matches(msg, defaultSyncKeys.Down):
		return m.moveDown(), nil
	case key.Matches(msg, defaultSyncKeys.Diff):
		return m, m.loadDiff()
	case key.Matches(msg, defaultSyncKeys.Add):
		return m.togglePending(reconcile.ActionAdd), nil
	case key.Matches(msg, defaultSyncKeys.Apply):
		return m.togglePending(reconcile.ActionApply), nil
	case key.Matches(msg, defaultSyncKeys.Merge):
		return m.togglePending(reconcile.ActionMerge), nil
	case key.Matches(msg, defaultSyncKeys.Skip):
		return m.clearPending(), nil
	default:
		m.message = "unknown choice"
		return m, nil
	}
}

func (m syncTUIModel) updateConfirm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Key().Code == tea.KeyEscape || msg.Key().Code == tea.KeyEsc:
		m.mode = modeReview
		m.message = ""
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Quit):
		return m, tea.Quit
	case msg.String() == "y" || msg.String() == "Y":
		return m, m.executePending()
	default:
		m.message = "confirm with y, esc to review, q to quit"
		return m, nil
	}
}

func (m syncTUIModel) View() tea.View {
	return tea.NewView(m.viewString())
}

func (m syncTUIModel) viewString() string {
	if len(m.entries) == 0 {
		return "clean\n"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("cm sync"))
	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(sectionStyle.Render("Changed files"))
	b.WriteByte('\n')
	for i, entry := range m.entries {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		b.WriteString(cursor)
		b.WriteString(m.pendingLabel(entry.Path))
		b.WriteByte(' ')
		b.WriteString(entry.Path)
		b.WriteByte('\n')
	}

	b.WriteByte('\n')
	b.WriteString(sectionStyle.Render("Diff"))
	b.WriteByte('\n')
	if m.diff == "" || m.diffTarget != m.current().Path {
		b.WriteString("press d to show diff\n")
	} else {
		b.WriteString(renderDiff(m.diff))
		b.WriteByte('\n')
	}

	if m.mode == modeConfirm {
		b.WriteByte('\n')
		b.WriteString(sectionStyle.Render("Confirm"))
		b.WriteByte('\n')
		for _, action := range m.pendingActions() {
			b.WriteString(actionLabel(action.Kind))
			b.WriteByte(' ')
			b.WriteString(action.Target)
			b.WriteByte('\n')
		}
	}

	if m.message != "" {
		b.WriteByte('\n')
		b.WriteString(m.message)
		b.WriteByte('\n')
	}

	b.WriteByte('\n')
	b.WriteString(helpStyle.Render(m.help.ShortHelpView(defaultSyncKeys.ShortHelp())))
	b.WriteByte('\n')
	return b.String()
}

func renderDiff(diff string) string {
	if diff == "" {
		return ""
	}

	lines := strings.Split(diff, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "@@"):
			lines[i] = diffHunkStyle.Render(line)
		case strings.HasPrefix(line, "diff "), strings.HasPrefix(line, "---"), strings.HasPrefix(line, "+++"):
			lines[i] = diffHeaderStyle.Render(line)
		case strings.HasPrefix(line, "+"):
			lines[i] = diffAddStyle.Render(line)
		case strings.HasPrefix(line, "-"):
			lines[i] = diffRemoveStyle.Render(line)
		case strings.HasPrefix(line, `\ No newline`):
			lines[i] = diffMetaStyle.Render(line)
		}
	}
	return strings.Join(lines, "\n")
}

func (m syncTUIModel) current() chezmoi.StatusEntry {
	if len(m.entries) == 0 {
		return chezmoi.StatusEntry{}
	}
	return m.entries[m.cursor]
}

func (m syncTUIModel) moveUp() syncTUIModel {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m syncTUIModel) moveDown() syncTUIModel {
	if m.cursor < len(m.entries)-1 {
		m.cursor++
	}
	return m
}

func (m syncTUIModel) togglePending(kind reconcile.ActionKind) syncTUIModel {
	target := m.current().Path
	if current, ok := m.pending[target]; ok && current == kind {
		delete(m.pending, target)
		m.message = "cleared " + target
		return m
	}
	m.pending[target] = kind
	m.message = actionLabel(kind) + " " + target
	return m
}

func (m syncTUIModel) clearPending() syncTUIModel {
	target := m.current().Path
	delete(m.pending, target)
	m.message = "skipped " + target
	return m
}

func (m syncTUIModel) pendingActions() []reconcile.Action {
	actions := make([]reconcile.Action, 0, len(m.pending))
	for _, entry := range m.entries {
		kind, ok := m.pending[entry.Path]
		if !ok {
			continue
		}
		actions = append(actions, reconcile.Action{Target: entry.Path, Kind: kind})
	}
	return actions
}

func (m syncTUIModel) loadDiff() tea.Cmd {
	target := m.current().Path
	if target == "" {
		return nil
	}
	if m.diffTarget == target && m.diff != "" {
		return nil
	}
	return func() tea.Msg {
		out, err := m.service.DiffOutput(target)
		return syncActionMsg{target: target, diff: string(out), err: err}
	}
}

func (m syncTUIModel) executePending() tea.Cmd {
	actions := m.pendingActions()
	return func() tea.Msg {
		if len(actions) == 0 {
			return executeMsg{}
		}
		targets := make([]string, 0, len(actions))
		for _, action := range actions {
			targets = append(targets, action.Target)
		}
		fresh, err := m.service.Status(targets)
		if err != nil {
			return executeMsg{err: err}
		}
		dirty := make(map[string]struct{}, len(fresh))
		for _, entry := range fresh {
			dirty[entry.Path] = struct{}{}
		}
		kept := actions[:0]
		for _, action := range actions {
			if _, ok := dirty[action.Target]; ok {
				kept = append(kept, action)
			}
		}
		if len(kept) == 0 {
			return executeMsg{skipped: len(actions)}
		}
		if err := m.service.Execute(kept); err != nil {
			return executeMsg{err: err}
		}
		return executeMsg{executed: kept, skipped: len(actions) - len(kept)}
	}
}

func (m syncTUIModel) pendingLabel(target string) string {
	kind, ok := m.pending[target]
	if !ok {
		return "[ ]"
	}
	return "[" + actionMarker(kind) + "]"
}

func actionLabel(kind reconcile.ActionKind) string {
	switch kind {
	case reconcile.ActionAdd:
		return "add"
	case reconcile.ActionApply:
		return "apply"
	case reconcile.ActionMerge:
		return "merge"
	default:
		return "?"
	}
}

func actionMarker(kind reconcile.ActionKind) string {
	switch kind {
	case reconcile.ActionAdd:
		return "A"
	case reconcile.ActionApply:
		return "P"
	case reconcile.ActionMerge:
		return "M"
	default:
		return "?"
	}
}

type syncKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Diff    key.Binding
	Add     key.Binding
	Apply   key.Binding
	Merge   key.Binding
	Skip    key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

func (k syncKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Diff, k.Add, k.Apply, k.Merge, k.Skip, k.Confirm, k.Quit}
}

func (k syncKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var (
	defaultSyncKeys = syncKeyMap{
		Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move")),
		Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move")),
		Diff:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Apply:   key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "apply")),
		Merge:   key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "merge")),
		Skip:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip")),
		Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sectionStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	diffHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	diffHunkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	diffMetaStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
