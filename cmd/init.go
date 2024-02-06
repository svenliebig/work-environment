/*
Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/utils"
	"github.com/svenliebig/work-environment/pkg/we"
)

var (
	initCmd = &cobra.Command{
		Use: "init",

		Short: "Initialize your work environemnt",
		Long: `Initialize your work environemnt in the current directory, it will search into the
subdirectories to find all git based projects there might be. After that, a work-environment
config directory will be created.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			p, err := utils.GetPath(args)

			if err != nil {
				log.Fatal(err)
			}

			o := &we.InitializeOptions{
				Override: cmd.Flag("override").Value.String() == "true",
			}

			err = we.Do(p, o)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates your work environemnt",
		Long: `Updates your work environment, it will search into the
subdirectories to find all git based projects there might be. After that, a work-environment
config directory will be updated by deleting projects that are not available anymore
and adding projects that are new.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = we.Update(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	weListCmd = &cobra.Command{
		Use:   "list",
		Short: "List your work environemnt",
		Long:  `Lists all your projects in your current path by default.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = we.List(ctx, &we.ListOptions{
				All:    cmd.Flag("all").Value.String() == "true",
				Tags:   cmd.Flag("tags").Value.String(),
				Filter: cmd.Flag("filter").Value.String(),
			})

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	useCmd = &cobra.Command{
		Use:   "use",
		Short: "Uses a work environment configuration to clone projects",
		Long: `Uses a work environment configuration to clone projects that are
defined in the there. This command should be used to redownload an
work environment after a work environment was created with 'we ci init'.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = we.Use(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("override", "o", false, "will override an existing work-environment configuration")

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(useCmd)

	rootCmd.AddCommand(weListCmd)
	weListCmd.Flags().BoolP("all", "a", false, "will list all projects instead of only the ones in the current path (no implemented yet)")
	weListCmd.Flags().StringP("filter", "f", "", "will filter the projects by the given string")
	weListCmd.Flags().StringP("tags", "t", "", "will filter the projects by the given tags (not implemented yet)")
}
