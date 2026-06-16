package chezmoi

import "os"

type ContentLoader struct {
	Client Client
}

func (l ContentLoader) TargetContent(target string) ([]byte, error) {
	return l.Client.Output("cat", target)
}

func (ContentLoader) LocalContent(target string) ([]byte, error) {
	content, err := os.ReadFile(target)
	if err != nil {
		return nil, err
	}
	return content, nil
}
