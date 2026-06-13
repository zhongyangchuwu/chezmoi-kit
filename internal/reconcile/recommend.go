package reconcile

import "github.com/zhongyangchuwu/cm/internal/chezmoi"

func Describe(chezmoi.StatusEntry) string {
	return "differs from chezmoi"
}
