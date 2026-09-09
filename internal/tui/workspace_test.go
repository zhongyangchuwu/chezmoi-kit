package tui

import (
	"errors"
	"io"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func TestWorkspaceRetainsCleanEntriesAndLoadsInitialPreview(t *testing.T) {
	clean := workspaceEntry("clean.toml", app.FileClean, app.TargetFile)
	service := &fakeWorkspaceService{}
	model := newWorkspaceModel(service, workspaceSnapshot(clean))

	current := model.current()
	if len(model.entries) != 1 || current.Path != clean.Path || current.RelativePath != clean.RelativePath || current.State != app.FileClean || current.Type != app.TargetFile {
		t.Fatalf("entries = %#v, want persisted clean entry %#v", model.entries, clean)
	}
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init command is nil, want initial preview load")
	}
	msg, ok := cmd().(workspacePreviewMsg)
	if !ok {
		t.Fatalf("Init message = %T, want workspacePreviewMsg", msg)
	}
	if len(service.previewCalls) != 1 {
		t.Fatalf("preview calls = %#v, want one", service.previewCalls)
	}
	call := service.previewCalls[0]
	if call.entry.Path != clean.Path || call.kind != app.PreviewDiff || call.reveal {
		t.Fatalf("preview call = %#v, want clean diff without reveal", call)
	}

	model, _ = model.applyWorkspacePreview(msg)
	state := model.currentDiffState()
	if state.loading || state.preview.Entry.Path != clean.Path || !strings.Contains(strings.Join(state.lines, "\n"), "diff preview") {
		t.Fatalf("preview state = %#v", state)
	}
}

func TestWorkspaceDirectoryBrowserBuildsVirtualAncestorsAndRestoresSelection(t *testing.T) {
	child := workspaceEntry(".config/app/config.toml", app.FileDirty, app.TargetFile)
	otherChild := workspaceEntry(".config/theme.toml", app.FileClean, app.TargetFile)
	outside := workspaceEntry(".zshrc", app.FileClean, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(child, otherChild, outside))

	if got := strings.Join(workspacePaths(model.entries), ","); got != ".config,.zshrc" {
		t.Fatalf("root entries = %q, want virtual directory and sibling", got)
	}
	if entry := model.current(); entry.Type != app.TargetDirectory || entry.State != app.FileDirty {
		t.Fatalf("virtual directory = %#v, want dirty directory", entry)
	}
	model = workspaceKeyUpdate(t, model, 'l')
	if model.currentDir != ".config" {
		t.Fatalf("enter config = directory:%q entries:%#v", model.currentDir, workspacePaths(model.entries))
	}
	if got := strings.Join(workspacePaths(model.entries), ","); got != ".config/app,.config/theme.toml" {
		t.Fatalf("config entries = %q", got)
	}
	model = workspaceKeyUpdate(t, model, 'l')
	if model.currentDir != ".config/app" {
		t.Fatalf("enter app = directory:%q", model.currentDir)
	}
	if got := strings.Join(workspacePaths(model.entries), ","); got != ".config/app/config.toml" {
		t.Fatalf("app entries = %q", got)
	}
	model = workspaceKeyUpdate(t, model, 'h')
	if model.currentDir != ".config" || model.currentTarget() != "/home/me/.config/app" {
		t.Fatalf("leave app = directory:%q target:%q", model.currentDir, model.currentTarget())
	}
	model = workspaceKeyUpdate(t, model, 'h')
	if model.currentDir != "" || model.currentTarget() != "/home/me/.config" {
		t.Fatalf("leave config = directory:%q target:%q", model.currentDir, model.currentTarget())
	}
}

func TestWorkspaceAncestorEntryBecomesAggregateDirectory(t *testing.T) {
	ancestor := workspaceEntry(".config", app.FileClean, app.TargetUnknown)
	child := workspaceEntry(".config/app.toml", app.FileDirty, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(ancestor, child))

	entry := model.current()
	if entry.Type != app.TargetDirectory || entry.State != app.FileDirty {
		t.Fatalf("ancestor entry = %#v, want dirty navigable directory", entry)
	}
	model.currentDir = ".config"
	model.filter = filterIgnored
	model.rebuildEntries("")
	parents := model.parentEntries()
	if len(parents) != 1 || parents[0].RelativePath != ".config" {
		t.Fatalf("filtered parent context = %#v, want current directory retained", parents)
	}
}

