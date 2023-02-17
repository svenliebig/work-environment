package ci

import (
	gocontext "context"
	"fmt"
	"strings"

	"github.com/svenliebig/work-environment/pkg/api/bamboo"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/cli"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

// @Alex so before i had pointer everywhere, no I need to use the interface
// what are the implifications on the memory and data structure underneath this?
func Create(ctx context.BaseContext, url string, ciType string, name string, auth string) error {
	override := false
	config := ctx.Configuration()

	if config.HasCI(name) {
		q := fmt.Sprintf("\nThe Identifier '%s' is already declared in your configuration.\nDo you want to overwrite it? [y/n] ", cli.Colorize(cli.Purple, name))
		answer := cli.Question(q, []string{"y", "n"})
		if answer == "n" {
			fmt.Printf("%s the process.\n", cli.Colorize(cli.Red, "Abort"))
			return nil
		} else {
			override = true
		}
	}

	// TODO validate parameters
	if ciType != "bamboo" {
		fmt.Printf("the type %q is not a valid ci type\n", ciType)
		return nil
	}

	client := &bamboo.Client{
		BaseUrl:   url,
		AuthToken: auth,
	}

	version, err := client.GetInfo(gocontext.Background())

	if err != nil {
		return err
	}

	ci := &core.CI{
		CiType:     ciType,
		Identifier: name,
		AuthToken:  auth,
		Url:        url,
		Version:    version.Version,
	}

	if override {
		err = config.UpdateCI(ci)
	} else {
		err = config.AddCI(ci)
	}

	if err != nil {
		return err
	}

	err = ctx.Close()

	if err != nil {
		return err
	}

	fmt.Printf("\n%s added a new CI to your work environment:\n\n", cli.Colorize(cli.Green, "Successfully"))

	w := &tablewriter.TableWriter{}
	fmt.Fprintf(w, "  Identifier: \t%s", ci.Identifier)
	fmt.Fprintf(w, "  Type: \t%s", ci.CiType)
	fmt.Fprintf(w, "  URL: \t%s", ci.Url)
	fmt.Fprintf(w, "  Version: \t%s", version.Version)
	fmt.Fprintf(w, "  Token: \t%s", strings.Repeat("*", len(ci.AuthToken)))
	w.Print()
	fmt.Println()

	return nil
}
