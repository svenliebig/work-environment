package gitazuredevops

import (
	goctx "context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
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
	client, err := c.gitClient()

	if err != nil {
		return err
	}

	config, err := c.configuration()

	if err != nil {
		return err
	}

	repoId, err := c.getRepositoryId()

	if err != nil {
		return err
	}

	repository, err := client.GetRepository(goctx.Background(), git.GetRepositoryArgs{
		Project:      &config.Project,
		RepositoryId: &repoId,
	})

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("project"), *repository.Project.Name)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("ssh_url"), *repository.SshUrl)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("web_url"), *repository.WebUrl)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("default_branch"), *repository.DefaultBranch)
	fmt.Fprintf(w, "  %s\t%d\n", cli.Bold("size"), *repository.Size)

	w.Flush()

	fmt.Println()

	return nil
}

func init() {
	vcs.RegisterClient(identifier, func(ctx context.ProjectContext) vcs.Client {
		return &client{
			ctx: ctx,
		}
	})
}