func TestWorkspaceStartsAtSingleDirectoryScope(t *testing.T) {
	entry := workspaceEntry(".config/app.toml", app.FileClean, app.TargetFile)
	snapshot := workspaceSnapshot(entry)
	snapshot.Scopes = []string{"/home/me/.config"}
	model := newWorkspaceModel(nil, snapshot)

	if model.currentDir != ".config" || model.currentTarget() != entry.Path {
		t.Fatalf("scoped start = directory:%q target:%q", model.currentDir, model.currentTarget())
	}
}

func TestWorkspaceStartsAtCommonDirectoryForMultipleScopes(t *testing.T) {
	nvim := workspaceEntry(".config/nvim/init.lua", app.FileClean, app.TargetFile)
	git := workspaceEntry(".config/git/config", app.FileClean, app.TargetFile)
	snapshot := workspaceSnapshot(nvim, git)
	snapshot.Scopes = []string{"/home/me/.config/nvim", "/home/me/.config/git"}
	model := newWorkspaceModel(nil, snapshot)

	if model.currentDir != ".config" || strings.Join(workspacePaths(model.entries), ",") != ".config/git,.config/nvim" {
		t.Fatalf("multi-scope start = directory:%q entries:%q", model.currentDir, strings.Join(workspacePaths(model.entries), ","))
	}
}

