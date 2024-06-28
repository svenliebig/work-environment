package vcs

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Add(ctx context.ProjectContext) error {
	c := ctx.Configuration()

	options := make([]string, 0, len(c.VCSEnvironments))
	for _, v := range c.VCSEnvironments {
		options = append(options, v.Identifier)
	}
	env := cli.Select("Select the VCS environment to add", options)

	vcse, err := c.GetVCSEnvironmentById(env)

	if err != nil {
		return err
	}

	err = ConfigureClient(ctx, vcse)

	if err != nil {
		return err
	}

	err = ctx.Close()

	if err != nil {
		return err
	}

	fmt.Printf(
		"%s the vsce '%s' to the project '%s'.\n",
		cli.Colorize(cli.Green, "Added"),
		cli.Colorize(cli.Purple, vcse.Identifier),
		cli.Colorize(cli.Purple, ctx.Project().Identifier),
	)

	return nil
}
