package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWorkspaceServiceIntegrationInventoryAndPreviews(t *testing.T) {
	h := newChezmoiIntegration(t)
	identity := filepath.Join(h.root, "age-identity.txt")
	h.run(t, "age-keygen", "--output", identity)
	recipient := strings.TrimSpace(string(h.output(t, "age-keygen", "--convert", identity)))
	if recipient == "" {
		t.Fatal("age-keygen returned an empty recipient")
	}
	writeIntegrationFile(t, h.config, fmt.Sprintf("encryption = \"age\"\n[age]\nidentity = %q\nrecipient = %q\n", identity, recipient), 0o600)

	selected := filepath.Join(h.destination, "workspace")
	if err := os.MkdirAll(selected, 0o755); err != nil {
		t.Fatalf("create selected workspace: %v", err)
	}
	if err := os.Chmod(selected, 0o755); err != nil {
		t.Fatalf("chmod selected workspace: %v", err)
	}
	clean := filepath.Join(selected, "clean")
	dirty := filepath.Join(selected, "dirty")
	template := filepath.Join(selected, "template")
	encrypted := filepath.Join(selected, "encrypted")
	link := filepath.Join(selected, "link")
	managedDirectory := filepath.Join(selected, "managed-directory")
	unmanaged := filepath.Join(selected, "unmanaged")
	nestedUnmanaged := filepath.Join(selected, "nested", "unmanaged")

	writeIntegrationFile(t, clean, "clean managed content\n", 0o644)
	h.run(t, "add", clean)
	cleanSource := strings.TrimSpace(string(h.output(t, "source-path", clean)))

	writeIntegrationFile(t, dirty, "managed source content\n", 0o644)
	h.run(t, "add", dirty)
	dirtySource := strings.TrimSpace(string(h.output(t, "source-path", dirty)))
	writeIntegrationFile(t, dirty, "dirty destination content\n", 0o644)

	writeIntegrationFile(t, template, "template=destination\n", 0o644)
	h.run(t, "add", "--template", template)
	templateSourcePath := strings.TrimSpace(string(h.output(t, "source-path", template)))
	writeIntegrationFile(t, templateSourcePath, "template={{ .chezmoi.os }}\n", 0o644)

	const encryptedContent = "age encrypted integration content\n"
	writeIntegrationFile(t, encrypted, encryptedContent, 0o644)
	h.run(t, "add", "--encrypt", encrypted)
	encryptedSourcePath := strings.TrimSpace(string(h.output(t, "source-path", encrypted)))

	linkTarget := filepath.Join("..", "symlink-referent")
	writeIntegrationFile(t, filepath.Join(h.destination, "symlink-referent"), "not the symlink content\n", 0o644)
	if err := os.Symlink(linkTarget, link); err != nil {
		t.Fatalf("create managed symlink: %v", err)
	}
	h.run(t, "add", link)
	linkSource := strings.TrimSpace(string(h.output(t, "source-path", link)))

	writeIntegrationFile(t, filepath.Join(managedDirectory, "kept"), "managed directory content\n", 0o644)
	if err := os.Chmod(managedDirectory, 0o755); err != nil {
		t.Fatalf("chmod managed directory: %v", err)
	}
	h.run(t, "add", managedDirectory)
	directorySource := strings.TrimSpace(string(h.output(t, "source-path", managedDirectory)))
	writeIntegrationFile(t, filepath.Join(h.source, "run_once_workspace-pending.sh"), "#!/bin/sh\nexit 0\n", 0o700)

	writeIntegrationFile(t, filepath.Join(h.source, ".chezmoiignore"), "workspace/ignored\n", 0o600)
	writeIntegrationFile(t, filepath.Join(h.source, "workspace", "ignored"), "ignored source content\n", 0o600)
	writeIntegrationFile(t, unmanaged, "unmanaged direct content\n", 0o600)
	writeIntegrationFile(t, nestedUnmanaged, "unmanaged nested content\n", 0o600)

	unscoped, err := h.service.Inventory(nil)
	if err != nil {
		t.Fatalf("unscoped Inventory returned error: %v", err)
	}
	if !strings.Contains(unscoped.Notice, "unmanaged discovery disabled") {
		t.Fatalf("unscoped Inventory notice does not explain that unmanaged discovery is disabled")
	}
	if workspaceIntegrationHasEntry(unscoped, "workspace/unmanaged") || workspaceIntegrationHasEntry(unscoped, "workspace/nested/unmanaged") {
		t.Fatal("unscoped Inventory included unmanaged candidates")
	}

	cleanEntry := workspaceIntegrationEntry(t, unscoped, "workspace/clean")
	workspaceIntegrationAssertEntry(t, cleanEntry, clean, cleanSource, FileClean, TargetFile)
	dirtyEntry := workspaceIntegrationEntry(t, unscoped, "workspace/dirty")
	workspaceIntegrationAssertEntry(t, dirtyEntry, dirty, dirtySource, FileDirty, TargetFile)
	templateEntry := workspaceIntegrationEntry(t, unscoped, "workspace/template")
	workspaceIntegrationAssertEntry(t, templateEntry, template, templateSourcePath, FileDirty, TargetFile)
	if !templateEntry.Template || templateEntry.Encrypted {
		t.Fatal("template inventory metadata is incorrect")
	}
	encryptedEntry := workspaceIntegrationEntry(t, unscoped, "workspace/encrypted")
	workspaceIntegrationAssertEntry(t, encryptedEntry, encrypted, encryptedSourcePath, FileUninspected, TargetFile)
	if !encryptedEntry.Encrypted || encryptedEntry.Template {
		t.Fatal("encrypted inventory metadata is incorrect")
	}
	linkEntry := workspaceIntegrationEntry(t, unscoped, "workspace/link")
	workspaceIntegrationAssertEntry(t, linkEntry, link, linkSource, FileClean, TargetSymlink)
	directoryEntry := workspaceIntegrationEntry(t, unscoped, "workspace/managed-directory")
	workspaceIntegrationAssertEntry(t, directoryEntry, managedDirectory, directorySource, FileClean, TargetDirectory)
	scriptEntry := workspaceIntegrationEntry(t, unscoped, "workspace-pending.sh")
	workspaceIntegrationAssertEntry(t, scriptEntry, filepath.Join(h.destination, "workspace-pending.sh"), filepath.Join(h.source, "run_once_workspace-pending.sh"), FileScript, TargetScript)
	ignoredEntry := workspaceIntegrationEntry(t, unscoped, "workspace/ignored")
	workspaceIntegrationAssertEntry(t, ignoredEntry, filepath.Join(selected, "ignored"), "", FileIgnored, TargetUnknown)

	scoped, err := h.service.Inventory([]string{selected})
	if err != nil {
		t.Fatalf("scoped Inventory returned error: %v", err)
	}
	if len(scoped.Scopes) != 1 || scoped.Scopes[0] != selected {
		t.Fatal("scoped Inventory did not retain the selected directory scope")
	}
	workspaceIntegrationAssertEntry(t, workspaceIntegrationEntry(t, scoped, "workspace/unmanaged"), unmanaged, "", FileUnmanaged, TargetFile)
	workspaceIntegrationAssertEntry(t, workspaceIntegrationEntry(t, scoped, "workspace/nested/unmanaged"), nestedUnmanaged, "", FileUnmanaged, TargetFile)

	cleanDiff, err := h.service.Preview(cleanEntry, PreviewDiff, false)
	if err != nil {
		t.Fatalf("clean diff Preview returned error: %v", err)
	}
	if cleanDiff.Content != "" || cleanDiff.Notice != "destination matches rendered target" {
		t.Fatal("clean diff Preview did not report a matching rendered target")
	}

	dirtyDiff, err := h.service.Preview(dirtyEntry, PreviewDiff, false)
	if err != nil {
		t.Fatalf("dirty diff Preview returned error: %v", err)
	}
	if !dirtyDiff.Review.Dirty || !strings.Contains(dirtyDiff.Content, "-managed source content") || !strings.Contains(dirtyDiff.Content, "+dirty destination content") {
		t.Fatal("dirty diff Preview did not present the target-to-destination direction")
	}
	destination, err := h.service.Preview(dirtyEntry, PreviewDestination, false)
	if err != nil {
		t.Fatalf("destination Preview returned error: %v", err)
	}
	if destination.Content != "dirty destination content\n" {
		t.Fatal("destination Preview did not return destination content")
	}

	templateTarget, err := h.service.Preview(templateEntry, PreviewTarget, false)
	if err != nil {
		t.Fatalf("withheld template target Preview returned error: %v", err)
	}
	if !templateTarget.Withheld || templateTarget.Content != "" {
		t.Fatal("template target Preview did not withhold rendered content")
	}
	templateTarget, err = h.service.Preview(templateEntry, PreviewTarget, true)
	if err != nil {
		t.Fatalf("revealed template target Preview returned error: %v", err)
	}
	if templateTarget.Withheld || !strings.Contains(templateTarget.Content, "template="+runtime.GOOS) || strings.Contains(templateTarget.Content, "{{ .chezmoi.os }}") {
		t.Fatal("revealed template target Preview was not rendered")
	}
	templateSource, err := h.service.Preview(templateEntry, PreviewSource, false)
	if err != nil {
		t.Fatalf("template source Preview returned error: %v", err)
	}
	if templateSource.Withheld || templateSource.Content != "template={{ .chezmoi.os }}\n" {
		t.Fatal("template source Preview did not expose the unrendered source")
	}

	encryptedTarget, err := h.service.Preview(encryptedEntry, PreviewTarget, false)
	if err != nil {
		t.Fatalf("withheld encrypted target Preview returned error: %v", err)
	}
	encryptedDiff, err := h.service.Preview(encryptedEntry, PreviewDiff, false)
	if err != nil {
		t.Fatalf("withheld encrypted diff Preview returned error: %v", err)
	}
	if !encryptedDiff.Withheld || encryptedDiff.Content != "" {
		t.Fatal("encrypted diff Preview did not remain uninspected before reveal")
	}
	if !encryptedTarget.Withheld || encryptedTarget.Content != "" {
		t.Fatal("encrypted target Preview did not withhold rendered content")
	}
	encryptedSource, err := h.service.Preview(encryptedEntry, PreviewSource, false)
	if err != nil {
		t.Fatalf("withheld encrypted source Preview returned error: %v", err)
	}
	if !encryptedSource.Withheld || encryptedSource.Content != "" {
		t.Fatal("encrypted source Preview did not withhold plaintext")
	}
	encryptedSource, err = h.service.Preview(encryptedEntry, PreviewSource, true)
	if err != nil {
		t.Fatalf("revealed encrypted source Preview returned error: %v", err)
	}
	if encryptedSource.Withheld || encryptedSource.Content != encryptedContent {
		t.Fatal("revealed encrypted source Preview did not return the decrypted content")
	}

	linkDestination, err := h.service.Preview(linkEntry, PreviewDestination, false)
	if err != nil {
		t.Fatalf("symlink destination Preview returned error: %v", err)
	}
	if strings.TrimSuffix(linkDestination.Content, "\n") != linkTarget {
		t.Fatal("symlink destination Preview did not return the link target")
	}
	linkRenderedTarget, err := h.service.Preview(linkEntry, PreviewTarget, false)
	if err != nil {
		t.Fatalf("symlink target Preview returned error: %v", err)
	}
	if strings.TrimSuffix(linkRenderedTarget.Content, "\n") != linkTarget {
		t.Fatal("symlink target Preview did not return the rendered link target")
	}
}

