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
	ctx          *context.Context
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

	return browser.Open(fmt.Sprintf("%s/deploy/viewDeploymentProjectEnvironments.action?id=%d", c.bambooClient.BaseUrl, c.ctx.Project().CD.ProjectId))
}

func (c *client) Info() (*cd.ClientInfo, error) {
	_, err := c.bamboo()

	if err != nil {
		return nil, err
	}

	return &cd.ClientInfo{
		Identifier: c.ctx.Project().CD.Id,
		Type:       c.ci.CiType,
		URL:        c.bambooClient.BaseUrl,
		Version:    c.ci.Version,
		ProjectId:  c.ctx.Project().CD.ProjectId,
	}, nil
}

func init() {
	cd.RegisterClient("bamboo", func(ctx *context.Context) cd.Client {
		return &client{
			ctx: ctx,
		}
	})
}