func TestWorkspaceDirectoryFilterAndGlobalSearchKeepResultsReachable(t *testing.T) {
	dirty := workspaceEntry(".config/app.toml", app.FileDirty, app.TargetFile)
	clean := workspaceEntry(".config/clean.toml", app.FileClean, app.TargetFile)
	ignored := workspaceEntry(".ignored", app.FileIgnored, app.TargetFile)
	service := &fakeWorkspaceService{}
	model := newWorkspaceModel(service, workspaceSnapshot(dirty, clean, ignored))

	model.filter = filterDirty
	model.rebuildEntries("")
	if got := strings.Join(workspacePaths(model.entries), ","); got != ".config" {
		t.Fatalf("dirty root entries = %q, want matching ancestor", got)
	}
	model.enterCurrentDirectory()
	if got := strings.Join(workspacePaths(model.entries), ","); got != ".config/app.toml" {
		t.Fatalf("dirty directory entries = %q", got)
	}

	model.filter = filterAll
	model.currentDir = ""
	model.rebuildEntries("")
	model = model.beginSearch()
	model = typeWorkspaceSearch(t, model, "app")
	if !model.searchResults || model.currentTarget() != dirty.Path {
		t.Fatalf("search result = results:%t target:%q", model.searchResults, model.currentTarget())
	}
	updated, load := model.updateSearch(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = updated.(workspaceModel)
	if load == nil {
		t.Fatal("accepted search did not load the located file preview")
	}
	model, _ = model.applyWorkspacePreview(load().(workspacePreviewMsg))
	if model.searchResults || model.currentDir != ".config" || model.currentTarget() != dirty.Path || len(service.previewCalls) != 1 {
		t.Fatalf("accepted search = results:%t directory:%q target:%q calls:%d", model.searchResults, model.currentDir, model.currentTarget(), len(service.previewCalls))
	}

	model.currentDir = ""
	model.rebuildEntries("")
	model = model.beginSearch()
	model = typeWorkspaceSearch(t, model, "app")
	updated, _ = model.updateSearch(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
	model = updated.(workspaceModel)
	if model.searchResults || model.currentDir != "" || model.currentTarget() != "/home/me/.config" {
		t.Fatalf("escaped search = results:%t directory:%q target:%q", model.searchResults, model.currentDir, model.currentTarget())
	}
	model.currentDir = ".config"
	model.rebuildEntries(dirty.Path)
	model = model.beginSearch()
	model = typeWorkspaceSearch(t, model, "missing")
	updated, _ = model.updateSearch(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = updated.(workspaceModel)
	if model.searchResults || model.currentDir != ".config" || model.currentTarget() != dirty.Path || model.message != "no path matches" {
		t.Fatalf("empty search = results:%t directory:%q target:%q message:%q", model.searchResults, model.currentDir, model.currentTarget(), model.message)
	}
}

func TestWorkspaceRejectsReconciliationActions(t *testing.T) {
	model := newWorkspaceModel(nil, workspaceSnapshot(workspaceEntry(".config/app.toml", app.FileDirty, app.TargetFile)))

	for _, kind := range []app.ActionKind{app.ActionAdd, app.ActionApply, app.ActionMerge} {
		model = model.togglePending(kind)
		if len(model.pending) != 0 {
			t.Fatalf("pending after %s = %#v, want no mutations", kind, model.pending)
		}
		if !strings.Contains(model.message, "read-only") {
			t.Fatalf("message after %s = %q, want read-only rejection", kind, model.message)
		}
	}
}

func TestWorkspaceDirectoryPreviewIsLocalSummary(t *testing.T) {
	child := workspaceEntry("./.config/app.toml/", app.FileDirty, app.TargetFile)
	service := &fakeWorkspaceService{}
	model := newWorkspaceModel(service, workspaceSnapshot(child))

	command := model.Init()
	if command == nil {
		t.Fatal("directory summary command is nil")
	}
	message, ok := command().(workspacePreviewMsg)
	if !ok {
		t.Fatalf("directory summary message = %T", message)
	}
	model, _ = model.applyWorkspacePreview(message)
	lines := strings.Join(model.currentDiffState().lines, "\n")
	if len(service.previewCalls) != 0 || !strings.Contains(lines, "view: directory summary") || !strings.Contains(lines, "direct entries: 1") || !strings.Contains(lines, "matching descendants: 1") {
		t.Fatalf("directory preview called service or omitted summary: calls=%#v lines=%#v", service.previewCalls, model.currentDiffState().lines)
	}

	model.filter = filterIgnored
	model.rebuildEntries(model.currentTarget())
	model.filter = filterDirty
	model.rebuildEntries(model.currentTarget())
	model, command = model.startWorkspacePreviewLoad(false)
	if command == nil {
		t.Fatal("filtered directory summary reused stale cache")
	}
	model, _ = model.applyWorkspacePreview(command().(workspacePreviewMsg))
	if !strings.Contains(strings.Join(model.currentDiffState().lines, "\n"), "direct entries: 1") {
		t.Fatalf("filtered directory summary = %#v", model.currentDiffState().lines)
	}
}

func TestWorkspaceSourceEditHandoffRefreshesAndClearsStalePreviewState(t *testing.T) {
	entry := workspaceEntry(".config/app.toml", app.FileClean, app.TargetFile)
	entry.SourcePath = "/home/me/source/dot_config/app.toml"
	updated := entry
	updated.State = app.FileDirty
	service := &fakeWorkspaceService{snapshot: app.WorkspaceSnapshot{
		Root:    "/home/me",
		Scopes:  []string{"/home/me/.config"},
		Entries: []app.WorkspaceEntry{updated},
	}}
	initial := workspaceSnapshot(entry)
	initial.Scopes = []string{"/home/me/.config"}
	model := newWorkspaceModel(service, initial)
	model.currentDir = ".config"
	model.rebuildEntries(entry.Path)
	model.filter = filterManaged
	model.focus = focusDiff
	model.previewKind = app.PreviewTarget
	model.previewFull = true
	model.diffs[model.previewStateKey()] = diffState{lines: []string{"stale preview"}}
	model.revealedPreviews[model.revealKey()] = true
	model.previewQuery = "stale"
	model.previewMatches = []int{0}
	model.diffScroll, model.previewX = 3, 2
	stalePreview := workspacePreviewMsg{
		key:     model.previewStateKey(),
		epoch:   model.previewEpoch,
		preview: app.WorkspacePreview{Entry: entry, Kind: model.previewKind, Content: "stale response"},
	}

	updatedModel, prepare := model.Update(tea.KeyPressMsg(tea.Key{Code: 'e'}))
	model = updatedModel.(workspaceModel)
	if !model.workspaceBusy || prepare == nil || model.message != "opening source editor..." {
		t.Fatalf("source edit start = busy:%t command:%t message:%q", model.workspaceBusy, prepare != nil, model.message)
	}
	lockedCursor := model.cursor
	locked, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'j'}))
	model = locked.(workspaceModel)
	if model.cursor != lockedCursor {
		t.Fatalf("editor handoff allowed input: cursor=%d want %d", model.cursor, lockedCursor)
	}
	request, ok := prepare().(workspaceHandoffRequestMsg)
	if !ok || len(service.sourceEditCalls) != 1 || service.sourceEditCalls[0].Path != entry.Path {
		t.Fatalf("source edit request = %#v calls=%#v", request, service.sourceEditCalls)
	}
	model, _ = model.applyWorkspaceHandoffRequest(request)
	model, refresh := model.applyWorkspaceHandoffDone(workspaceHandoffDoneMsg{target: entry.Path})
	if refresh == nil || model.message != "refreshing workspace..." {
		t.Fatalf("handoff return = refresh:%t message:%q", refresh != nil, model.message)
	}
	model, preview := model.applyWorkspaceRefresh(refresh().(workspaceRefreshMsg))
	model, _ = model.applyWorkspacePreview(stalePreview)
	if lines := model.currentDiffState().lines; len(lines) != 0 {
		t.Fatalf("stale preview response survived refresh: %q", lines)
	}
	if model.workspaceBusy || model.currentTarget() != updated.Path || model.currentDir != ".config" || model.filter != filterManaged || model.focus != focusDiff || model.previewKind != app.PreviewTarget || !model.previewFull {
		t.Fatalf("refreshed workspace state = %#v", model)
	}
	if len(model.diffs) != 1 || len(model.revealedPreviews) != 0 || model.previewQuery != "" || len(model.previewMatches) != 0 || model.diffScroll != 0 || model.previewX != 0 {
		t.Fatalf("stale preview state survived refresh: diffs=%#v revealed=%#v query=%q matches=%#v scroll=%d x=%d", model.diffs, model.revealedPreviews, model.previewQuery, model.previewMatches, model.diffScroll, model.previewX)
	}
	if preview == nil {
		t.Fatal("refresh did not reload selected preview")
	}
	model, _ = model.applyWorkspacePreview(preview().(workspacePreviewMsg))
	if !strings.Contains(model.workspaceNotice, "workspace refreshed") || !strings.Contains(strings.Join(model.currentDiffState().lines, "\n"), "target preview") {
		t.Fatalf("refreshed preview = notice:%q lines=%#v", model.workspaceNotice, model.currentDiffState().lines)
	}
	if got := service.inventoryCalls; len(got) != 1 || strings.Join(got[0], ",") != "/home/me/.config" {
		t.Fatalf("inventory calls = %#v", got)
	}
}

