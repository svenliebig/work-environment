package context

import (
	"errors"
	"strings"

	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils"
)

var (
	ErrNoSuchProjectInDirectory = errors.New("there is no project in the directory")
	ErrProjectHasNoCI           = errors.New("project has no ci environment defined")
)

var (
	_ projectContex = &Context{}
)

type projectContex interface {
	baseContext
}

type Context struct {
	// the cwd path
	Cwd string

	baseContext *BaseContext

	ciId    string
	project *core.Project
}

func CreateContext() (*Context, error) {
	p, err := utils.GetPath([]string{})
	c := &Context{
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

func (c *Context) Validate() error {
	c.baseContext = &BaseContext{
		Cwd: c.Cwd,
	}

	err := c.baseContext.Validate()

	if err != nil {
		return err
	}

	if c.project != nil {
		return nil
	}

	config := c.baseContext.Configuration()

	for _, project := range config.Projects {
		if strings.Contains(c.Cwd, project.Path) {
			// TODO multiple projects found
			c.project = project
			return nil
		}
	}

	return ErrNoSuchProjectInDirectory
}

func (c *Context) Configuration() *core.Configuration {
	return c.baseContext.Configuration()
}

func (c *Context) ConfigurationPath() string {
	return c.baseContext.ConfigurationPath()
}

func (c *Context) Close() error {
	return c.baseContext.Close()
}

func (c *Context) Project() *core.Project {
	return c.project
}

// tell the context to use a ci with a specific id, overrides the default, to take
// the ci id from the project
func (c *Context) UseCI(id string) error {
	config := c.Configuration()

	if _, err := config.GetCIEnvironmentById(id); err != nil {
		return err
	} else {
		c.ciId = id
		return nil
	}
}

// returns the CI for the current project in the cwd.
func (c *Context) GetCI() (*core.CI, error) {
	if c.ciId != "" {
		return c.Configuration().GetCIEnvironmentById(c.ciId)
	}

	p := c.Project()

	if p.CI == nil {
		return nil, ErrProjectHasNoCI
	}

	return c.Configuration().GetCIEnvironmentById(p.CI.Id)
}
