package vcs

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

type CreateParameter = core.VCS

func Create(ctx context.BaseContext, p CreateParameter) error {
	config := ctx.Configuration()

	vcsEnvironment, err := SetupClient(ctx, core.VCS{
		Type:       p.Type,
		Identifier: p.Identifier,
	})

	if err != nil {
		return err
	}

	err = config.AddVCS(&vcsEnvironment)

	if errors.Is(err, core.ErrVCSAlreadyExists) {
		q := fmt.Sprintf("\nThe Identifier '%s' is already declared in your configuration.\nDo you want to overwrite it? [y/n] ", cli.Colorize(cli.Purple, p.Identifier))
		answer := cli.Question(q, []string{"y", "n"})

		if answer == "n" {
			return nil
		}

		// do override
	} else if err != nil {
		return err
	}

	err = ctx.Close()

	if err != nil {
		return err
	}

	fmt.Printf("\n%s added a new VCS to your work environment:\n\n", cli.Colorize(cli.Green, "Successfully"))

	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  Identifier: \t%s", vcsEnvironment.Identifier)
	fmt.Fprintf(w, "  Type: \t%s", vcsEnvironment.Type)
	fmt.Fprintf(w, "  Type: \t%s", vcsEnvironment.Configuration)
	w.Print()
	fmt.Println()

	return nil
}