func TestWorkspaceSourceEditRejectsUnsupportedEntryAndRefreshesAfterEditorError(t *testing.T) {
	unmanaged := workspaceEntry("new.conf", app.FileUnmanaged, app.TargetFile)
	model := newWorkspaceModel(&fakeWorkspaceService{}, workspaceSnapshot(unmanaged))
	model = workspaceKeyUpdate(t, model, 'e')
	if !strings.Contains(model.message, "unmanaged target has no chezmoi source") {
		t.Fatalf("unmanaged source edit message = %q", model.message)
	}

	entry := workspaceEntry("app.toml", app.FileClean, app.TargetFile)
	entry.SourcePath = "/home/me/source/dot_app.toml"
	service := &fakeWorkspaceService{snapshot: workspaceSnapshot(entry)}
	model = newWorkspaceModel(service, workspaceSnapshot(entry))
	model, prepare := model.startWorkspaceSourceEdit()
	request := prepare().(workspaceHandoffRequestMsg)
	model, _ = model.applyWorkspaceHandoffRequest(request)
	model, refresh := model.applyWorkspaceHandoffDone(workspaceHandoffDoneMsg{target: entry.Path, err: errors.New("editor cancelled")})
	model, _ = model.applyWorkspaceRefresh(refresh().(workspaceRefreshMsg))
	if model.workspaceBusy || !strings.Contains(model.workspaceNotice, "editor exited with error: editor cancelled") || !strings.Contains(model.workspaceNotice, "workspace refreshed") {
		t.Fatalf("editor error refresh = busy:%t notice:%q", model.workspaceBusy, model.workspaceNotice)
	}
}

func TestWorkspaceSourceEditRefreshFailureKeepsBrowsableSnapshot(t *testing.T) {
	entry := workspaceEntry("app.toml", app.FileClean, app.TargetFile)
	entry.SourcePath = "/home/me/source/dot_app.toml"
	service := &fakeWorkspaceService{inventoryErr: errors.New("inventory unavailable")}
	model := newWorkspaceModel(service, workspaceSnapshot(entry))
	model.diffs[model.previewStateKey()] = diffState{lines: []string{"stale"}}
	model.revealedPreviews[model.revealKey()] = true

	model, refresh := model.applyWorkspaceHandoffDone(workspaceHandoffDoneMsg{target: entry.Path})
	model, _ = model.applyWorkspaceRefresh(refresh().(workspaceRefreshMsg))
	if model.workspaceBusy || len(model.entries) != 1 || len(model.diffs) != 0 || len(model.revealedPreviews) != 0 || !strings.Contains(model.message, "workspace refresh failed: inventory unavailable") {
		t.Fatalf("refresh failure state = entries:%#v diffs:%#v revealed:%#v message:%q", model.entries, model.diffs, model.revealedPreviews, model.message)
	}
}
func TestWorkspaceSourceEditRefreshUsesSelectionFallbackAndContextualFooter(t *testing.T) {
	oldEntry := workspaceEntry(".config/old.toml", app.FileClean, app.TargetFile)
	oldEntry.SourcePath = "/home/me/source/dot_config/old.toml"
	newEntry := workspaceEntry(".config/new.toml", app.FileDirty, app.TargetFile)
	newEntry.SourcePath = "/home/me/source/dot_config/new.toml"
	service := &fakeWorkspaceService{snapshot: workspaceSnapshot(newEntry)}
	model := newWorkspaceModel(service, workspaceSnapshot(oldEntry))
	model.currentDir = ".config"
	model.rebuildEntries(oldEntry.Path)

	model, refresh := model.applyWorkspaceHandoffDone(workspaceHandoffDoneMsg{target: oldEntry.Path})
	model, _ = model.applyWorkspaceRefresh(refresh().(workspaceRefreshMsg))
	if model.currentDir != ".config" || model.currentTarget() != newEntry.Path {
		t.Fatalf("selection fallback = directory:%q target:%q", model.currentDir, model.currentTarget())
	}
	if !strings.Contains(ansi.Strip(model.workspaceFooter()), "e edit source") {
		t.Fatalf("eligible footer = %q", ansi.Strip(model.workspaceFooter()))
	}

	unmanaged := workspaceEntry("new.conf", app.FileUnmanaged, app.TargetFile)
	model = newWorkspaceModel(&fakeWorkspaceService{}, workspaceSnapshot(unmanaged))
	if strings.Contains(ansi.Strip(model.workspaceFooter()), "e edit source") {
		t.Fatalf("unmanaged footer advertises source edit: %q", ansi.Strip(model.workspaceFooter()))
	}
}

