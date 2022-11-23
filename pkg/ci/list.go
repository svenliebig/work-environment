package ci

import (
	"fmt"

	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

func List(p string) error {
	c, err := GetConfig(p)

	if err != nil {
		return fmt.Errorf("%w: error while trying to read the ci configuration", err)
	}

	w := &tablewriter.TableWriter{}
	fmt.Printf("\nAvailable CI Environments:\n\n")
	fmt.Fprintf(w, "| %sID%s \t| Type \t| URL \t|", cli.Blue, cli.Reset)
	for _, e := range c.Environments {
		fmt.Fprintf(w, "| %s \t| %s \t| %s \t|", e.Identifier, e.CiType, e.Url)
	}
	w.Print()
	fmt.Println()

	return nil
}
