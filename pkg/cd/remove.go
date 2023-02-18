package cd

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Remove(ctx context.ProjectContext) error {
	config := ctx.Configuration()
	p := ctx.Project()

	if p.CD == nil {
		fmt.Printf("The project '%s' has no CD configured.\n", cli.Colorize(cli.Purple, p.Identifier))
		return nil
	}

	cd := p.CD
	err := config.UpdateProjectCD(p.Identifier, nil)

	if err != nil {
		return fmt.Errorf("%w: error while trying to update the project cd", err)
	}

	err = ctx.Close()

	if err != nil {
		return fmt.Errorf("%w: error while trying to update the config", err)
	}

	fmt.Printf("Removed the CD with the ID '%s' from the project '%s'.\n", cli.Colorize(cli.Purple, fmt.Sprintf("%d", cd.ProjectId)), cli.Colorize(cli.Purple, p.Identifier))
	return nil
}
