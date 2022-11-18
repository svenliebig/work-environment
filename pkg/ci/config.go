package ci

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/utils/wepath"
)

const ConfigFilename = "ci.json"

type CI struct {
	CiType     string
	Identifier string
	AuthToken  string
	Url        string
}

type CIConfig struct {
	Environments []*CI
}

func (C *CIConfig) Contains(identifier string) bool {
	for _, e := range C.Environments {
		if e.Identifier == identifier {
			return true
		}
	}
	return false
}

// returns the environment directly if there is only one ci environment
// if there are multiple environment, the identifier is used to find the
// correct environment.
func (c *CIConfig) GetEnvironment(identifier string) (*CI, error) {
	if len(c.Environments) == 0 {
		return nil, ErrCIDoesNotExist
	} else if len(c.Environments) == 1 {
		return c.Environments[0], nil
	} else {
		if identifier == "" {
			return nil, ErrRequiredIdentifierCI
		}

		return c.GetEnvironmentByIdentifier(identifier)
	}
}

func (C *CIConfig) GetEnvironmentByIdentifier(identifier string) (*CI, error) {
	for _, e := range C.Environments {
		if e.Identifier == identifier {
			return e, nil
		}
	}
	return nil, ErrCIDoesNotExist
}

var (
	ErrConfigDoesNotExist = errors.New("the ci configuration does not exists in .work-environment")
	ErrCIDoesNotExist     = errors.New("there is no ci with that identifier in the configuration")
)

// TODO basically the same code in CI and CONFIG package
func GetConfig(p string) (*CIConfig, error) {
	wer, err := wepath.GetWorkEnvironmentRoot(p)

	if err != nil {
		return nil, err
	}

	cp := filepath.Join(wer, ConfigFilename)

	if wepath.Exists(cp) {
		return ReadConfig(cp)
	}

	return nil, ErrConfigDoesNotExist
}

// reads the ci configuration, needs a direct absolute path to the configuration
func ReadConfig(configPath string) (*CIConfig, error) {
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

	var config CIConfig

	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
