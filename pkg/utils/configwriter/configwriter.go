package configwriter

import (
	"errors"
	"fmt"

	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

var (
	ErrSectionAlreadyExists = errors.New("section already exists")
)

type ConfigWriter struct {
	Title          string
	sections       []*Section
	currentSection *Section
}

type sectionWriter struct {
}

type Section struct {
	Name   string
	Values []Value
}

type Value struct {
	Key   string
	Value string
}

func (c *ConfigWriter) AddSection(name string) (func(key, value string), error) {
	for _, s := range c.sections {
		if s.Name == name {
			return nil, ErrSectionAlreadyExists
		}
	}

	section := Section{Name: name}
	c.sections = append(c.sections, &section)

	return func(key, value string) {
		section.Values = append(section.Values, Value{Key: key, Value: value})
	}, nil
}

func title(title string) string {
	return cli.Bold(cli.Underline(cli.Colorize(cli.White, title)))
}

func section(name string) string {
	return cli.Bold(cli.Underline(cli.Colorize(cli.White, name)))
}

func (c *ConfigWriter) Print() {
	if c.Title != "" {
		fmt.Print("\n", title(c.Title), "\n\n")
	}

	for _, s := range c.sections {
		fmt.Println(section(s.Name))

		if len(s.Values) == 0 {
			fmt.Println("  None")
		}

		tw := tablewriter.TableWriter{}
		for _, v := range s.Values {
			fmt.Fprintf(&tw, "  %s  \t%s", cli.Bold(v.Key), v.Value)
		}
		tw.Print()

		fmt.Println()
	}
}
