package context

import (
	"encoding/json"
	"os"

	"github.com/svenliebig/work-environment/pkg/core"
)

// this contains a base context for the work environment commands
// it secures to contain a work environment configuration set but
// not a project
type BaseContext struct {
	// the current working directory from where the command is called
	Cwd string

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

	return nil
}

func (c *BaseContext) Configuration() *core.Configuration {
	return c.configuration
}

func (c *BaseContext) ConfigurationPath() string {
	return c.configurationPath
}

func (c *BaseContext) Close() error {
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
