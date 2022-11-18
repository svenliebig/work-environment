package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/utils/wepath"
)

func GetConfig(p string) (*WorkEnvironmentConfiguration, error) {
	wer, err := wepath.GetWorkEnvironmentRoot(p)

	if err != nil {
		return nil, err
	}

	cp := filepath.Join(wer, ConfigFilename)

	if wepath.Exists(cp) {
		c, err := ReadConfig(cp)

		if err != nil {
			return nil, err
		}

		return &WorkEnvironmentConfiguration{
			location:              cp,
			workEnvironmentConfig: c,
		}, nil
	}

	return nil, ErrConfigDoesNotExist
}

// reads the ci configuration, needs a direct absolute path to the configuration
func ReadConfig(configPath string) (*workEnvironmentConfig, error) {
	file, err := os.Open(configPath)
	// @Comm
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("%w: err opening config file", err)
	}

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("%w: err reading config file", err)
	}

	var config workEnvironmentConfig

	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
