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
		// TODO I don't know...
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{}, cobra.ShellCompDirectiveFilterDirs
		},
		Short: "Initialize your work environemnt",
		Long: `Initialize your work environemnt in the current directory, it will search into the
subdirectories to find all git based projects there might be. After that, a work-environment
config directory will be created.`,
		Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
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
		Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			p, err := utils.GetPath([]string{})

			if err != nil {
				log.Fatal(err)
			}

			ctx := &context.BaseContext{Cwd: p}

			err = ctx.Validate()

			if err != nil {
				log.Fatal(err)
			}

			err = we.Update(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(updateCmd)
	initCmd.Flags().BoolP("override", "o", false, "Will override an existing work-environment configuration")
}
