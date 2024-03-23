package gitazuredevops

import (
	goctx "context"
	"encoding/json"
	"fmt"

	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/svenliebig/work-environment/pkg/vcs"
)

type configuration struct {
	Project string
	WebURL  string
}

func (c *client) configuration() (*configuration, error) {
	str := c.ctx.Project().VCS.Configuration

	var configuration configuration

	if err := json.Unmarshal([]byte(str), &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (c *client) Configure() (string, error) {
	p := c.ctx.Project()
	fmt.Printf("Configure '%s' for project '%s'...\n", p.VCS.Id, p.Identifier)

	coreClient, err := c.coreClient()
	gitClient, err := c.gitClient()

	if err != nil {
		return "", err
	}

	projects, err := coreClient.GetProjects(goctx.Background(), azcore.GetProjectsArgs{})

	if err != nil {
		return "", err
	}

	for _, project := range projects.Value {
		repositories, err := gitClient.GetRepositories(goctx.Background(), git.GetRepositoriesArgs{
			Project: project.Name,
		})

		if err != nil {
			return "", err
		}

		for _, repository := range *repositories {
			if *repository.Name == p.Identifier {
				fmt.Println("Repository found.")

				result, err := json.Marshal(configuration{
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
