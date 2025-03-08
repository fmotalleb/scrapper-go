/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cobra"
)

var (
	browsers            []string
	driverDirectory     string
	onlyInstallShell    bool
	skipInstallBrowsers bool
	verbose             bool
	dryRun              bool
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs every integration needed for application to bootup",
	Run: func(cmd *cobra.Command, args []string) {
		playwright.Install(
			&playwright.RunOptions{
				OnlyInstallShell:    onlyInstallShell,
				DriverDirectory:     driverDirectory,
				SkipInstallBrowsers: skipInstallBrowsers,
				Verbose:             verbose,
				DryRun:              dryRun,
				Browsers:            browsers,
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringArrayVarP(&browsers, "browsers", "b", []string{"chromium"}, "Only Install Selected Browsers (chromium, firefox, webkit)")
	setupCmd.Flags().StringVarP(&driverDirectory, "driver-directory", "d", "", "where to put drivers (defaults to cache directory based on os)")
	setupCmd.Flags().BoolVar(&onlyInstallShell, "just-shell", false, "only install shell")
	setupCmd.Flags().BoolVar(&skipInstallBrowsers, "skip-browsers", false, "skip browser installation")
	setupCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	setupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry-run (wont install anything)")
}
