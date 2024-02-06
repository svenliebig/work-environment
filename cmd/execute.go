/*
 *Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
 */
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/context"

	"github.com/svenliebig/work-environment/pkg/execute"
)

// ciCmd represents the ci command
var (
	executeCmd = &cobra.Command{
		Use:   "execute",
		Short: "executes predefined commands on one or multiple projects",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = execute.Do(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(executeCmd)
}
