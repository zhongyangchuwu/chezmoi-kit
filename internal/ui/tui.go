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
	"golang.org/x/term"
)

type syncActionKind int

const (
	syncActionDiff syncActionKind = iota
	syncActionAdd
	syncActionApply
	syncActionMerge
	syncActionSkip
	syncActionQuit
)

type syncAction struct {
	kind   syncActionKind
	target string
}

type syncActionMsg struct {
	action syncAction
	err    error
	diff   string
}

type syncTUIModel struct {
	service SyncTUIService
	entries []chezmoi.StatusEntry
	diffs   []string
	cursor  int
	help    help.Model
	message string
	err     error
}

func RunSyncTUI(service SyncTUIService, targets []string, input io.Reader, output io.Writer) error {
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

func newSyncTUIModel(service SyncTUIService, entries []chezmoi.StatusEntry) syncTUIModel {
	model := newSyncModel(entries, nil)
	model.service = service
	return model
}

func newSyncModel(entries []chezmoi.StatusEntry, diffs []string) syncTUIModel {
	return syncTUIModel{
		entries: append([]chezmoi.StatusEntry(nil), entries...),
		diffs:   append([]string(nil), diffs...),
		help:    help.New(),
	}
}

func (m syncTUIModel) Init() tea.Cmd {
	return nil
}

func (m syncTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Key().Code == tea.KeyEnter || msg.Key().Code == tea.KeyReturn:
			return m, nil
		case key.Matches(msg, defaultSyncKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, defaultSyncKeys.Up):
			return m.moveUp(), nil
		case key.Matches(msg, defaultSyncKeys.Down):
			return m.moveDown(), nil
		default:
			action, ok := m.actionForKey(msg.String())
			if !ok {
				m.message = "unknown choice"
				return m, nil
			}
			if action.kind == syncActionQuit {
				return m, tea.Quit
			}
			return m, m.runAction(action)
		}
	case syncActionMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m = m.applyActionResult(msg.action, msg.diff)
		if len(m.entries) == 0 {
			m.message = "clean"
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
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
		b.WriteString(entry.Path)
		b.WriteByte('\n')
	}

	b.WriteByte('\n')
	b.WriteString(sectionStyle.Render("Diff"))
	b.WriteByte('\n')
	if len(m.diffs) == 0 {
		b.WriteString("press d to show diff\n")
	} else {
		b.WriteString(renderDiff(strings.Join(m.diffs, "\n")))
		b.WriteByte('\n')
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

func (m syncTUIModel) actionForKey(name string) (syncAction, bool) {
	if len(m.entries) == 0 {
		return syncAction{}, false
	}
	target := m.current().Path
	switch name {
	case "d":
		return syncAction{kind: syncActionDiff, target: target}, true
	case "a":
		return syncAction{kind: syncActionAdd, target: target}, true
	case "p":
		return syncAction{kind: syncActionApply, target: target}, true
	case "m":
		return syncAction{kind: syncActionMerge, target: target}, true
	case "s":
		return syncAction{kind: syncActionSkip, target: target}, true
	case "q":
		return syncAction{kind: syncActionQuit, target: target}, true
	default:
		return syncAction{}, false
	}
}

func (m syncTUIModel) runAction(action syncAction) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch action.kind {
		case syncActionDiff:
			var out []byte
			out, err = m.service.DiffOutput([]string{action.target})
			return syncActionMsg{action: action, diff: string(out), err: err}
		case syncActionAdd:
			err = m.service.Add(action.target)
		case syncActionApply:
			err = m.service.Apply(action.target)
		case syncActionMerge:
			err = m.service.Merge(action.target)
		}
		return syncActionMsg{action: action, err: err}
	}
}

func (m syncTUIModel) applyActionResult(action syncAction, diff string) syncTUIModel {
	m.message = ""
	switch action.kind {
	case syncActionDiff:
		if diff == "" {
			diff = "diff shown for " + action.target
		}
		m.diffs = []string{strings.TrimRight(diff, "\n")}
		return m
	case syncActionSkip:
		return m.withEntryClean(action.target)
	case syncActionAdd, syncActionApply, syncActionMerge:
		if m.service == nil {
			return m.withEntryClean(action.target)
		}
		fresh, err := m.service.Status([]string{action.target})
		if err != nil {
			m.err = err
			return m
		}
		if len(fresh) == 0 {
			return m.withEntryClean(action.target)
		}
		m.entries[m.cursor] = fresh[0]
	}
	return m
}

func (m syncTUIModel) withEntryClean(target string) syncTUIModel {
	for i, entry := range m.entries {
		if entry.Path != target {
			continue
		}
		m.entries = append(m.entries[:i], m.entries[i+1:]...)
		if m.cursor >= len(m.entries) && m.cursor > 0 {
			m.cursor--
		}
		return m
	}
	return m
}

type syncKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Diff  key.Binding
	Add   key.Binding
	Apply key.Binding
	Merge key.Binding
	Skip  key.Binding
	Quit  key.Binding
}

func (k syncKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Diff, k.Add, k.Apply, k.Merge, k.Skip, k.Quit}
}

func (k syncKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var (
	defaultSyncKeys = syncKeyMap{
		Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move")),
		Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move")),
		Diff:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		Add:   key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Apply: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "apply")),
		Merge: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "merge")),
		Skip:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip")),
		Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
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
