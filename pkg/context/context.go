package context

import (
	"errors"
	"fmt"
	"strings"

	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/configwriter"
)

var (
	ErrNoSuchProjectInDirectory              = errors.New("there is no project in the directory")
	ErrProjectHasNoCI                        = errors.New("project has no ci environment defined")
	ErrProjectHasNoVCS                       = errors.New("project has no vcs environment defined")
	ErrProjectWithTheGivenIdentifierNotFound = errors.New("project with the given identifier not found")
)

var (
	_ ProjectContext = &projectContext{}
)

type ProjectContext interface {
	BaseContext

	Project() *core.Project
	GetCI() (*core.CI, error)
	UseCI(id string) error

	GetVCS() (*core.VCS, error)
	UpdateVCS(vcs *core.VCS, configuration string) error
}

// TODO same as in base.go
type projectContext struct {
	// the cwd path
	Cwd string

	baseContext *baseContext

	ciId    string
	project *core.Project
}

// create a ProjectContext out of the cwd and validates
// if a project and a configuration can be found.
func CreateProjectContext() (*projectContext, error) {
	p, err := utils.GetPath([]string{})
	c := &projectContext{
		Cwd: p,
	}

	if err != nil {
		return nil, err
	}

	err = c.validate("")

	if err != nil {
		return nil, err
	}

	return c, nil
}

// create a ProjectContext in this cwd and sets the given project name
// as project in this context, when the name is set, if not the default
// will take place, and the the fn will search for a project in the cwd
func CreateProjectContextWithProjectName(name string) (*projectContext, error) {
	p, err := utils.GetPath([]string{})
	c := &projectContext{
		Cwd: p,
	}

	if err != nil {
		return nil, err
	}

	err = c.validate(name)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *projectContext) validate(name string) error {
	bc, err := CreateBaseContext()

	if err != nil {
		return err
	}

	c.baseContext = bc.(*baseContext)
	config := c.baseContext.Configuration()

	if name != "" {
		contains, p := config.ContainsProject(name)
		if contains {
			c.project = p
			return nil
		} else {
			return fmt.Errorf("%w: identifier: %s", ErrProjectWithTheGivenIdentifierNotFound, name)
		}
	}

	if c.project != nil {
		return nil
	}

	var found *core.Project
	for _, project := range config.Projects {
		if strings.Contains(c.Cwd, project.Path) {
			if found == nil {
				found = project
				continue
			}

			if len(project.Path) > len(found.Path) {
				found = project
			}
		}
	}

	if found != nil {
		c.project = found
		return nil
	}

	return ErrNoSuchProjectInDirectory
}

func (c *projectContext) Configuration() *core.Configuration {
	return c.baseContext.Configuration()
}

func (c *projectContext) ConfigurationPath() string {
	return c.baseContext.ConfigurationPath()
}

func (c *projectContext) Info() *configwriter.ConfigWriter {
	cw := c.baseContext.Info()

	add, err := cw.AddSection("project")
	if err != nil {
		return nil
	}

	project := c.Project()
	branch, err := project.GetBranchName()

	if err != nil {
		return nil
	}

	add("name", project.Identifier)
	add("path", project.Path)
	add("remote_url", project.Git.RemoteUrl)
	add("branch", branch)

	defaultBranch, err := project.GetDefaultBranchName()

	if err != nil {
		return nil
	}

	add("default_branch", defaultBranch)

	if project.VCS != nil {
		add("vcs", project.VCS.Id+cli.Italic(fmt.Sprintf(" (configuration: '%s')", project.VCS.Configuration)))
	}

	if project.CI != nil {
		add("ci", project.CI.Id+cli.Italic(fmt.Sprintf(" (projectkey: '%s')", project.CI.ProjectKey)))
	}

	if project.CD != nil {
		add("cd", project.CD.Id+cli.Italic(fmt.Sprintf(" (projectid: '%d')", project.CD.ProjectId)))
	}

	return cw
}

func (c *projectContext) Close() error {
	return c.baseContext.Close()
}

func (c *projectContext) Project() *core.Project {
	return c.project
}

// tell the context to use a ci with a specific id, overrides the default, to take
// the ci id from the project
func (c *projectContext) UseCI(id string) error {
	config := c.Configuration()

	if _, err := config.GetCIEnvironmentById(id); err != nil {
		return err
	} else {
		c.ciId = id
		return nil
	}
}

// returns the CI for the current project in the cwd.
func (c *projectContext) GetCI() (*core.CI, error) {
	if c.ciId != "" {
		return c.Configuration().GetCIEnvironmentById(c.ciId)
	}

	p := c.Project()

	if p.CI == nil {
		return nil, ErrProjectHasNoCI
	}

	return c.Configuration().GetCIEnvironmentById(p.CI.Id)
}

func (c *projectContext) GetVCS() (*core.VCS, error) {
	p := c.Project()

	if p.VCS == nil {
		return nil, ErrProjectHasNoVCS
	}

	return c.Configuration().GetVCSEnvironmentById(p.VCS.Id)
}

func (c *projectContext) GetProjectsInPath() []*core.Project {
	return c.baseContext.GetProjectsInPath()
}

func (c *projectContext) UpdateVCS(vcs *core.VCS, configuration string) error {
	if vcs == nil {
		return c.Configuration().UpdateProjectVCS(c.Project().Identifier, nil)
	}

	return c.Configuration().UpdateProjectVCS(c.Project().Identifier, &core.ProjectVCS{
		Id:            vcs.Identifier,
		Configuration: configuration,
	})
}
