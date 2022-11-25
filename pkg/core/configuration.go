package core

import (
	"errors"
)

const ConfigurationFileName = ".work-environment.json"

type ConfigurationNew struct {
	Projects       []*Project
	CIEnvironments []*CI
}

type Configuration interface {
	Write() error
}

// returns the ci environment with the given identifier, if it exists
func (c *ConfigurationNew) GetCIEnvironmentById(identifier string) (*CI, error) {
	if len(c.CIEnvironments) == 0 {
		return nil, errors.New("no ci environments defined")
	} else {
		if identifier == "" {
			return nil, errors.New("parameter identifier is empty, but required when there are more then one ci environment")
		}

		for _, e := range c.CIEnvironments {
			if e.Identifier == identifier {
				return e, nil
			}
		}
		return nil, errors.New("no such ci environment with the given identifier")
	}
}

func (c *ConfigurationNew) HasCI(id string) bool {
	for _, ci := range c.CIEnvironments {
		if ci.Identifier == id {
			return true
		}
	}
	return false
}

// Put this into context?
func (c *ConfigurationNew) UpdateProjectCI(identifier string, pci *ProjectCI) error {
	for _, p := range c.Projects {
		if p.Identifier == identifier {
			p.CI = pci
			return nil
		}
	}

	return errors.New("project not found")
}
