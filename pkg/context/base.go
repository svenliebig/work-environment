package context

import (
	"encoding/json"
	"log"
	"os"

	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils"
)

var (
	_ baseContext = &BaseContext{}
)

type baseContext interface {
	Configuration() *core.Configuration
	ConfigurationPath() string
	Close() error
	Validate() error
}

func CreateBaseContext() (*BaseContext, error) {
	p, err := utils.GetPath([]string{})
	c := &BaseContext{
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

// this contains a base context for the work environment commands
// it secures to contain a work environment configuration set but
// not a project.
//
// the current working directory is required.
type BaseContext struct {
	// the current working directory from where the command is called
	Cwd string

	validated         bool
	configurationPath string
	configuration     *core.Configuration
}

// ensures that the context is setuo correct, this needs to be used
// after the initialization of the context
func (c *BaseContext) Validate() error {
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

func (c *BaseContext) Configuration() *core.Configuration {
	if !c.validated {
		notValidated()
	}

	return c.configuration
}

func (c *BaseContext) ConfigurationPath() string {
	if !c.validated {
		notValidated()
	}

	return c.configurationPath
}

// writes the configuration if it's dirty
func (c *BaseContext) Close() error {
	if !c.validated {
		notValidated()
	}

	if c.configuration.IsDirty() {
		return c.updateConfig()
	}

	return nil
}

func (c *BaseContext) updateConfig() error {
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

func notValidated() {
	log.Fatal("context used before calling Validate()")
}
