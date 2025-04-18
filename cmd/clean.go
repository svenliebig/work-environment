/*
Copyright © 2025 Sven Liebig <liebigsv@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/we"
)

var (
	clean = &cobra.Command{
		Use: "clean",

		Short: "Cleans your work environment",
		Long: `Cleans your work environment by removing target directories, node_modules and git ignores files from
		projects that are part of the subdirectories of the current path.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			o := &we.CleanOptions{
				Preview: cmd.Flag("preview").Value.String() == "true",
			}

			err = we.Clean(ctx, o)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(clean)
	clean.Flags().BoolP("preview", "p", false, "will preview the changes that would be made and the space that would be saved")
}
