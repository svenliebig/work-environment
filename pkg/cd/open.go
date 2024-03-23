package cd

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
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

		err = client.Open()

		if err != nil {
			return fmt.Errorf("%w: error while trying to open the CD", err)
		}
	} else {
		return fmt.Errorf("project has no ci defined")
	}

	return nil
}
