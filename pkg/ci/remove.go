package ci

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Remove(ctx *context.Context) error {
	config, err := ctx.GetConfiguration()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get the config", err)
	}

	p, err := ctx.GetProject()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get the project", err)
	}

	ci, err := ctx.GetCI()

	if err != nil {
		if errors.Is(err, context.ErrProjectHasNoCI) {
			fmt.Printf("The project '%s' has no ci configured.\n", cli.Colorize(cli.Purple, p.Identifier))
			return nil
		}

		return err
	}

	config.UpdateProjectCI(p.Identifier, nil)

	if err != nil {
		return fmt.Errorf("%w: error while trying to update the project ci", err)
	}

	err = ctx.UpdateConfig()

	if err != nil {
		return fmt.Errorf("%w: error while trying to update the config", err)
	}

	fmt.Printf("Removed '%s' from the project '%s'.\n", cli.Colorize(cli.Purple, ci.Identifier), cli.Colorize(cli.Purple, p.Identifier))
	return nil
}
