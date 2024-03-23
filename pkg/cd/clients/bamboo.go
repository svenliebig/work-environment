package clients

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/api/bamboo"
	"github.com/svenliebig/work-environment/pkg/cd"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

type client struct {
	ctx          context.ProjectContext
	bambooClient *bamboo.Client
	ci           *core.CI
}

func (c *client) bamboo() (*bamboo.Client, error) {
	if c.bambooClient == nil {
		ci, err := c.ctx.GetCI()

		if err != nil {
			return nil, err
		}

		c.bambooClient = &bamboo.Client{
			BaseUrl:   ci.Url,
			AuthToken: ci.AuthToken,
		}

		c.ci = ci
	}

	if c.ctx.Project().CD == nil {
		return nil, fmt.Errorf("no CD configured for project %q", c.ctx.Project().Identifier)
	}

	return c.bambooClient, nil
}

func (c *client) Open() error {
	_, err := c.bamboo()

	if err != nil {
		return err
	}

	url, err := c.GetPlanUrl()

	if err != nil {
		return err
	}

	return browser.Open(url)
}

func (c *client) Info() (*cd.ClientInfo, error) {
	_, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	return &cd.ClientInfo{
		Identifier: c.ctx.Project().CD.Id,
		Type:       c.ci.Type,
		URL:        c.bambooClient.BaseUrl,
		Version:    c.ci.Version,
		ProjectId:  c.ctx.Project().CD.ProjectId,
	}, nil
}

func (c *client) DeployResult(environmentId int) (*cd.DeployResult, error) {
	_, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	res, err := c.bambooClient.DeploymentResultsByEnvironmentId(environmentId, 1)

	if err != nil {
		return nil, err
	}

	if len(res.Results) != 1 {
		return nil, fmt.Errorf("there are no results for the environment with the id %d", environmentId)
	} else {
		return &cd.DeployResult{
			Id:              res.Results[0].Id,
			Version:         res.Results[0].DeploymentVersionName,
			DeploymentState: res.Results[0].DeploymentState,
			Finished:        res.Results[0].FinishedDate,
		}, nil
	}
}

func (c *client) GetPlanUrl() (string, error) {
	_, err := c.bamboo()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/deploy/viewDeploymentProjectEnvironments.action?id=%d", c.bambooClient.BaseUrl, c.ctx.Project().CD.ProjectId), nil
}

func (c *client) Environments() ([]*cd.Environment, error) {
	_, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	res, err := c.bambooClient.DeployProjectById(c.ctx.Project().CD.ProjectId)

	if err != nil {
		return nil, err
	}

	envs := res.Environments
	ret := make([]*cd.Environment, len(envs))

	for i, v := range envs {
		ret[i] = &cd.Environment{
			Name: v.Name,
			Id:   v.Id,
		}
	}

	return ret, nil
}

func init() {
	cd.RegisterClient("bamboo", func(ctx context.ProjectContext) cd.Client {
		return &client{
			ctx: ctx,
		}
	})
}
