package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type diffState struct {
	review  app.Review
	preview app.WorkspacePreview
	lines   []string
	loading bool
	err     error
}

func (m workspaceModel) startDiffLoad(refresh bool) (workspaceModel, tea.Cmd) {
	if m.isWorkspace() {
		return m.startWorkspacePreviewLoad(refresh)
	}
	target := m.currentTarget()
	if target == "" || m.service == nil {
		return m, nil
	}
	key := m.previewStateKey()
	if state, ok := m.diffs[key]; ok && state.review.Fingerprint != "" && !refresh {
		return m, nil
	}
	m.diffs[key] = diffState{loading: true}
	return m, loadReviewCmd(m.service, target, m.timing)
}

func (m workspaceModel) startWorkspacePreviewLoad(refresh bool) (workspaceModel, tea.Cmd) {
	entry := m.current()
	if entry.Path == "" {
		return m, nil
	}
	key := m.previewStateKey()
	if entry.Type == app.TargetDirectory {
		m.diffs[key] = diffState{loading: true}
		direct := len(m.directoryEntries(entry.RelativePath))
		return m, loadDirectoryPreviewCmd(entry, direct, m.allEntries, m.filter, m.previewKind, m.previewEpoch)
	}
	if state, ok := m.diffs[key]; ok && !state.loading && state.err == nil && !refresh {
		return m, nil
	}
	m.diffs[key] = diffState{loading: true}
	if m.workspace == nil {
		return m, nil
	}
	return m, loadWorkspacePreviewCmd(m.workspace, entry, m.previewKind, m.currentPreviewRevealed(), m.previewEpoch, m.timing)
}

func loadDirectoryPreviewCmd(entry app.WorkspaceEntry, direct int, allEntries []app.WorkspaceEntry, filter workspaceFilter, kind app.PreviewKind, epoch uint64) tea.Cmd {
	directory := cleanWorkspaceRelative(entry.RelativePath)
	descendants := 0
	counts := map[app.FileState]int{}
	for _, candidate := range allEntries {
		candidateRelative := cleanWorkspaceRelative(candidate.RelativePath)
		if !strings.HasPrefix(candidateRelative, directory+"/") || !workspaceEntryMatchesFilter(candidate, filter) {
			continue
		}
		descendants++
		counts[candidate.State]++
	}
	lines := []string{
		fmt.Sprintf("direct entries: %d", direct),
		fmt.Sprintf("matching descendants: %d", descendants),
	}
	for _, state := range []app.FileState{app.FileDirty, app.FileUninspected, app.FileUnmanaged, app.FileIgnored, app.FileScript, app.FileClean} {
		if counts[state] > 0 {
			lines = append(lines, fmt.Sprintf("%s: %d", stateMarker(state), counts[state]))
		}
	}
	message := workspacePreviewMsg{
		key:     workspacePreviewKey(entry.Path, kind, false),
		epoch:   epoch,
		preview: app.WorkspacePreview{Entry: entry, Kind: kind, Content: strings.Join(lines, "\n")},
	}
	return func() tea.Msg {
		return message
	}
}

func loadReviewCmd(service interface {
	Review(string) (app.Review, error)
}, target string, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		review, err := service.Review(target)
		timing.Info("sync review", "target", target, "dirty", review.Dirty, "type", review.Type, "template", review.Template, "duration", elapsed(start), "err", err)
		return syncReviewMsg{target: target, review: review, err: err}
	}
}

func loadWorkspacePreviewCmd(service app.WorkspaceService, entry app.WorkspaceEntry, kind app.PreviewKind, reveal bool, epoch uint64, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		preview, err := service.Preview(entry, kind, reveal)
		timing.Info("workspace preview", "target", entry.Path, "kind", kind, "reveal", reveal, "bytes", len(preview.Content), "withheld", preview.Withheld, "duration", elapsed(start), "err", err)
		return workspacePreviewMsg{key: workspacePreviewKey(entry.Path, kind, reveal), epoch: epoch, preview: preview, err: err}
	}
}

func (m workspaceModel) applyReview(msg syncReviewMsg) (workspaceModel, tea.Cmd) {
	if msg.err != nil {
		m.diffs[msg.target] = diffState{err: msg.err}
		m.message = msg.err.Error()
		return m, nil
	}
	if !msg.review.Dirty {
		m = m.removeEntry(msg.target)
		if len(m.entries) == 0 {
			m.completed = true
			m.stopped = true
			m.message = "clean"
			return m, tea.Quit
		}
		m.message = "already clean: " + m.displayPath(msg.target)
		return m.startDiffLoad(false)
	}
	m.storeReview(msg.review)
	m.message = ""
	return m, nil
}

func (m workspaceModel) applyWorkspacePreview(msg workspacePreviewMsg) (workspaceModel, tea.Cmd) {
	if msg.epoch != m.previewEpoch {
		return m, nil
	}
	current := msg.key == m.previewStateKey()
	if msg.err != nil {
		m.diffs[msg.key] = diffState{err: msg.err}
		if current {
			m.message = msg.err.Error()
		}
		return m, nil
	}
	m.diffs[msg.key] = diffState{
		review:  msg.preview.Review,
		preview: msg.preview,
		lines:   workspacePreviewLines(msg.preview),
	}
	m.replaceWorkspaceEntry(msg.preview.Entry)
	if current {
		m.message = ""
		m.updatePreviewMatches()
	}
	return m, nil
}

