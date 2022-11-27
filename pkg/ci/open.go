package ci

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

func Open(ctx *context.Context) error {
	p := ctx.Project()
	c, err := ctx.GetCI()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get the ci", err)
	}

	if p.CI != nil {
		client, err := UseClient(ctx, c.CiType)

		if err != nil {
			return fmt.Errorf("error while trying to use client: %w", err)
		}

		plans, err := client.GetBranchPlans()

		if err != nil {
			return fmt.Errorf("%w: error while searching for branch plans", err)
		}

		if len(plans) == 0 {
			url := fmt.Sprintf("%s/browse/%s", c.Url, p.CI.ProjectKey)
			err = browser.Open(url)
		} else if len(plans) == 1 {
			url := fmt.Sprintf("%s/browse/%s", c.Url, plans[0].Key)
			err = browser.Open(url)
		} else {
			return fmt.Errorf("found more branch plans, but selecting one is not supported")
		}

		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("project has no ci defined")
	}

}
