/*
Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/cd"
	"github.com/svenliebig/work-environment/pkg/context"
)

// cdCmd represents the cd command
var (
	cdCmd = &cobra.Command{
		Use:   "cd",
		Short: "Configure and use continuous delivery environments",
		Long:  `Configure and use continuous delivery environments, create CD environments and add them to your project.`,
	}
	cdAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a CD to your project",
		Long: `Adds a CD to your project, you have to be inside the project path or specify the
project identifier. The CD identifier is required, when you have more than one CD specified in yur
work environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				return err
			}

			return cd.Add(c)
		},
	}
	cdOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Opens the CD environment of you current project path in the browser",
		Long: `Opens the CD environment of you current project path in the browser
the project path is your current working directory, the project needs to
have a CD configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				return err
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
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			url, err := cmd.Flags().GetBool("url")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				return err
			}

			err = cd.Info(c, &cd.InfoOptions{
				Url: url,
			})
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
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				return err
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

	cdCmd.PersistentFlags().StringP("project", "p", "", "The project where you want to execute your command. It's the current project folder by default.")

	cdInfoCmd.Flags().BoolP("url", "u", false, "Prints the URL of the CD environment")
}
