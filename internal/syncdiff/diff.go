package syncdiff

import (
	"bytes"
	"fmt"
	"os"

	internaldiff "github.com/rogpeppe/go-internal/diff"
)

const (
	DefaultMaxFileSize = 1 << 20
	binaryProbeSize    = 8192
)

type ContentSource interface {
	TargetContent(target string) ([]byte, error)
	LocalContent(target string) ([]byte, error)
}

type Differ struct {
	Source      ContentSource
	MaxFileSize int64
}

func (d Differ) Diff(target string) ([]byte, error) {
	base, err := d.Source.TargetContent(target)
	if err != nil {
		return nil, fmt.Errorf("read chezmoi target %s: %w", target, err)
	}
	maxFileSize := d.maxFileSize()
	if tooLarge(base, maxFileSize) {
		return []byte(fmt.Sprintf("file too large to diff: %s\n", target)), nil
	}
	local, err := d.Source.LocalContent(target)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read local target %s: %w", target, err)
	}
	return DiffBytes(target, base, local, maxFileSize), nil
}

func DiffBytes(target string, base, local []byte, maxFileSize int64) []byte {
	if maxFileSize <= 0 {
		maxFileSize = DefaultMaxFileSize
	}
	if tooLarge(base, maxFileSize) || tooLarge(local, maxFileSize) {
		return []byte(fmt.Sprintf("file too large to diff: %s\n", target))
	}
	if binary(base) || binary(local) {
		return []byte(fmt.Sprintf("binary file differs: %s\n", target))
	}
	out := internaldiff.Diff("chezmoi:"+target, base, "local:"+target, local)
	if len(out) == 0 {
		return []byte(fmt.Sprintf("no diff: %s\n", target))
	}
	return out
}

func (d Differ) maxFileSize() int64 {
	if d.MaxFileSize > 0 {
		return d.MaxFileSize
	}
	return DefaultMaxFileSize
}

func tooLarge(content []byte, maxFileSize int64) bool {
	return int64(len(content)) > maxFileSize
}

func binary(content []byte) bool {
	if len(content) > binaryProbeSize {
		content = content[:binaryProbeSize]
	}
	return bytes.IndexByte(content, 0) >= 0
}

type ChezmoiContentLoader struct {
	Client ChezmoiClient
}

type ChezmoiClient interface {
	Output(args ...string) ([]byte, error)
}

func (l ChezmoiContentLoader) TargetContent(target string) ([]byte, error) {
	return l.Client.Output("cat", target)
}

func (ChezmoiContentLoader) LocalContent(target string) ([]byte, error) {
	content, err := os.ReadFile(target)
	if err != nil {
		return nil, err
	}
	return content, nil
}