func TestWorkspaceDirectoryRowsAlignAndResponsivePanes(t *testing.T) {
	directory := workspaceEntry(".config", app.FileClean, app.TargetDirectory)
	file := workspaceEntry(".zshrc", app.FileClean, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(directory, file))

	labelDirectory := ansi.Strip(model.entryLabel(directory))
	labelFile := ansi.Strip(model.entryLabel(file))
	if !strings.HasSuffix(labelDirectory, ".config/") || strings.Contains(labelDirectory, "▾") || strings.Contains(labelDirectory, "▸") {
		t.Fatalf("directory label = %q", labelDirectory)
	}
	if strings.Index("> "+labelDirectory, "C:") != strings.Index("  "+labelFile, "C:") {
		t.Fatalf("state columns are misaligned: directory=%q file=%q", labelDirectory, labelFile)
	}

	model.width, model.height = 120, 24
	wide := ansi.Strip(model.viewString())
	for _, pane := range []string{"Parent", "Current", "Preview"} {
		if !strings.Contains(wide, pane) {
			t.Fatalf("wide browser omits %q: %q", pane, wide)
		}
	}
	model.width = 80
	medium := ansi.Strip(model.viewString())
	if strings.Contains(medium, "Parent") || !strings.Contains(medium, "Current") || !strings.Contains(medium, "Preview") {
		t.Fatalf("medium browser panes = %q", medium)
	}
}

func TestWorkspacePreviewKindSwitchLoadsSelectedEntry(t *testing.T) {
	entry := workspaceEntry("app.toml", app.FileDirty, app.TargetFile)
	entry.SourcePath = "/home/me/.local/share/chezmoi/dot_config/app.toml"
	service := &fakeWorkspaceService{}
	model := newWorkspaceModel(service, workspaceSnapshot(entry))

	for _, kind := range []app.PreviewKind{app.PreviewDestination, app.PreviewTarget, app.PreviewSource, app.PreviewDiff} {
		var cmd tea.Cmd
		model, cmd = model.setPreviewKind(kind)
		if cmd == nil {
			t.Fatalf("switch to %q returned nil preview command", kind)
		}
		msg, ok := cmd().(workspacePreviewMsg)
		if !ok {
			t.Fatalf("switch to %q message = %T, want workspacePreviewMsg", kind, msg)
		}
		call := service.previewCalls[len(service.previewCalls)-1]
		if call.entry.Path != entry.Path || call.entry.SourcePath != entry.SourcePath || call.kind != kind || call.reveal {
			t.Fatalf("switch to %q preview call = %#v", kind, call)
		}
		model, _ = model.applyWorkspacePreview(msg)
		if got := model.currentDiffState().preview.Kind; got != kind {
			t.Fatalf("preview kind after switch = %q, want %q", got, kind)
		}
	}
}

