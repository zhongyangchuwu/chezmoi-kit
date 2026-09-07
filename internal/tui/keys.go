package tui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type syncKeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Left            key.Binding
	Right           key.Binding
	PageUp          key.Binding
	PageDown        key.Binding
	Tab             key.Binding
	Diff            key.Binding
	Add             key.Binding
	Apply           key.Binding
	Merge           key.Binding
	Skip            key.Binding
	Enter           key.Binding
	Execute         key.Binding
	Back            key.Binding
	Quit            key.Binding
	Tree            key.Binding
	Filter          key.Binding
	Search          key.Binding
	Full            key.Binding
	ViewDiff        key.Binding
	ViewDestination key.Binding
	ViewTarget      key.Binding
	ViewSource      key.Binding
	Reveal          key.Binding
	HunkPrevious    key.Binding
	HunkNext        key.Binding
	MatchPrevious   key.Binding
	MatchNext       key.Binding
	ToggleDirectory key.Binding
}

func (k syncKeyMap) reviewHelp(focus syncFocus, review app.Review) []key.Binding {
	bindings := []key.Binding{k.Tab, k.Up, k.Down, k.Diff}
	if review.Allows(app.ActionAdd) {
		bindings = append(bindings, k.Add)
	}
	if review.Allows(app.ActionApply) {
		bindings = append(bindings, k.Apply)
	}
	if review.Allows(app.ActionMerge) {
		bindings = append(bindings, k.Merge)
	}
	if focus == focusFiles {
		bindings = append(bindings, k.Skip)
	}
	return append(bindings, k.Enter, k.Quit)
}

func (k syncKeyMap) confirmHelp() []key.Binding {
	return []key.Binding{k.Execute, k.Back, k.Quit}
}

func (k syncKeyMap) executingHelp() []key.Binding {
	return nil
}

func (k syncKeyMap) ShortHelp() []key.Binding {
	return k.reviewHelp(focusFiles, app.Review{})
}

func (k syncKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var (
	defaultSyncKeys = syncKeyMap{
		Up:              key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:            key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:            key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right:           key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		PageUp:          key.NewBinding(key.WithKeys("ctrl+u", "pgup"), key.WithHelp("ctrl+u", "page up")),
		PageDown:        key.NewBinding(key.WithKeys("ctrl+d", "pgdown"), key.WithHelp("ctrl+d", "page down")),
		Tab:             key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focus")),
		Diff:            key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "reload diff")),
		Add:             key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Apply:           key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "apply")),
		Merge:           key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "merge")),
		Skip:            key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip")),
		Enter:           key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		Execute:         key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "execute")),
		Back:            key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Quit:            key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Tree:            key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tree/flat")),
		Filter:          key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "filter")),
		Search:          key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Full:            key.NewBinding(key.WithKeys("z"), key.WithHelp("z", "full preview")),
		ViewDiff:        key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "diff")),
		ViewDestination: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "destination")),
		ViewTarget:      key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "target")),
		ViewSource:      key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "source")),
		Reveal:          key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "reveal")),
		HunkPrevious:    key.NewBinding(key.WithKeys("["), key.WithHelp("[", "previous hunk")),
		HunkNext:        key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "next hunk")),
		MatchPrevious:   key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "previous match")),
		MatchNext:       key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next match")),
		ToggleDirectory: key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "collapse")),
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