func TestWorkspaceInventoryDoesNotResolveSecretTemplatesUntilReveal(t *testing.T) {
	h := newChezmoiIntegration(t)
	marker := filepath.Join(h.root, "secret-provider-called")
	provider := filepath.Join(h.root, "fake-pass")
	writeIntegrationFile(t, provider, fmt.Sprintf("#!/bin/sh\nprintf called > %q\nprintf 'RESEARCH_ONLY_VALUE\\n'\n", marker), 0o700)
	writeIntegrationFile(t, h.config, fmt.Sprintf("[pass]\ncommand = %q\n", provider), 0o600)

	target := filepath.Join(h.destination, ".secret-preview")
	writeIntegrationFile(t, target, "local value\n", 0o600)
	writeIntegrationFile(t, filepath.Join(h.source, "dot_secret-preview.tmpl"), "{{ pass \"workspace-test\" }}\n", 0o600)

	snapshot, err := h.service.Inventory(nil)
	if err != nil {
		t.Fatalf("Inventory returned error: %v", err)
	}
	entry := workspaceIntegrationEntry(t, snapshot, ".secret-preview")
	if entry.State != FileUninspected || !entry.Template {
		t.Fatal("secret-backed template was not marked uninspected")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("inventory invoked the secret provider")
	}

	preview, err := h.service.Preview(entry, PreviewDiff, false)
	if err != nil {
		t.Fatalf("withheld Preview returned error: %v", err)
	}
	if !preview.Withheld || preview.Content != "" {
		t.Fatal("secret-backed diff was not withheld")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("withheld preview invoked the secret provider")
	}

	preview, err = h.service.Preview(entry, PreviewDiff, true)
	if err != nil {
		t.Fatalf("revealed Preview returned error: %v", err)
	}
	if preview.Withheld || preview.Entry.State != FileDirty || preview.Content == "" {
		t.Fatal("revealed preview did not inspect the secret-backed target")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("explicit reveal did not invoke the secret provider: %v", err)
	}
}

func workspaceIntegrationEntry(t *testing.T, snapshot WorkspaceSnapshot, relativePath string) WorkspaceEntry {
	t.Helper()
	for _, entry := range snapshot.Entries {
		if entry.RelativePath == relativePath {
			return entry
		}
	}
	t.Fatalf("Inventory entry %q is missing", relativePath)
	return WorkspaceEntry{}
}

func workspaceIntegrationHasEntry(snapshot WorkspaceSnapshot, relativePath string) bool {
	for _, entry := range snapshot.Entries {
		if entry.RelativePath == relativePath {
			return true
		}
	}
	return false
}

func workspaceIntegrationAssertEntry(t *testing.T, entry WorkspaceEntry, path, sourcePath string, state FileState, targetType TargetType) {
	t.Helper()
	if entry.Path != path || entry.SourcePath != sourcePath || entry.State != state || entry.Type != targetType {
		t.Fatalf("Inventory metadata is incorrect for %q", entry.RelativePath)
	}
}
