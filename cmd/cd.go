/*
Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/cd"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils"
)

// cdCmd represents the cd command
var (
	cdCmd = &cobra.Command{
		Use:   "cd",
		Short: "Configure and use continuous integrations",
		Long: `Configure and use continuous integrations, add specficis CI's to your project or
create a globally available work environment CI.`,
	}
	cdAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a CD to your project",
		Long: `Adds a CD to your project, you have to be inside the project path or specify the
project identifier. The CD identifier is required, when you have more than one CD specified in yur
work environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			p, err := utils.GetPath([]string{})
			c := &context.Context{
				Cwd: p,
			}

			if err != nil {
				log.Fatal(err)
			}

			err = c.Validate()

			if err != nil {
				log.Fatal(err)
			}

			err = cd.Add(c)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cdOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Opens the CD environment of you current project path in the browser",
		Long: `Opens the CD environment of you current project path in the browser
the project path is your current working directory, the project needs to
have a CD configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := utils.GetPath([]string{})
			c := &context.Context{
				Cwd: p,
			}

			if err != nil {
				log.Fatal(err)
			}

			err = c.Validate()

			if err != nil {
				log.Fatal(err)
			}

			err = cd.Open(c)

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}
	cdInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Lists the available information for the CD that is configured for the project",
		Long: `Lists the available information for the CD that is configured for the project,
the project path is your current working directory, the project needs to
have a CD configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := utils.GetPath([]string{})
			c := &context.Context{
				Cwd: p,
			}

			if err != nil {
				log.Fatal(err)
			}

			err = c.Validate()

			if err != nil {
				log.Fatal(err)
			}

			err = cd.Info(c)
			return err
		},
	}
	cdRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Removes the CD configuration from a project",
		Long: `Removes the CD configuration from a project,
the project path is your current working directory, the project needs to
have a CD configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := utils.GetPath([]string{})
			c := &context.Context{
				Cwd: p,
			}

			if err != nil {
				log.Fatal(err)
			}

			err = c.Validate()

			if err != nil {
				log.Fatal(err)
			}

			err = cd.Remove(c)
			return err
		},
	}
)

func init() {
	rootCmd.AddCommand(cdCmd)
	cdCmd.AddCommand(cdAddCmd)
	cdCmd.AddCommand(cdOpenCmd)
	cdCmd.AddCommand(cdRemoveCmd)
	cdCmd.AddCommand(cdInfoCmd)
}