func TestWorkspaceIgnoresStalePreviewMessagesForActiveViewState(t *testing.T) {
	first := workspaceEntry(".first", app.FileClean, app.TargetFile)
	second := workspaceEntry(".second", app.FileClean, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(first, second))
	model.cursor = 1
	model.previewQuery = "needle"
	model.previewMatches = []int{3}
	model.previewMatch = 0
	model.diffScroll = 3
	model.message = "active preview status"
	model.diffs[model.previewStateKey()] = diffState{lines: []string{"needle"}}

	stale := workspacePreviewMsg{
		key: workspacePreviewKey(first.Path, app.PreviewDiff, false),
		preview: app.WorkspacePreview{
			Entry:   first,
			Kind:    app.PreviewDiff,
			Content: "stale preview",
		},
	}
	model, _ = model.applyWorkspacePreview(stale)

	if model.message != "active preview status" || model.diffScroll != 3 || len(model.previewMatches) != 1 || model.previewMatches[0] != 3 {
		t.Fatalf("stale preview disturbed active view state: message=%q scroll=%d matches=%#v", model.message, model.diffScroll, model.previewMatches)
	}
	if got := model.diffs[stale.key].preview.Content; got != "stale preview" {
		t.Fatalf("stale preview was not cached: %q", got)
	}
}

func TestWorkspaceStatusLineShowsCommittedPathFilter(t *testing.T) {
	model := newWorkspaceModel(nil, workspaceSnapshot(workspaceEntry(".config/app.toml", app.FileClean, app.TargetFile)))
	model.fileQuery = "app"
	if got := model.statusLine(); got != "path filter: app" {
		t.Fatalf("statusLine() = %q, want committed path filter", got)
	}
}

func TestWorkspaceWithheldPreviewRequiresExplicitReveal(t *testing.T) {
	entry := workspaceEntry(".secrets.env", app.FileClean, app.TargetFile)
	service := &fakeWorkspaceService{previews: map[workspacePreviewResponseKey]app.WorkspacePreview{
		{target: entry.Path, kind: app.PreviewDiff, reveal: false}: {
			Withheld: true,
			Notice:   "preview withheld",
		},
		{target: entry.Path, kind: app.PreviewDiff, reveal: true}: {
			Content: "revealed content",
		},
	}}
	model := newWorkspaceModel(service, workspaceSnapshot(entry))

	initial := model.Init()
	if initial == nil {
		t.Fatal("Init command is nil, want withheld preview load")
	}
	model, _ = model.applyWorkspacePreview(initial().(workspacePreviewMsg))
	if !model.currentDiffState().preview.Withheld {
		t.Fatalf("initial preview = %#v, want withheld", model.currentDiffState().preview)
	}

	var reveal tea.Cmd
	model, reveal = model.revealCurrentPreview()
	if reveal == nil {
		t.Fatal("reveal command is nil")
	}
	model, _ = model.applyWorkspacePreview(reveal().(workspacePreviewMsg))
	if !model.currentPreviewRevealed() || model.currentDiffState().preview.Withheld || model.currentDiffState().preview.Content != "revealed content" {
		t.Fatalf("revealed preview state = %#v", model.currentDiffState())
	}
	if len(service.previewCalls) != 2 || service.previewCalls[0].reveal || !service.previewCalls[1].reveal {
		t.Fatalf("preview calls = %#v, want withheld then explicit reveal", service.previewCalls)
	}
}

