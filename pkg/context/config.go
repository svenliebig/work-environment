package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/wepath"
)

func readConfig(p string) (*core.Configuration, error) {
	file, err := os.Open(p)

	// @Comm
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("%w: err opening config file", err)
	}

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("%w: err reading config file", err)
	}

	var config core.Configuration

	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getConfigurationPath(from string) (string, error) {
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

	return getConfiguration(from)
}

func getConfiguration(p string) (string, error) {
	x := filepath.Join(p, core.ConfigurationFileName)

	if wepath.Exists(x) {
		return x, nil
	}

	if wepath.IsRoot(p) {
		return "", wepath.ErrWorkEnvironmentNotFound
	}

	return getConfiguration(filepath.Dir(p))
}
