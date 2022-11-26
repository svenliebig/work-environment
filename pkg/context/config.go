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

func (c *Context) GetConfiguration() (*core.Configuration, error) {
	if c.configuration == nil {
		p, err := getConfigurationPath(c.Cwd)

		if err != nil {
			return nil, err
		}

		config, err := readConfig(p)

		if err != nil {
			return nil, err
		}

		c.configurationPath = p
		c.configuration = config
	}

	return c.configuration, nil
}

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

func (c *Context) UpdateConfig() error {
	result, err := json.MarshalIndent(c.configuration, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(c.configurationPath, result, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (c *Context) AddCIEnvironmentToConfiguration(ci *core.CI) error {
	if c.configuration.CIEnvironments == nil {
		c.configuration.CIEnvironments = []*core.CI{ci}
	} else {
		for _, cie := range c.configuration.CIEnvironments {
			if cie.Identifier == ci.Identifier {
				return ErrCIAlreadyExists
			}
		}

		c.configuration.CIEnvironments = append(c.configuration.CIEnvironments, ci)
	}

	return nil
}

func (c *Context) UpdateCIEnvironmentToConfiguration(ci *core.CI) error {
	for i, cie := range c.configuration.CIEnvironments {
		if cie.Identifier == ci.Identifier {
			c.configuration.CIEnvironments[i] = ci
			return nil
		}
	}

	return ErrNoSuchCI
}
