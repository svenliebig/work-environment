package vcs

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Remove(ctx context.ProjectContext) error {
	p := ctx.Project()

	if p.VCS == nil {
		fmt.Printf("The project '%s' has no vcs configured.\n", cli.Colorize(cli.Purple, p.Identifier))
		return nil
	}

	vcs := p.VCS
	err := ctx.UpdateVCS(nil, "")

	if err != nil {
		return fmt.Errorf("%w: error while trying to remove the vcs from the project", err)
	}

	err = ctx.Close()

	if err != nil {
		return fmt.Errorf("%w: error while trying to update the config", err)
	}

	fmt.Printf(
		"%s the vcs with the ID '%s' from the project '%s'.\n",
		cli.Colorize(cli.Red, "Removed"),
		cli.Colorize(cli.Purple, vcs.Id),
		cli.Colorize(cli.Purple, p.Identifier),
	)

	return nil
}
