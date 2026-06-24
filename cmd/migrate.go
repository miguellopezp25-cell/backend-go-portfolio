package cmd

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/miguel/go-back-portfolo/config"
)

// migrateCmd ejecuta las migraciones pendientes usando golang-migrate.
// Las migraciones están en schema/migrations/ como archivos .up.sql y .down.sql.
// La convención de nombres (001_nombre) define el orden de ejecución.
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(cfgPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		m, err := migrate.New(
			"file://schema/migrations",
			cfg.Database.URL(),
		)
		if err != nil {
			log.Fatalf("Failed to create migrator: %v", err)
		}
		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migrations applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
