package context

import (
	"errors"
	"fmt"
	"strings"

	"github.com/svenliebig/work-environment/pkg/core"
)

var (
	ErrNoSuchProject              = errors.New("no such project in configuration")
	ErrConfigurationDoesNotExists = errors.New("configuration does not exists")
	ErrProjectHasNoCI             = errors.New("project has no ci environment defined")
	ErrCIAlreadyExists            = errors.New("a ci with that identifier does already exists")
	ErrNoSuchCI                   = errors.New("there is no ci with that identifier available")
)

// TODO IDEA:
// Project Context, that init's and throws when necessary things are not provided
// Work Environment Context, that init's and just wants the work environment
// Project Context is inherits the WE Context
type Context struct {
	// the cwd path
	Path string

	ciId              string
	configurationPath string
	configuration     *core.Configuration
	project           *core.Project
}

func (c *Context) Defer() {
	// if configuration is diry, then do things
}

// tell the context to use a ci with a specific id, overrides the default, to take
// the ci id from the project
func (c *Context) UseCI(id string) error {
	config, err := c.GetConfiguration()

	if err != nil {
		return err
	}

	if _, err := config.GetCIEnvironmentById(id); err != nil {
		return err
	} else {
		c.ciId = id
		return nil
	}
}

// uses the path of Context to get the first project that is
// in the current path.
func (c *Context) GetProject() (*core.Project, error) {
	if c.project != nil {
		return c.project, nil
	}

	config, err := c.GetConfiguration()

	if err != nil {
		return nil, err
	}

	for _, project := range config.Projects {
		if strings.Contains(c.Path, project.Path) {
			// TODO multiple projects found
			c.project = project
			return c.project, nil
		}
	}

	return nil, ErrNoSuchProject
}

// returns the CI for the current project in the cwd.
func (c *Context) GetCI() (*core.CI, error) {
	fmt.Println(c.ciId)
	if c.ciId != "" {
		return c.configuration.GetCIEnvironmentById(c.ciId)
	}

	p, err := c.GetProject()

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if p.CI == nil {
		return nil, ErrProjectHasNoCI
	}

	return c.configuration.GetCIEnvironmentById(p.CI.Id)
}
