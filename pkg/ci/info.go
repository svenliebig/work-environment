package ci

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

type InfoOptions struct {
	// prints only the url of the plan if true
	Url bool
}

func Info(ctx context.ProjectContext, options *InfoOptions) error {
	p := ctx.Project()
	ci, err := ctx.GetCI()

	if err != nil {
		if errors.Is(err, context.ErrProjectHasNoCI) {
			return fmt.Errorf("project has no ci configured, use:\n\n\twe ci add --suggest\n\nto add a ci configuration")
		}

		return err
	}

	client, err := UseClient(ctx, "bamboo")

	if err != nil {
		return err
	}

	if options.Url {
		url, err := client.GetBranchPlanUrl()

		if err != nil {
			return err
		}

		fmt.Println(url)
		return nil
	}

	r, err := client.LatestBuildResult()

	if err != nil {
		return err
	}

	var buildResult string

	if r.Success {
		buildResult = cli.Colorize(cli.Green, "Success")
	} else {
		buildResult = cli.Colorize(cli.Red, "Failed")
	}

	fmt.Printf("\nConfigured CI for '%s':\n\n", cli.Colorize(cli.Purple, p.Identifier))
	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  %s: \t%s", "CI Identifier", ci.Identifier)
	fmt.Fprintf(w, "  %s: \t%s", "CI Type", ci.CiType)
	fmt.Fprintf(w, "  %s: \t%s", "CI URL", ci.Url)
	fmt.Fprintf(w, "  %s: \t%s", "CI Version", ci.Version)
	fmt.Fprintf(w, "  %s: \t%s", "Project Key", p.CI.ProjectKey)
	fmt.Fprintf(w, "")
	w.Print()

	fmt.Printf("Latest Build (%s): %s\n", r.BuildNumber, buildResult)
	if !r.Success {
		fmt.Printf("Logs: %s\n", r.LogUrl)
		fmt.Println()
		for _, l := range r.Logs {
			fmt.Printf("  > %s\n", l)
		}
	}
	if r.IsBuilding {
		fmt.Println()
		fmt.Printf("A build is currently %s!\n", cli.Colorize(cli.Blue, "running"))
	}

	fmt.Println()

	return nil
}