func TestWorkspacePreviewNavigationReachesLineTailsHunksAndMatches(t *testing.T) {
	entry := workspaceEntry("app.toml", app.FileDirty, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(entry))
	longLine := strings.Repeat("x", 20) + "THE-TAIL"
	model.diffs[model.previewStateKey()] = diffState{lines: []string{
		"header",
		"@@ first hunk",
		"needle one",
		"@@ second hunk",
		"needle two",
		longLine,
	}}
	model.focus = focusDiff

	model = model.scrollPreviewHorizontal(100, 8)
	if want := lipgloss.Width(longLine) - 8; model.previewX != want {
		t.Fatalf("horizontal scroll = %d, want tail offset %d", model.previewX, want)
	}
	if rendered := model.renderDiffPane(rect{width: 10, height: 9}); !strings.Contains(rendered, "THE-TAIL") {
		t.Fatalf("tail is not visible after horizontal scroll: %q", rendered)
	}

	model = workspaceKeyUpdate(t, model, ']')
	if model.diffScroll != 1 {
		t.Fatalf("next hunk scroll = %d, want 1", model.diffScroll)
	}
	model = workspaceKeyUpdate(t, model, ']')
	if model.diffScroll != 3 {
		t.Fatalf("second next hunk scroll = %d, want 3", model.diffScroll)
	}
	model = workspaceKeyUpdate(t, model, '[')
	if model.diffScroll != 1 {
		t.Fatalf("previous hunk scroll = %d, want 1", model.diffScroll)
	}

	model = model.beginSearch()
	model = typeWorkspaceSearch(t, model, "needle")
	if model.diffScroll != 2 {
		t.Fatalf("first preview match scroll = %d, want 2", model.diffScroll)
	}
	updated, _ := model.updateSearch(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = updated.(workspaceModel)
	model = workspaceKeyUpdate(t, model, 'n')
	if model.diffScroll != 4 {
		t.Fatalf("next preview match scroll = %d, want 4", model.diffScroll)
	}
	model = workspaceKeyUpdate(t, model, 'N')
	if model.diffScroll != 2 {
		t.Fatalf("previous preview match scroll = %d, want 2", model.diffScroll)
	}
}
func TestWorkspaceHelpShowsGuidanceAndBlocksWorkspaceActions(t *testing.T) {
	first := workspaceEntry(".config/first.toml", app.FileClean, app.TargetFile)
	second := workspaceEntry(".config/second.toml", app.FileDirty, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(first, second))
	model.width = 100
	model.height = 30

	model = workspaceKeyUpdate(t, model, '?')
	if !model.helpVisible {
		t.Fatal("help is not visible after ?")
	}
	for _, text := range []string{"Quick help", "Start here", "C clean", "[E] encrypted", "Esc, ?, or q closes help"} {
		if !strings.Contains(ansi.Strip(model.viewString()), text) {
			t.Fatalf("help does not contain %q", text)
		}
	}
	cursor, filter := model.cursor, model.filter
	model = workspaceKeyUpdate(t, model, 'j')
	model = workspaceKeyUpdate(t, model, 'f')
	if model.cursor != cursor || model.filter != filter {
		t.Fatalf("help allowed workspace mutation: cursor=%d filter=%v", model.cursor, model.filter)
	}

	model = workspaceKeyUpdate(t, model, '?')
	if model.helpVisible {
		t.Fatal("help remains visible after ?")
	}
	model = workspaceKeyUpdate(t, model, '?')
	model = workspaceKeyUpdate(t, model, 'q')
	if model.helpVisible || model.stopped {
		t.Fatal("q did not close help without quitting workspace")
	}
}

func TestWorkspaceSemanticLegendSurvivesNoColorMode(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	model := newWorkspaceModel(nil, workspaceSnapshot(workspaceEntry(".config/app.toml", app.FileDirty, app.TargetFile)))
	legend := ansi.Strip(model.workspaceLegend())
	for _, text := range []string{"C clean", "D dirty", "U unmanaged", "I ignored", "R script", "? uninspected", "[T] template", "[E] encrypted"} {
		if !strings.Contains(legend, text) {
			t.Fatalf("no-color legend does not contain %q: %q", text, legend)
		}
	}
}

func TestWorkspacePaletteDistinguishesFileStates(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	styles := newTUIStyles()
	rendered := map[app.FileState]string{}
	for _, state := range []app.FileState{app.FileClean, app.FileDirty, app.FileUnmanaged, app.FileIgnored, app.FileScript, app.FileUninspected} {
		rendered[state] = styles.stateStyle(state).Render("x")
	}
	for state, value := range rendered {
		for otherState, otherValue := range rendered {
			if state != otherState && value == otherValue {
				t.Fatalf("states %q and %q have identical visual markers", state, otherState)
			}
		}
	}
}

func TestWorkspaceFooterRetainsHelpAndQuitAtEverySupportedWidth(t *testing.T) {
	model := newWorkspaceModel(nil, workspaceSnapshot(workspaceEntry(".config/app.toml", app.FileClean, app.TargetFile)))
	for _, width := range []int{24, 60, 80, 120} {
		model.width = width
		footer := ansi.Strip(model.workspaceFooter())
		if !strings.Contains(footer, "?") || !strings.Contains(footer, "q") || lipgloss.Width(footer) > width {
			t.Fatalf("footer at width %d = %q width=%d", width, footer, lipgloss.Width(footer))
		}
	}
}

func TestWorkspaceFullPreviewOmitsFilesAndFitsNarrowWidth(t *testing.T) {
	entry := workspaceEntry("workspace-file.txt", app.FileClean, app.TargetFile)
	model := newWorkspaceModel(nil, workspaceSnapshot(entry))
	model.width = 80
	model.height = 14
	model.diffs[model.previewStateKey()] = diffState{lines: []string{"preview payload"}}

	if rendered := ansi.Strip(model.viewString()); !strings.Contains(rendered, "C:f workspace-file.txt") {
		t.Fatalf("normal workspace view omits file pane entry: %q", rendered)
	}
	model.previewFull = true
	if !model.previewFull {
		t.Fatal("full preview flag is false")
	}
	if rendered := ansi.Strip(model.viewString()); strings.Contains(rendered, "C:f workspace-file.txt") {
		t.Fatalf("full preview still renders file pane entry: %q", rendered)
	}

	model.width = 24
	for _, line := range strings.Split(model.viewString(), "\n") {
		if width := lipgloss.Width(line); width > model.width {
			t.Fatalf("narrow view line width = %d, want at most %d: %q", width, model.width, line)
		}
	}
}

func workspaceSnapshot(entries ...app.WorkspaceEntry) app.WorkspaceSnapshot {
	return app.WorkspaceSnapshot{Root: "/home/me", Entries: entries}
}

func workspaceEntry(relativePath string, state app.FileState, targetType app.TargetType) app.WorkspaceEntry {
	return app.WorkspaceEntry{
		Path:         "/home/me/" + relativePath,
		RelativePath: relativePath,
		State:        state,
		Type:         targetType,
	}
}

func workspacePaths(entries []app.WorkspaceEntry) []string {
	paths := make([]string, len(entries))
	for i, entry := range entries {
		paths[i] = entry.RelativePath
	}
	return paths
}

func typeWorkspaceSearch(t *testing.T, model workspaceModel, input string) workspaceModel {
	t.Helper()
	for _, r := range input {
		updated, _ := model.updateSearch(tea.KeyPressMsg(tea.Key{Text: string(r)}))
		model = updated.(workspaceModel)
	}
	return model
}

func workspaceKeyUpdate(t *testing.T, model workspaceModel, code rune) workspaceModel {
	t.Helper()
	updated, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: code}))
	return updated.(workspaceModel)
}

