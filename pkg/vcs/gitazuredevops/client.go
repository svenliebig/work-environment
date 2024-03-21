package gitazuredevops

import (
	goctx "context"
	"encoding/json"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/vcs"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

type client struct {
	ctx         context.ProjectContext
	connection  *azuredevops.Connection
	_coreClient azcore.Client
	_gitClient  git.Client
}

type configuration struct {
	Project string
	WebURL  string
}

func (c *client) conn() (*azuredevops.Connection, error) {
	if c.connection != nil {
		return c.connection, nil
	}

	vcs, err := c.ctx.GetVCS()

	if err != nil {
		return nil, err
	}

	c.connection = azuredevops.NewPatConnection(vcs.Url, vcs.AccessToken)

	return c.connection, nil
}

func (c *client) coreClient() (azcore.Client, error) {
	if c._coreClient != nil {
		return c._coreClient, nil
	}

	connection, err := c.conn()

	if err != nil {
		return nil, err
	}

	client, err := azcore.NewClient(goctx.Background(), connection)

	if err != nil {
		return nil, err
	}

	c._coreClient = client

	return c._coreClient, nil
}

func (c *client) gitClient() (git.Client, error) {
	if c._gitClient != nil {
		return c._gitClient, nil
	}

	connection, err := c.conn()

	if err != nil {
		return nil, err
	}

	client, err := git.NewClient(goctx.Background(), connection)

	if err != nil {
		return nil, err
	}

	c._gitClient = client

	return c._gitClient, nil
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

func (c *client) List() ([]string, error) {
	fmt.Println("List")

	client, err := c.gitClient()

	if err != nil {
		return nil, err
	}

	var bla = "tp7"

	repositories, err := client.GetRepositories(goctx.Background(), git.GetRepositoriesArgs{
		Project: &bla,
	})

	if err != nil {
		return nil, err
	}

	for _, repository := range *repositories {
		fmt.Println(*repository.Name)
	}

	return nil, nil
}

func init() {
	vcs.RegisterClient("gitazuredevops", func(ctx context.ProjectContext) vcs.Client {
		return &client{
			ctx: ctx,
		}
	})
}
