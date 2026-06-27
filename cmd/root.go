// Package cmd implementa los comandos CLI con Cobra. Cada comando es un archivo
// separado (serve, migrate, version) para mantener el código organizado.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Visitor management API",
	Long:  "A REST API for managing visitor registrations with PostgreSQL and JSONB storage.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Flag global --config para todos los subcomandos
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "config/config.yaml", "path to configuration file")
}
