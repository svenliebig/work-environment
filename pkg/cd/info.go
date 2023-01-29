package cd

import (
	"fmt"
	"time"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

func Info(ctx *context.Context) error {
	p := ctx.Project()

	client, err := UseClient(ctx, "bamboo")

	if err != nil {
		return err
	}

	info, err := client.Info()

	if err != nil {
		return err
	}

	envs, err := client.Environments()

	if err != nil {
		return err
	}

	// var buildResult string

	fmt.Printf("\nConfigured CD for '%s':\n\n", cli.Colorize(cli.Purple, p.Identifier))
	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  %s: \t%s", "CD Identifier", info.Identifier)
	fmt.Fprintf(w, "  %s: \t%s", "CD Type", info.Type)
	fmt.Fprintf(w, "  %s: \t%s", "CD URL", info.URL)
	fmt.Fprintf(w, "  %s: \t%s", "CD Version", info.Version)
	fmt.Fprintf(w, "  %s: \t%d", "Project Id", info.ProjectId)
	fmt.Fprintf(w, "")
	w.Print()
	w.Flush()
	fmt.Fprintf(w, "Environments:")
	fmt.Fprintf(w, "  Name \t Id \t Last Build\t\t Release Name")

	// TODO this can be a go routine
	for _, v := range envs {
		res, err := client.DeployResult(v.Id)

		if err != nil {
			return err
		}

		var buildResult string
		// this belongs into the bamboo client
		if res.DeploymentState == "SUCCESS" {
			buildResult = cli.Colorize(cli.Green, "Success")
		} else {
			buildResult = cli.Colorize(cli.Red, res.DeploymentState)
		}

		t := time.Unix(int64(res.Finished/1000), 0)
		fmt.Fprintf(w, "  %s: \t %d \t %s \t%s \t %s", v.Name, v.Id, buildResult, t.Format(time.RFC822), res.Version)
	}
	fmt.Fprintf(w, "")
	w.Print()

	return nil
}
