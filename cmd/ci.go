/*
Copyright © 2024 Sven Liebig
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/svenliebig/work-environment/pkg/ci"
	"github.com/svenliebig/work-environment/pkg/context"
)

// ciCmd represents the ci command
var (
	ciCmd = &cobra.Command{
		Use:   "ci",
		Short: "Configure and use continuous integrations environments",
		Long:  `Configure and use continuous integrations environments, create CI environments and add them to your project.`,
	}
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Adds a new CI to your work environment",
		Long:  `Adds a new CI to your work environment`,
		Run: func(cmd *cobra.Command, args []string) {
			u, err := cmd.Flags().GetString("url")

			if err != nil {
				log.Fatal(err)
			} else {
				ur, err := url.ParseRequestURI(u)

				if err != nil {
					log.Fatalf("the url parameter does not satisfy an url format\n%s", err)
				}

				u = (&url.URL{
					Scheme: ur.Scheme,
					Host:   ur.Host,
				}).String()
			}

			ciType, err := cmd.Flags().GetString("type")

			if err != nil {
				log.Fatal(err)
			}

			auth, err := cmd.Flags().GetString("auth")

			if err != nil {
				log.Fatal(err)
			}

			name, err := cmd.Flags().GetString("name")

			if err != nil {
				log.Fatal(err)
			}

			ctx, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = ci.Create(ctx, u, ciType, name, auth)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a CI to your project",
		Long: `Adds a CI to your project, you have to be inside the project path or specify the
project identifier. The CI identifier is required, when you have more than one CI specified in yur
work environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			ciId, err := cmd.Flags().GetString("ciIdentifier")

			if err != nil {
				log.Fatal(fmt.Errorf("err while trying to get variable %q. %w", "ciIdentifier", err))
			}

			bambooKey, err := cmd.Flags().GetString("key")

			if err != nil {
				log.Fatal(err)
			}

			suggest, err := cmd.Flags().GetBool("suggest")

			if err != nil {
				log.Fatal(err)
			}

			project, err := cmd.Flags().GetString("project")

			if err != nil {
				log.Fatal(err)
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				log.Fatal(err)
			}

			err = ci.Add(c, ciId, project, bambooKey, suggest)

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	openCmd = &cobra.Command{
		Use:   "open",
		Short: "Opens the CI environment of you current project path in the browser",
		Long: `Opens the CI environment of you current project path in the browser
the project path is your current working directory, the project needs to
have a CI configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				log.Fatal(err)
			}

			err = ci.Open(c)

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists the available CI environments in your work environment",
		Long:  `Lists the available CI environments in your work environment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := context.CreateBaseContext()

			if err != nil {
				log.Fatal(err)
			}

			err = ci.List(c)
			return err
		},
	}
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Lists the available information for the ci that is configured for the project",
		Long: `Lists the available information for the ci that is configured for the project,
the project path is your current working directory, the project needs to
have a CI configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				log.Fatal(err)
			}

			url, err := cmd.Flags().GetBool("url")

			if err != nil {
				log.Fatal(err)
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				log.Fatal(err)
			}

			err = ci.Info(c, &ci.InfoOptions{
				Url: url,
			})
			return err
		},
	}
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Removes the ci configuration from a project",
		Long: `Removes the ci configuration from a project,
the project path is your current working directory, the project needs to
have a CI configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := cmd.Flags().GetString("project")

			if err != nil {
				return err
			}

			c, err := context.CreateProjectContextWithProjectName(project)

			if err != nil {
				log.Fatal(err)
			}

			err = ci.Remove(c)
			return err
		},
	}

// 	resultCmd = &cobra.Command{
// 		Use:   "result",
// 		Short: "Prints the latest CI build results or the currently running build",
// 		Long: `Prints the latest CI build results or the currently running build,
// if the build is currently running. See also 'we ci results' to look at
// a history of build results.`,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := context.CreateProjectContext()

// 			if err != nil {
// 				log.Fatal(err)
// 			}

//			err = ci.Remove(c)
//			return err
//		},
//	}
)

func init() {
	rootCmd.AddCommand(ciCmd)
	ciCmd.AddCommand(createCmd)
	ciCmd.AddCommand(addCmd)
	ciCmd.AddCommand(openCmd)
	ciCmd.AddCommand(listCmd)
	ciCmd.AddCommand(infoCmd)
	ciCmd.AddCommand(removeCmd)

	ciCmd.PersistentFlags().StringP("project", "p", "", "The project where you want to execute your command. It's the current project folder by default.")

	createCmd.Flags().StringP("url", "u", "", "the URL of the CI you want to add\nexample: 'https://bamboo.company.com'")
	createCmd.Flags().StringP("type", "t", "", "the CI type, currently available types are 'bamboo'\nexample: 'bamboo'")
	createCmd.Flags().StringP("auth", "a", "", "your base64 auth token for the CI environment\nexample: '8fmiam7dm/2o3m8cunskeswefwe'")
	createCmd.Flags().StringP("name", "n", "", "the unique Identifier for the CI in your work environment\nexample: 'my-bamboo-ci'")

	createCmd.MarkFlagRequired("url")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("auth")
	createCmd.MarkFlagRequired("name")

	addCmd.Flags().StringP("ciIdentifier", "c", "", "the identifier of the ci\nexample: 'my-bamboo'")
	addCmd.Flags().BoolP("suggest", "s", false, "if set, you will get suggestions of bamboo project keys")
	addCmd.Flags().StringP("key", "b", "", "the key identifier for the project in the ci, not relevant if suggest is set\nexmaple: 'PRS-SZ'")

	infoCmd.Flags().BoolP("url", "u", false, "Prints the URL of the CI environment")

	// addCmd÷MarkFlagsMutuallyExclusive("suggest", "bambooKey")
}
