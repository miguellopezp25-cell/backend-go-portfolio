package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/miguel/go-back-portfolo/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := api.NewServer(cfgPath)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
