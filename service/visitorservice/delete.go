package visitorservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	_, err := s.store.DeleteVisitor(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to delete visitor: %w", err)
	}

	return nil
}
