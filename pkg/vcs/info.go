package vcs

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
)

func Info(ctx context.ProjectContext) error {
	vcse, err := ctx.GetVCS()

	if err != nil {
		return err
	}

	client, err := UseClient(ctx, vcse)

	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cli.Bold(cli.Underline("Version Control System")))
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("identifier"), vcse.Identifier)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("type"), vcse.Type)
	fmt.Fprintf(w, "  %s\t%s\n", cli.Bold("configuration"), vcse.Configuration)
	w.Flush()
	fmt.Println()

	return client.Info()
}
