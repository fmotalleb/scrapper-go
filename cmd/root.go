/*
Copyright © 2025 Motalleb Fallahnezhad

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	cfg          config.ExecutionConfig
	outputFormat utils.Output
	logLevel     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrapper-go",
	Short: "A Simple playwright wrapper that executes a simple yaml pipeline",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := log.SetupLogger(logLevel); err != nil {
			slog.Error("failed to set log level", slog.Any("err", err))
			panic(err)
		}
		slog.Debug("Level Set To", slog.String("level", logLevel))
	},
	Run: func(cmd *cobra.Command, args []string) {
		result, err := engine.ExecuteConfig(cfg)
		if err != nil {
			slog.Error("failed to execute command", slog.Any("err", err))
		}
		formatted, err := outputFormat.Format(result)
		if err != nil {
			slog.Error("failed to format", slog.Any("err", err))
			return
		}
		fmt.Fprint(os.Stdout, formatted)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.scrapper-go.yaml)")
	format := rootCmd.Flags().String("format", "json", "output format (json,yaml) defaults to json")
	outputFormat = utils.Output(*format)
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "WARN", "Log Level (DEBUG INFO WARN ERROR) set to DEBUG for verbose logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".scrapper-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".scrapper-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("using config file", slog.String("config", viper.ConfigFileUsed()))
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("unable to decode into struct: %v", slog.Any("err", err))
		panic("configuration error")
	}
	slog.Debug("loaded config file", slog.Any("config", cfg))
}
