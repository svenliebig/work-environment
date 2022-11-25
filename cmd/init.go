/*
Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/initialize"
	"github.com/svenliebig/work-environment/pkg/utils"
)

var initCmd = &cobra.Command{
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

		o := &initialize.InitializeOptions{
			Override: cmd.Flag("override").Value.String() == "true",
		}

		err = initialize.Do(p, o)

		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("override", "o", false, "Will override an existing work-environment configuration")
}
