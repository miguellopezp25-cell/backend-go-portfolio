package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type SQLStore struct {
	pool *pgxpool.Pool
	*Queries
}

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}

// ExecTx executes fn within a database transaction.
func (s *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
