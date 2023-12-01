package cd

import (
	"fmt"
	"sync"
	"time"

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

	client, err := UseClient(ctx, "bamboo")

	if err != nil {
		return err
	}

	if options.Url {
		url, err := client.GetPlanUrl()

		if err != nil {
			return err
		}

		fmt.Println(url)
		return nil
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

	var wg sync.WaitGroup
	max := 3
	c := make(chan *Environment, max)

	go func() {
		for _, e := range envs {
			c <- e
		}
		close(c)
	}()

	for i := 0; i < max; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for v := range c {

				res, err := client.DeployResult(v.Id)

				if err != nil {
					fmt.Println("ERROR: cry", err)
				}

				var buildResult string
				// this belongs into the bamboo client
				if res.DeploymentState == "SUCCESS" {
					buildResult = cli.Colorize(cli.Green, "Success")
				} else {
					buildResult = cli.Colorize(cli.Red, res.DeploymentState)
				}

				// of course... this is not sorted now...
				t := time.Unix(int64(res.Finished/1000), 0)
				fmt.Fprintf(w, "  %s: \t %d \t %s \t%s \t %s", v.Name, v.Id, buildResult, t.Format(time.RFC822), res.Version)
			}
		}()
	}

	wg.Wait()
	fmt.Fprintf(w, "")
	w.Print()

	return nil
}
