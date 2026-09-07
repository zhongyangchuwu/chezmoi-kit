package tui

import (
	"os"
	"path/filepath"
	"strings"
)

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Clean(home)
}

func (m workspaceModel) displayPath(path string) string {
	if path == "" || m.homeDir == "" {
		return path
	}
	clean := filepath.Clean(path)
	if clean == m.homeDir {
		return "~"
	}
	prefix := m.homeDir + string(os.PathSeparator)
	if strings.HasPrefix(clean, prefix) {
		return "~" + string(os.PathSeparator) + strings.TrimPrefix(clean, prefix)
	}
	return path
}
