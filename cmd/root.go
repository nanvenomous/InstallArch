/*
Copyright © 2025 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagVersion bool

	version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "InstallArch",
	Short: "CLI to help install Arch Linux on a computer",
	Long: `InstallArch is a CLI tool for installing Arch Linux.

You can run individual commands or use 'run-all' to execute all steps in sequence.
Each command can be run individually if a step fails, allowing you to debug and continue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagVersion {
			fmt.Println(version)
			return nil
		}

		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	version = ver
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "print the version of the cli")
}
