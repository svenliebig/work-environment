package ci

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

func Info(ctx *context.Context) error {
	p, err := ctx.GetProject()

	if err != nil {
		return err
	}

	ci, err := ctx.GetCI()

	if err != nil {
		if errors.Is(err, context.ErrProjectHasNoCI) {
			return fmt.Errorf("project has no ci configured, use:\n\n\twe ci add --suggest\n\nto add a ci configuration")
		}

		return err
	}

	fmt.Printf("\nConfigured CI for %q:\n\n", p.Identifier)
	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  %s: \t%s", "CI Identifier", ci.Identifier)
	fmt.Fprintf(w, "  %s: \t%s", "CI URL", ci.Url)
	fmt.Fprintf(w, "  %s: \t%s", "Project Key", p.CI.ProjectKey)
	w.Print()
	fmt.Println()

	return nil
}
