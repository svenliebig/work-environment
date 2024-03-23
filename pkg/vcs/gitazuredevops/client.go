package gitazuredevops

import (
	goctx "context"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/vcs"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

type client struct {
	ctx         context.ProjectContext
	_connection *azuredevops.Connection
	_env        *envConfig
	_coreClient azcore.Client
	_gitClient  git.Client
}

const identifier = "gitazuredevops"

func (c *client) PullRequestWebURL() (string, error) {
	config, err := c.configuration()

	if err != nil {
		return "", err
	}

	project := c.ctx.Project()
	branch, err := project.GetBranchName()

	if err != nil {
		return "", err
	}

	defaultBranch, err := project.GetDefaultBranchName()

	fmt.Println(defaultBranch)

	return fmt.Sprintf("%s/pullrequestcreate?sourceRef=%s&targetRef=%s", config.WebURL, branch, defaultBranch), nil
}

func (c *client) WebURL() (string, error) {
	config, err := c.configuration()

	if err != nil {
		return "", err
	}

	return config.WebURL, nil
}

func (c *client) Info() error {
	fmt.Println("List")

	client, err := c.gitClient()

	if err != nil {
		return err
	}

	var bla = "tp7"

	repositories, err := client.GetRepositories(goctx.Background(), git.GetRepositoriesArgs{
		Project: &bla,
	})

	if err != nil {
		return err
	}

	for _, repository := range *repositories {
		fmt.Println(*repository.Name)
	}

	return nil
}

func init() {
	vcs.RegisterClient(identifier, func(ctx context.ProjectContext) vcs.Client {
		return &client{
			ctx: ctx,
		}
	})
}
