package cmd

import (
	"log/slog"
	"strings"

	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cobra"
)

var (
	browsers            []string
	driverDirectory     string
	onlyInstallShell    bool
	skipInstallBrowsers bool
	dryRun              bool
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs every integration needed for application to bootup",
	Run: func(cmd *cobra.Command, args []string) {
		if err := playwright.Install(
			&playwright.RunOptions{
				OnlyInstallShell:    onlyInstallShell,
				DriverDirectory:     driverDirectory,
				SkipInstallBrowsers: skipInstallBrowsers,
				Verbose:             strings.ToLower(logLevel) == "debug",
				DryRun:              dryRun,
				Browsers:            browsers,
			},
		); err != nil {
			slog.Error("failed to install playwright's dependencies", slog.Any("err", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringArrayVarP(&browsers, "browsers", "b", []string{"chromium"}, "Only Install Selected Browsers (chromium, firefox, webkit)")
	setupCmd.Flags().StringVarP(&driverDirectory, "driver-directory", "d", "", "where to put drivers (defaults to cache directory based on os)")
	setupCmd.Flags().BoolVar(&onlyInstallShell, "just-shell", false, "only install shell")
	setupCmd.Flags().BoolVar(&skipInstallBrowsers, "skip-browsers", false, "skip browser installation")
	setupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry-run (wont install anything)")
}
