package ci

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Add(ctx context.ProjectContext, _ string, projectId string, key string, suggest bool) error {
	var err error
	config := ctx.Configuration()
	project := ctx.Project()

	options := make([]string, 0, len(config.CIEnvironments))
	for _, ci := range config.CIEnvironments {
		options = append(options, ci.Identifier)
	}

	switch len(options) {
	case 0:
		fmt.Printf("%s CI environments found in configuration. Please use '%s' to create an environment.\n", cli.Colorize(cli.Red, "No"), cli.Colorize(cli.Purple, "vcs ci create"))
		return nil
	case 1:
		err = ctx.UseCI(options[0])
	default:
		ciId := cli.Select("Select the CI environment to use", options)
		err = ctx.UseCI(ciId)
	}

	if err != nil {
		return fmt.Errorf("%w: error while trying to call use ci", err)
	}

	ci, err := ctx.GetCI()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get ci", err)
	}

	if suggest {
		client, err := UseClient(ctx, ci.Type)

		if err != nil {
			return err
		}

		key, err := client.GetPlanSuggestion()

		if err != nil {
			return err
		}

		err = config.UpdateProjectCI(project.Identifier, &core.ProjectCI{Id: ci.Identifier, ProjectKey: key})

		if err != nil {
			return err
		}

		// TODO set dirty and defer ctx.object.something() das dann schreibt
		err = ctx.Close()

		if err != nil {
			return err
		}

		fmt.Printf("\n%s %q ci with project key %q to project %q.\n",
			cli.Colorize(cli.Green, "Added"),
			ci.Identifier,
			key,
			project.Identifier,
		)
	} else {
		if key == "" {
			return errors.New("no key")
		}

		// TODO validate if key is a ci plan
	}

	return nil
}
