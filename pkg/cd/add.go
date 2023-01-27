package cd

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/ci"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Add(ctx *context.Context) error {
	config := ctx.Configuration()
	project := ctx.Project()
	ci_, err := ctx.GetCI()

	if err != nil {
		if errors.Is(err, context.ErrProjectHasNoCI) {
			return fmt.Errorf("project has no ci configured, currently only bamboo is supported as CD, use:\n\n\twe ci add --suggest\n\nto add a ci configuration")
		}

		return err
	}

	client, err := ci.UseClient(ctx, ci_.CiType)

	if err != nil {
		return fmt.Errorf("%w: error while trying to get client", err)
	}

	id, err := client.GetCD()

	if err != nil {
		return err
	}

	err = config.UpdateProjectCD(project.Identifier, &core.ProjectCD{Id: ci_.Identifier, ProjectId: id})

	if err != nil {
		return err
	}

	// TODO set dirty and defer ctx.object.something() das dann schreibt
	err = ctx.Close()

	if err != nil {
		return err
	}

	fmt.Printf("\nAdded '%s' CD with project id '%s' to project '%s'.\n",
		cli.Colorize(cli.Purple, ci_.Identifier),
		cli.Colorize(cli.Purple, fmt.Sprintf("%d", id)),
		cli.Colorize(cli.Purple, project.Identifier),
	)

	return nil
}
