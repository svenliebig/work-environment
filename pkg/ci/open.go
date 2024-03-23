package ci

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/browser"
)

func Open(ctx context.ProjectContext) error {
	p := ctx.Project()
	c, err := ctx.GetCI()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get the ci", err)
	}

	if p.CI != nil {
		client, err := UseClient(ctx, c.Type)

		if err != nil {
			return fmt.Errorf("error while trying to use client: %w", err)
		}

		url, err := client.GetBranchPlanUrl()

		if err != nil {
			return fmt.Errorf("%w: error while searching for branch plans", err)
		}

		err = browser.Open(url)

		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("project has no ci defined")
	}

}
