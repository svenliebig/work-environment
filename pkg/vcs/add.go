package vcs

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Add(ctx context.ProjectContext) error {
	vcse, err := ctx.Configuration().GetVCSEnvironmentById("azure-tp6")

	if err != nil {
		return err
	}

	err = ctx.UpdateVCS(vcse, "")

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcse.Type)

	if err != nil {
		return err
	}

	configuration, err := client.Configure()

	if err != nil {
		return err
	}

	err = ctx.UpdateVCS(vcse, configuration)

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
