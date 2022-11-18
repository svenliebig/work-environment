package config

import (
	"errors"
	"strings"

	"github.com/svenliebig/work-environment/pkg/core"
)

const ConfigFilename = "config.json"

var (
	_ core.Configuration = &WorkEnvironmentConfiguration{}

	// errors
	ErrConfigDoesNotExist = errors.New("the configuration does not exists in .work-environment")
)

type workEnvironmentConfig struct {
	Projects []*core.Project
}

type WorkEnvironmentConfiguration struct {
	*workEnvironmentConfig
	location string
}

func (c *WorkEnvironmentConfiguration) GetProjectByPath(p string) (*core.Project, error) {
	for _, project := range c.Projects {
		if strings.Contains(p, project.Path) {
			return project, nil
		}
	}
	return nil, ErrConfigDoesNotExist
}

func (c *WorkEnvironmentConfiguration) Write() error {
	return Write(c.location, c.Projects)
}
