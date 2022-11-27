package ci

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
)

// TODO suggest could be an extra cmd
// TODO configurations should be like a context object
func Add(ctx *context.Context, ciId string, projectId string, key string, suggest bool) error {
	config := ctx.Configuration()
	project := ctx.Project()

	if ciId != "" {
		err := ctx.UseCI(ciId)

		if err != nil {
			return fmt.Errorf("%w: error while trying to call use ci", err)
		}
	} else {
		fmt.Println("use", config.CIEnvironments[0].Identifier)
		err := ctx.UseCI(config.CIEnvironments[0].Identifier)

		if err != nil {
			return fmt.Errorf("%w: error while trying to call use ci", err)
		}
	}

	ci, err := ctx.GetCI()

	if err != nil {
		return fmt.Errorf("%w: error while trying to get ci", err)
	}

	if suggest {
		client, err := UseClient(ctx, ci.CiType)

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

		fmt.Printf("\nAdded %q ci with project key %q to project %q.\n", ci.Identifier, key, project.Identifier)
	} else {
		if key == "" {
			return errors.New("no key")
		}

		// TODO validate if key is a ci plan
	}

	return nil
}
