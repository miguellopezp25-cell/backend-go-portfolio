// Package database maneja la conexión a PostgreSQL usando pgxpool (pool de
// conexiones). El pool se configura con límites para evitar saturar la BD y
// se verifica con Ping antes de retornar.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/miguel/go-back-portfolo/config"
)

// Connect crea un pool de conexiones y verifica que la BD responda.
// pgxpool es thread-safe y recomendado para servers HTTP.
func Connect(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	// Límites conservadores para desarrollo; ajustar según carga en producción.
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
