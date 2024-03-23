/*
Copyright © 2024 Sven Liebig
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/vcs"

	_ "github.com/svenliebig/work-environment/pkg/vcs/gitazuredevops"
)

var (
	vcsCmd = &cobra.Command{
		Use:   "vcs",
		Short: "Use this command to manage your version control system.",
		Long: `Use this command to manage your vcs (version control system).

	You can add new vcs, list all available vcs, remove a vcs or switch between vcs.`,
	}
	vcsCreate = &cobra.Command{
		Use:   "create",
		Short: "Adds a new vcs to your work environment",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = vcs.Create(ctx, vcs.CreateParameter{
				Type:        cmd.Flag("type").Value.String(),
				Identifier:  cmd.Flag("identifier").Value.String(),
				AccessToken: cmd.Flag("access-token").Value.String(),
				Url:         cmd.Flag("url").Value.String(),
			})

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	vcsOpen = &cobra.Command{
		Use:   "open",
		Short: "Opens the repository in the browser",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateProjectContextWithProjectName(cmd.Flag("project").Value.String())

			if err != nil {
				log.Fatal(err)
			}

			err = vcs.Open(ctx, vcs.OpenParameters{
				PullRequest: cmd.Flag("pull-request").Value.String() == "true",
			})

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	vcsInfo = &cobra.Command{
		Use:   "info",
		Short: "Shows information about the vcs",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateProjectContextWithProjectName(cmd.Flag("project").Value.String())

			if err != nil {
				log.Fatal(err)
			}

			err = vcs.Info(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	vcsAdd = &cobra.Command{
		Use:   "add",
		Short: "Adds a vcs to the current project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateProjectContextWithProjectName(cmd.Flag("project").Value.String())

			if err != nil {
				log.Fatal(err)
			}

			err = vcs.Add(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	vcsRemove = &cobra.Command{
		Use:   "remove",
		Short: "Removes the configured vcs from the current project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := context.CreateProjectContextWithProjectName(cmd.Flag("project").Value.String())

			if err != nil {
				log.Fatal(err)
			}

			err = vcs.Remove(ctx)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(vcsCmd)
	vcsCmd.AddCommand(vcsCreate)
	vcsCmd.AddCommand(vcsInfo)
	vcsCmd.AddCommand(vcsAdd)
	vcsCmd.AddCommand(vcsRemove)
	vcsCmd.AddCommand(vcsOpen)

	vcsCmd.PersistentFlags().StringP("project", "p", "", "The project where you want to execute your command. It's the current project folder by default.")

	vcsCreate.Flags().StringP("type", "t", "", fmt.Sprintf("The type of the vcs. Can be one of: '%s'", strings.Join(vcs.AvailableClients(), "', '")))
	vcsCreate.Flags().StringP("identifier", "i", "", "The identifier of the vcs.")
	vcsCreate.Flags().StringP("access-token", "a", "", "The access token of the vcs.")
	vcsCreate.Flags().StringP("url", "u", "", "The url of the vcs. For example https://dev.azure.com/{organization}/{project}")

	vcsCreate.MarkFlagRequired("type")
	vcsCreate.MarkFlagRequired("identifier")
	vcsCreate.MarkFlagRequired("access-token")
	vcsCreate.MarkFlagRequired("url")

	vcsOpen.Flags().BoolP("pull-request", "r", false, "Opens the pull request page instead of the repository page.")
}
