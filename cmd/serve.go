/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fmotalleb/scrapper-go/log"
	"github.com/fmotalleb/scrapper-go/server"
	"github.com/spf13/cobra"
)

type serveArgs struct {
	address string
	port    uint32
}

var (
	serverArg serveArgs
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve service as an api endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.StartServer(fmt.Sprintf("%s:%d", serverArg.address, serverArg.port)); err != nil {
			slog.Error("error starting server", log.ErrVal(err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&serverArg.address, "address", "a", "127.0.0.1", "change this value if you want to expose server (since this app does not support authentication keep it behind a reverse proxy)")
	serveCmd.Flags().Uint32VarP(&serverArg.port, "port", "p", 8080, "port on which the service will be exposed (since this app does not support authentication keep it behind a reverse proxy)")
}
