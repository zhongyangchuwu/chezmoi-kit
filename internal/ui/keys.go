package ui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

type syncKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Tab     key.Binding
	Diff    key.Binding
	Add     key.Binding
	Apply   key.Binding
	Merge   key.Binding
	Skip    key.Binding
	Enter   key.Binding
	Execute key.Binding
	Back    key.Binding
	Quit    key.Binding
}

func (k syncKeyMap) reviewHelp(focus syncFocus) []key.Binding {
	if focus == focusDiff {
		return []key.Binding{k.Tab, k.Up, k.Down, k.Diff, k.Add, k.Apply, k.Merge, k.Enter, k.Quit}
	}
	return []key.Binding{k.Tab, k.Up, k.Down, k.Diff, k.Add, k.Apply, k.Merge, k.Skip, k.Enter, k.Quit}
}

func (k syncKeyMap) confirmHelp() []key.Binding {
	return []key.Binding{k.Execute, k.Back, k.Quit}
}

func (k syncKeyMap) executingHelp() []key.Binding {
	return nil
}

func (k syncKeyMap) ShortHelp() []key.Binding {
	return k.reviewHelp(focusFiles)
}

func (k syncKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var (
	defaultSyncKeys = syncKeyMap{
		Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Tab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focus")),
		Diff:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "refresh diff")),
		Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Apply:   key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "apply")),
		Merge:   key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "merge")),
		Skip:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip")),
		Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		Execute: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "execute")),
		Back:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sectionStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	activePaneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("12"))
	paneStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8"))
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	diffHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	diffHunkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	diffMetaStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
