package gitazuredevops

import (
	goctx "context"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	azcore "github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

func (c *client) connection() (*azuredevops.Connection, error) {
	if c._connection != nil {
		return c._connection, nil
	}

	env, err := c.environment()

	if err != nil {
		return nil, err
	}

	c._connection = azuredevops.NewPatConnection(fmt.Sprintf("https://dev.azure.com/%s", env.Organization), env.AccessToken)

	return c._connection, nil
}

func (c *client) coreClient() (azcore.Client, error) {
	if c._coreClient != nil {
		return c._coreClient, nil
	}

	connection, err := c.connection()

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

	connection, err := c.connection()

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
