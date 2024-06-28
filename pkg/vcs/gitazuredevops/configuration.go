package gitazuredevops

import (
	goctx "context"
	"encoding/json"
	"fmt"

	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/vcs"
)

type projectConfiguration struct {
	Project string
	WebURL  string
}

func (c *client) configuration() (*projectConfiguration, error) {
	str := c.ctx.Project().VCS.Configuration

	var configuration projectConfiguration

	if err := json.Unmarshal([]byte(str), &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (c *client) Configure(vcse *core.VCS) (string, error) {
	p := c.ctx.Project()
	fmt.Printf("%s '%s' for project '%s'...\n",
		cli.Colorize(cli.Blue, "Configure"),
		cli.Colorize(cli.Purple, vcse.Identifier),
		cli.Colorize(cli.Purple, p.Identifier),
	)

	env, err := environment(vcse)

	if err != nil {
		return "", err
	}

	c._env = env

	coreClient, err := c.coreClient()

	if err != nil {
		return "", fmt.Errorf("failed to create core client: %w", err)
	}

	gitClient, err := c.gitClient()

	if err != nil {
		return "", fmt.Errorf("failed to create git client: %w", err)
	}

	projects, err := coreClient.GetProjects(goctx.Background(), azcore.GetProjectsArgs{})

	if err != nil {
		return "", fmt.Errorf("failed to get projects: %w", err)
	}

	projectIdentifier, err := p.GetGitIdentifier()

	if err != nil {
		return "", fmt.Errorf("failed to get project identifier: %w", err)
	}

	for _, project := range projects.Value {
		repositories, err := gitClient.GetRepositories(goctx.Background(), git.GetRepositoriesArgs{
			Project: project.Name,
		})

		if err != nil {
			return "", err
		}

		for _, repository := range *repositories {
			if *repository.Name == projectIdentifier {

				result, err := json.Marshal(projectConfiguration{
					Project: *project.Name,
					WebURL:  *repository.WebUrl,
				})

				if err != nil {
					return "", err
				}

				return string(result), nil
			}
		}
	}

	return "", vcs.ErrRepositoryNotFound
}
