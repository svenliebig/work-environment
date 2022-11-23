/*
Copyright © 2022 Sven Liebig <liebigsv@gmail.com>
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "we",
	Short: "the work environment cli help you to organize and maintain a productive work environment",
	Long: `the work environment cli help you to organize and maintain a productive work environment
by providing a link between the tool you use the most (the cli) and the applications
that you have to use (like CI, CD, Dev/QA Stages, etc).

The goal is, to get less disrupted in your workflow and more time for important 
things, rather than waiting for lists of ci plans to fetch or waiting for bloated
web applications (*cough cough* jira *cough*) to execute a base functionality.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
