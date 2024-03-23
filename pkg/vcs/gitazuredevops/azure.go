package gitazuredevops

import (
	goctx "context"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

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
