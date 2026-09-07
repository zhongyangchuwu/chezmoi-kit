package tui

import (
	"os"

	"charm.land/lipgloss/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type tuiStyles struct {
	title          lipgloss.Style
	section        lipgloss.Style
	activePane     lipgloss.Style
	pane           lipgloss.Style
	help           lipgloss.Style
	muted          lipgloss.Style
	clean          lipgloss.Style
	dirty          lipgloss.Style
	unmanaged      lipgloss.Style
	ignored        lipgloss.Style
	script         lipgloss.Style
	uninspected    lipgloss.Style
	template       lipgloss.Style
	encrypted      lipgloss.Style
	directory      lipgloss.Style
	symlink        lipgloss.Style
	error          lipgloss.Style
	loading        lipgloss.Style
	withheld       lipgloss.Style
	selected       lipgloss.Style
	selectedActive lipgloss.Style
	diffHeader     lipgloss.Style
	diffHunk       lipgloss.Style
	diffAdd        lipgloss.Style
	diffRemove     lipgloss.Style
	diffMeta       lipgloss.Style
}

func newTUIStyles() tuiStyles {
	styles := tuiStyles{
		title:          lipgloss.NewStyle().Bold(true),
		section:        lipgloss.NewStyle().Bold(true),
		activePane:     lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		pane:           lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		help:           lipgloss.NewStyle(),
		muted:          lipgloss.NewStyle().Faint(true),
		clean:          lipgloss.NewStyle().Bold(true),
		dirty:          lipgloss.NewStyle().Bold(true),
		unmanaged:      lipgloss.NewStyle().Bold(true),
		ignored:        lipgloss.NewStyle().Faint(true),
		script:         lipgloss.NewStyle().Bold(true),
		uninspected:    lipgloss.NewStyle().Bold(true),
		template:       lipgloss.NewStyle().Bold(true),
		encrypted:      lipgloss.NewStyle().Bold(true),
		directory:      lipgloss.NewStyle().Bold(true),
		symlink:        lipgloss.NewStyle(),
		error:          lipgloss.NewStyle().Bold(true),
		loading:        lipgloss.NewStyle().Faint(true),
		withheld:       lipgloss.NewStyle().Bold(true),
		selected:       lipgloss.NewStyle().Bold(true),
		selectedActive: lipgloss.NewStyle().Bold(true).Reverse(true),
		diffHeader:     lipgloss.NewStyle(),
		diffHunk:       lipgloss.NewStyle().Bold(true),
		diffAdd:        lipgloss.NewStyle(),
		diffRemove:     lipgloss.NewStyle(),
		diffMeta:       lipgloss.NewStyle().Faint(true),
	}
	if os.Getenv("NO_COLOR") != "" {
		return styles
	}
	return tuiStyles{
		title:          styles.title.Foreground(lipgloss.Color("12")),
		section:        styles.section.Foreground(lipgloss.Color("14")),
		activePane:     styles.activePane.BorderForeground(lipgloss.Color("12")),
		pane:           styles.pane.BorderForeground(lipgloss.Color("8")),
		help:           styles.help.Foreground(lipgloss.Color("8")),
		muted:          styles.muted.Foreground(lipgloss.Color("8")),
		clean:          styles.clean.Foreground(lipgloss.Color("10")),
		dirty:          styles.dirty.Foreground(lipgloss.Color("11")),
		unmanaged:      styles.unmanaged.Foreground(lipgloss.Color("14")),
		ignored:        styles.ignored.Foreground(lipgloss.Color("8")),
		script:         styles.script.Foreground(lipgloss.Color("13")),
		uninspected:    styles.uninspected.Foreground(lipgloss.Color("12")),
		template:       styles.template.Foreground(lipgloss.Color("13")),
		encrypted:      styles.encrypted.Foreground(lipgloss.Color("9")),
		directory:      styles.directory.Foreground(lipgloss.Color("12")),
		symlink:        styles.symlink.Foreground(lipgloss.Color("14")),
		error:          styles.error.Foreground(lipgloss.Color("9")),
		loading:        styles.loading.Foreground(lipgloss.Color("12")),
		withheld:       styles.withheld.Foreground(lipgloss.Color("11")),
		selected:       styles.selected.Background(lipgloss.Color("236")),
		selectedActive: styles.selectedActive.Background(lipgloss.Color("12")),
		diffHeader:     styles.diffHeader.Foreground(lipgloss.Color("12")),
		diffHunk:       styles.diffHunk.Foreground(lipgloss.Color("14")),
		diffAdd:        styles.diffAdd.Foreground(lipgloss.Color("10")),
		diffRemove:     styles.diffRemove.Foreground(lipgloss.Color("9")),
		diffMeta:       styles.diffMeta.Foreground(lipgloss.Color("8")),
	}
}

func (s tuiStyles) stateStyle(state app.FileState) lipgloss.Style {
	switch state {
	case app.FileClean:
		return s.clean
	case app.FileDirty:
		return s.dirty
	case app.FileUnmanaged:
		return s.unmanaged
	case app.FileIgnored:
		return s.ignored
	case app.FileScript:
		return s.script
	case app.FileUninspected:
		return s.uninspected
	default:
		return s.muted
	}
}

func (s tuiStyles) typeStyle(targetType app.TargetType) lipgloss.Style {
	switch targetType {
	case app.TargetDirectory:
		return s.directory
	case app.TargetSymlink:
		return s.symlink
	case app.TargetRemove:
		return s.error
	default:
		return lipgloss.NewStyle()
	}
}
