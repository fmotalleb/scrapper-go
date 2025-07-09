package cmd

import (
	"log/slog"
	"strings"

	"github.com/fmotalleb/scrapper-go/log"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cobra"
)

type setupArgs struct {
	browsers            []string
	driverDirectory     string
	onlyInstallShell    bool
	skipInstallBrowsers bool
	dryRun              bool
}

var setupArg setupArgs

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs every integration needed for application to bootup",
	Run: func(cmd *cobra.Command, args []string) {
		if err := playwright.Install(
			&playwright.RunOptions{
				OnlyInstallShell:    setupArg.onlyInstallShell,
				DriverDirectory:     setupArg.driverDirectory,
				SkipInstallBrowsers: setupArg.skipInstallBrowsers,
				Verbose:             strings.ToLower(logLevel) == "debug",
				DryRun:              setupArg.dryRun,
				Browsers:            setupArg.browsers,
			},
		); err != nil {
			slog.Error("failed to install playwright's dependencies", log.ErrVal(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringArrayVarP(&setupArg.browsers, "browsers", "b", []string{"chromium"}, "Only Install Selected Browsers (chromium, firefox, webkit)")
	setupCmd.Flags().StringVarP(&setupArg.driverDirectory, "driver-directory", "d", "", "where to put drivers (defaults to cache directory based on os)")
	setupCmd.Flags().BoolVar(&setupArg.onlyInstallShell, "just-shell", false, "only install shell")
	setupCmd.Flags().BoolVar(&setupArg.skipInstallBrowsers, "skip-browsers", false, "skip browser installation")
	setupCmd.Flags().BoolVar(&setupArg.dryRun, "dry-run", false, "dry-run (wont install anything)")
}
