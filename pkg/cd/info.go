package cd

import (
	"fmt"

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

	// var buildResult string

	// if r.Success {
	// 	buildResult = cli.Colorize(cli.Green, "Success")
	// } else {
	// 	buildResult = cli.Colorize(cli.Red, "Failed")
	// }

	fmt.Printf("\nConfigured CD for '%s':\n\n", cli.Colorize(cli.Purple, p.Identifier))
	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  %s: \t%s", "CD Identifier", info.Identifier)
	fmt.Fprintf(w, "  %s: \t%s", "CD Type", info.Type)
	fmt.Fprintf(w, "  %s: \t%s", "CD URL", info.URL)
	fmt.Fprintf(w, "  %s: \t%s", "CD Version", info.Version)
	fmt.Fprintf(w, "  %s: \t%d", "Project Id", info.ProjectId)
	fmt.Fprintf(w, "")
	w.Print()

	// fmt.Printf("Latest Build (%s): %s\n", r.BuildNumber, buildResult)
	// if !r.Success {
	// 	fmt.Printf("Logs: %s\n", r.LogUrl)
	// 	fmt.Println()
	// 	for _, l := range r.Logs {
	// 		fmt.Printf("  > %s\n", l)
	// 	}
	// }
	// if r.IsBuilding {
	// 	fmt.Println()
	// 	fmt.Printf("A build is currently %s!\n", cli.Colorize(cli.Blue, "running"))
	// }

	// fmt.Println()

	return nil
}
