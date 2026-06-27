package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/miguel/go-back-portfolo/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.NewServer(cfgPath)
		if err != nil {
			slog.Error("failed to create server", "error", err)
			os.Exit(1)
		}
		if err := srv.Start(); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
