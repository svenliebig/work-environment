package context

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils"
)

var (
	_ BaseContext = &baseContext{}
)

type BaseContext interface {
	Configuration() *core.Configuration
	ConfigurationPath() string
	Close() error
	GetProjectsInPath() []*core.Project

	Validate() error
}

// returns a base context for the work environment validating the most basic
// things that are needed to run a command, if no error is returned, it's
// ensured that you have a configuration available that you can work with.
//
// @Alex i would like to do this
// the cwd is an optional parameter to overwrite the default cwd, which is the
// current working directory from where the command is executed
func CreateBaseContext() (BaseContext, error) {
	p, err := utils.GetPath([]string{})

	c := &baseContext{
		Cwd: p,
	}

	if err != nil {
		return nil, err
	}

	err = c.Validate()

	if err != nil {
		return nil, err
	}

	return c, nil
}

type baseContext struct {
	// the current working directory from where the command is called
	Cwd string

	validated         bool
	configurationPath string
	configuration     *core.Configuration
}

// ensures that the context is setup correct, this needs to be used
// after the initialization of the context
func (c *baseContext) Validate() error {
	if c.configurationPath != "" && c.configuration != nil {
		return nil
	}

	p, err := getConfigurationPath(c.Cwd)

	if err != nil {
		return err
	}

	config, err := readConfig(p)

	if err != nil {
		return err
	}

	c.configurationPath = p
	c.configuration = config
	c.validated = true

	return nil
}

func (c *baseContext) Configuration() *core.Configuration {
	if !c.validated {
		notValidated()
	}

	return c.configuration
}

func (c *baseContext) ConfigurationPath() string {
	if !c.validated {
		notValidated()
	}

	return c.configurationPath
}

// writes the configuration if it's dirty
func (c *baseContext) Close() error {
	if !c.validated {
		notValidated()
	}

	if c.configuration.IsDirty() {
		// TODO close on configuration...
		// close() get's some parameter about path or something
		return c.updateConfig()
	}

	return nil
}

func (c *baseContext) updateConfig() error {
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

func (c *baseContext) GetProjectsInPath() []*core.Project {
	pip := make([]*core.Project, 0, len(c.Configuration().Projects))

	for _, p := range c.Configuration().Projects {
		if strings.Contains(p.Path, c.Cwd) {
			pip = append(pip, p)
		}
	}

	return pip
}

func notValidated() {
	log.Fatal("context used before calling Validate()")
}
