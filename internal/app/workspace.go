package app

// FileState separates membership from whether a destination was actually inspected.
type FileState string

const (
	FileClean       FileState = "clean"
	FileDirty       FileState = "dirty"
	FileUnmanaged   FileState = "unmanaged"
	FileIgnored     FileState = "ignored"
	FileUninspected FileState = "uninspected"
	FileScript      FileState = "script"
)

// WorkspaceEntry identifies a target across views without encoding chezmoi filenames.
type WorkspaceEntry struct {
	Path         string
	RelativePath string
	SourcePath   string
	State        FileState
	Type         TargetType
	Template     bool
	Encrypted    bool
	Code         string
}

// WorkspaceSnapshot contains sorted entries and the explicit discovery scope.
// Empty Scopes means managed/source-ignored inventory only, not a HOME scan.
type WorkspaceSnapshot struct {
	Root    string
	Scopes  []string
	Entries []WorkspaceEntry
	Notice  string
}

type PreviewKind string

const (
	PreviewDiff        PreviewKind = "diff"
	PreviewDestination PreviewKind = "destination"
	PreviewTarget      PreviewKind = "target"
	PreviewSource      PreviewKind = "source"
)

// WorkspacePreview never represents withheld content as a clean review.
// Review is actionable only when a complete authoritative diff has been inspected.
type WorkspacePreview struct {
	Entry    WorkspaceEntry
	Kind     PreviewKind
	Content  string
	Notice   string
	Withheld bool
	Review   Review
}

// WorkspaceService adds read-only browsing to the existing reconciliation service.
// reveal is an explicit, per-target user decision, not an automatic fallback.
type WorkspaceService interface {
	SyncService
	Inventory(scopes []string) (WorkspaceSnapshot, error)
	Preview(entry WorkspaceEntry, kind PreviewKind, reveal bool) (WorkspacePreview, error)
}
