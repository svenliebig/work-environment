package we

import (
	"fmt"
	"path/filepath"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

func Update(ctx *context.BaseContext) error {
	p := ctx.ConfigurationPath()
	d := filepath.Dir(p)

	fmt.Print("\nScanning for projects...\n\n")
	projects, err := scanForProjects(d)

	if err != nil {
		return err
	}

	c := ctx.Configuration()

	// check with bench tests, what does take more memory
	add := make([]*core.Project, 0, len(projects))
	delete := make([]*core.Project, 0, len(projects))

	for _, project := range projects {
		if ok, _ := c.ContainsProject(project.Identifier); ok {
			// do nothing i guess
		} else {
			add = append(add, project)
		}
	}

	for _, configProject := range c.Projects {
		contains := false
		for _, project := range projects {
			if project.Identifier == configProject.Identifier {
				contains = true
			}
		}

		if !contains {
			delete = append(delete, configProject)
		}
	}

	w := &tablewriter.TableWriter{}

	if len(add) > 0 {
		fmt.Print("Added projects:\n\n")

		for _, a := range add {
			c.Projects = append(c.Projects, a)
			fmt.Fprintf(w, "  %s \t-> %s", cli.Colorize(cli.Green, a.Identifier), a.Path)
		}
		w.Print()
		w.Flush()
	} else {
		fmt.Print("No new projects.\n\n")
	}

	if len(delete) > 0 {
		fmt.Print("\nRemoved projects:\n\n")

		for _, d := range delete {
			c.RemoveProject(d)
			fmt.Fprintf(w, "  %s \t-> %s", cli.Colorize(cli.Red, d.Identifier), d.Path)
		}
		w.Print()
		w.Flush()
	} else {
		fmt.Print("No removed projects.\n")
	}

	if len(c.Projects) > 0 {
		fmt.Print("\nWork Environment:\n\n")

		for _, p := range c.Projects {
			fmt.Fprintf(w, "  %s \t-> %s", cli.Colorize(cli.Purple, p.Identifier), p.Path)
		}
		w.Print()
		w.Flush()
	} else {
		fmt.Print("\nWork Environment: <empty>\n\n")
	}

	if len(add) == 0 && len(delete) == 0 {
		fmt.Print("\nYour work environment is up to date!\n\n")
	} else {

		fmt.Printf("\n%s updated your configuration in %q.\n", cli.Colorize(cli.Green, "Successfully"), ctx.ConfigurationPath())
	}

	return ctx.Close()
}