func (m *workspaceModel) storeReview(review app.Review) {
	if pending, ok := m.pending[review.Entry.Path]; ok && pending.Fingerprint != review.Fingerprint {
		delete(m.pending, review.Entry.Path)
	}
	lines := []string{"target: " + reviewLabel(review), ""}
	if review.Diff == "" {
		lines = append(lines, "no diff")
	} else {
		lines = append(lines, splitLines(strings.TrimRight(review.Diff, "\n"))...)
	}
	m.diffs[review.Entry.Path] = diffState{review: review, lines: lines}
}

func workspacePreviewLines(preview app.WorkspacePreview) []string {
	entry := preview.Entry
	attributes := make([]string, 0, 2)
	if entry.Template {
		attributes = append(attributes, "template")
	}
	if entry.Encrypted {
		attributes = append(attributes, "encrypted")
	}
	label := "view: " + string(preview.Kind) + " • state: " + string(entry.State) + " • type: " + string(entry.Type)
	if entry.Type == app.TargetDirectory {
		label = "view: directory summary • state: " + string(entry.State) + " • type: " + string(entry.Type)
	}
	if len(attributes) > 0 {
		label += " • " + strings.Join(attributes, ",")
	}
	lines := []string{label, "path: " + entry.Path}
	if entry.SourcePath != "" {
		lines = append(lines, "source: "+entry.SourcePath)
	}
	lines = append(lines, "")
	if preview.Notice != "" {
		lines = append(lines, preview.Notice)
	}
	if preview.Content != "" {
		lines = append(lines, splitLines(strings.TrimRight(preview.Content, "\n"))...)
	}
	if preview.Notice == "" && preview.Content == "" {
		lines = append(lines, "empty")
	}
	return lines
}

func reviewLabel(review app.Review) string {
	label := string(review.Type)
	if review.Template {
		label += " template"
	}
	return label
}

func (m workspaceModel) currentDiffState() diffState {
	return m.diffs[m.previewStateKey()]
}

func (m workspaceModel) scrollDiff(delta int, height int) workspaceModel {
	lines := m.currentDiffState().lines
	maxScroll := len(lines) - height
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.diffScroll += delta
	if m.diffScroll < 0 {
		m.diffScroll = 0
	}
	if m.diffScroll > maxScroll {
		m.diffScroll = maxScroll
	}
	return m
}

func (m workspaceModel) scrollPreviewHorizontal(delta int, width int) workspaceModel {
	maxWidth := 0
	for _, line := range m.currentDiffState().lines {
		if lineWidth := lipgloss.Width(line); lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}
	maxScroll := maxWidth - width
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.previewX += delta
	if m.previewX < 0 {
		m.previewX = 0
	}
	if m.previewX > maxScroll {
		m.previewX = maxScroll
	}
	return m
}

func (m workspaceModel) jumpHunk(direction int) workspaceModel {
	lines := m.currentDiffState().lines
	if direction >= 0 {
		for i := m.diffScroll + 1; i < len(lines); i++ {
			if strings.HasPrefix(lines[i], "@@") {
				m.diffScroll = i
				return m
			}
		}
		return m
	}
	for i := m.diffScroll - 1; i >= 0; i-- {
		if strings.HasPrefix(lines[i], "@@") {
			m.diffScroll = i
			return m
		}
	}
	return m
}

func (m *workspaceModel) updatePreviewMatches() {
	m.previewMatches = nil
	m.previewMatch = 0
	query := strings.ToLower(strings.TrimSpace(m.previewQuery))
	if query == "" {
		return
	}
	for i, line := range m.currentDiffState().lines {
		if strings.Contains(strings.ToLower(line), query) {
			m.previewMatches = append(m.previewMatches, i)
		}
	}
	if len(m.previewMatches) > 0 {
		m.diffScroll = m.previewMatches[0]
	}
}

func (m workspaceModel) movePreviewMatch(direction int) workspaceModel {
	if len(m.previewMatches) == 0 {
		return m
	}
	m.previewMatch += direction
	if m.previewMatch < 0 {
		m.previewMatch = len(m.previewMatches) - 1
	}
	if m.previewMatch >= len(m.previewMatches) {
		m.previewMatch = 0
	}
	m.diffScroll = m.previewMatches[m.previewMatch]
	return m
}

func (m workspaceModel) setPreviewKind(kind app.PreviewKind) (workspaceModel, tea.Cmd) {
	if !m.isWorkspace() || m.previewKind == kind {
		return m, nil
	}
	m.previewKind = kind
	m.resetPreviewPosition()
	return m.startDiffLoad(false)
}

func (m workspaceModel) revealCurrentPreview() (workspaceModel, tea.Cmd) {
	if !m.isWorkspace() {
		return m, nil
	}
	state := m.currentDiffState()
	if !state.preview.Withheld {
		m.message = "current preview is not withheld"
		return m, nil
	}
	m.revealedPreviews[m.revealKey()] = true
	m.resetPreviewPosition()
	return m.startDiffLoad(true)
}

func workspacePreviewKey(target string, kind app.PreviewKind, reveal bool) string {
	return target + "\x00" + string(kind) + "\x00" + strconv.FormatBool(reveal)
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
