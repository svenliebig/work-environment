package we

import (
	"fmt"
	"strings"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/git"
	"github.com/svenliebig/work-environment/pkg/utils/progress"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

type ListOptions struct {
	All    bool
	Tags   string
	Filter string
}

func List(ctx context.BaseContext, opts *ListOptions) error {
	projects := ctx.GetProjectsInPath()

	if opts.Filter != "" {
		p := projects
		projects = make([]*core.Project, 0, len(p))
		for _, a := range p {
			if strings.Contains(a.Identifier, opts.Filter) {
				projects = append(projects, a)
			}
		}
	}

	if len(projects) == 0 {
		msg := ""

		if opts.Filter != "" {
			msg = fmt.Sprintf("No projects found for filter '%s'.\n\n", opts.Filter)
		} else {
			msg = "No projects found.\n\n"
		}
		fmt.Print(msg)

		return nil
	}

	w := &tablewriter.TableWriter{}

	fmt.Fprintf(w, "")
	fmt.Fprintf(w, "project \tgit branch")
	fmt.Fprintf(w, "")

	max := 6
	wg := sync.WaitGroup{}
	c := make(chan *core.Project, max)

	go (func() {
		for _, project := range projects {
			c <- project
		}
		close(c)
	})()

	p := &progress.Progress{
		Max: len(projects),
	}
	pw := &progress.Writer{}

	fmt.Println()
	fmt.Println("Finding projects and evaluating git status...")
	fmt.Println()

	for i := 0; i < max; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for project := range c {
				b, _ := git.BranchName(project.Path)
				s, _ := git.Status(project.Path)

				var gs string
				if s.Dirty() {
					gs = cli.Colorize(cli.Yellow, fmt.Sprintf("%s [%s]", b, s.String()))
				} else {
					gs = cli.Colorize(cli.Green, b)
				}

				fmt.Fprintf(w, "%s \t(%s)", project.Identifier, gs)

				p.Add(1)
				pw.Print(p)
			}
		}()
	}

	wg.Wait()

	fmt.Println()
	w.Print()
	w.Flush()
	fmt.Println()

	return nil
}