type workspacePreviewResponseKey struct {
	target string
	kind   app.PreviewKind
	reveal bool
}

type workspacePreviewCall struct {
	entry  app.WorkspaceEntry
	kind   app.PreviewKind
	reveal bool
}

type fakeWorkspaceService struct {
	snapshot        app.WorkspaceSnapshot
	inventoryErr    error
	inventoryCalls  [][]string
	previews        map[workspacePreviewResponseKey]app.WorkspacePreview
	previewCalls    []workspacePreviewCall
	sourceEditErr   error
	sourceEditCalls []app.WorkspaceEntry
	sourceEditCmd   app.TerminalCommand
}

func (f *fakeWorkspaceService) Status([]string) (app.SyncStatus, error) {
	return app.SyncStatus{}, nil
}

func (f *fakeWorkspaceService) Review(target string) (app.Review, error) {
	return app.Review{Entry: app.ReconcileEntry{Path: target}}, nil
}

func (f *fakeWorkspaceService) ExecuteNonInteractive(app.Action) (app.ActionResult, error) {
	return app.ActionResult{}, nil
}

func (f *fakeWorkspaceService) TerminalCommand(app.Action) (app.TerminalCommand, error) {
	return workspaceTerminalCommand{}, nil
}

func (f *fakeWorkspaceService) Inventory(scopes []string) (app.WorkspaceSnapshot, error) {
	f.inventoryCalls = append(f.inventoryCalls, append([]string(nil), scopes...))
	if f.inventoryErr != nil {
		return app.WorkspaceSnapshot{}, f.inventoryErr
	}
	return f.snapshot, nil
}

func (f *fakeWorkspaceService) SourceEditCommand(entry app.WorkspaceEntry) (app.TerminalCommand, error) {
	f.sourceEditCalls = append(f.sourceEditCalls, entry)
	if f.sourceEditErr != nil {
		return nil, f.sourceEditErr
	}
	if f.sourceEditCmd != nil {
		return f.sourceEditCmd, nil
	}
	return workspaceTerminalCommand{}, nil
}

func (f *fakeWorkspaceService) Preview(entry app.WorkspaceEntry, kind app.PreviewKind, reveal bool) (app.WorkspacePreview, error) {
	f.previewCalls = append(f.previewCalls, workspacePreviewCall{entry: entry, kind: kind, reveal: reveal})
	if preview, ok := f.previews[workspacePreviewResponseKey{target: entry.Path, kind: kind, reveal: reveal}]; ok {
		preview.Entry = entry
		preview.Kind = kind
		return preview, nil
	}
	return app.WorkspacePreview{Entry: entry, Kind: kind, Content: string(kind) + " preview"}, nil
}

type workspaceTerminalCommand struct{}

func (workspaceTerminalCommand) Run() error          { return nil }
func (workspaceTerminalCommand) SetStdin(io.Reader)  {}
func (workspaceTerminalCommand) SetStdout(io.Writer) {}
func (workspaceTerminalCommand) SetStderr(io.Writer) {}
