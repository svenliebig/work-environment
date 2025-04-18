package git

import (
	"path/filepath"
	"strings"
)

func GetIdentifier(path string) (string, error) {
	r, err := getRepository(path)

	if err != nil {
		return "", err
	}

	c, err := r.Config()

	if err != nil {
		return "", err
	}

	re, ok := c.Remotes["origin"]

	if ok {
		if len(re.URLs) > 0 {
			url := re.URLs[0]
			b := filepath.Base(url)
			return strings.Replace(b, ".git", "", 1), nil
		}
	}

	wt, err := r.Worktree()

	if err != nil {
		return "", err
	}

	rpath := wt.Filesystem.Root()

	return filepath.Base(rpath), nil
}
