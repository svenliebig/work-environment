package wepath

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrWorkEnvironmentNotFound = errors.New("work environment not found")
)

// takes an absolute path and tries to find a work environment
func GetWorkEnvironmentRoot(from string) (string, error) {
	if !path.IsAbs(from) {
		return "", errors.New("expected an absolute path")
	}

	fs, err := os.Stat(from)

	if os.IsNotExist(err) {
		return "", errors.New("expected path to exist")
	}

	if !fs.IsDir() {
		return "", errors.New("expected path to refer to a directory")
	}

	return getWe(from)
}

func getWe(p string) (string, error) {
	x := filepath.Join(p, WorkEnvironmentDirectory)

	if Exists(x) {
		return x, nil
	}

	if IsRoot(p) {
		return "", ErrWorkEnvironmentNotFound
	}

	return getWe(filepath.Dir(p))
}
