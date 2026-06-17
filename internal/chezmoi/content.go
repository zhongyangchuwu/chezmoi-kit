package chezmoi

import (
	"io"
	"os"
)

type ContentLoader struct {
	Client Client
}

func (l ContentLoader) TargetContent(target string, limit int64) ([]byte, error) {
	return l.Client.OutputLimit(limit, "cat", target)
}

func (ContentLoader) LocalContent(target string, limit int64) ([]byte, error) {
	return readFileLimit(target, limit)
}

func readFileLimit(path string, limit int64) ([]byte, error) {
	if limit < 0 {
		limit = 0
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content := make([]byte, limit+1)
	n, err := file.Read(content)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return content[:n], nil
}
