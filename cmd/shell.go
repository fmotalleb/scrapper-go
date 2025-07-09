// Package cmd contains cli interface logic
package cmd

import (
	"github.com/fmotalleb/scrapper-go/shell"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start scrapper in an interactive shell",
	Run: func(cmd *cobra.Command, args []string) {
		shell.RunShell()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
