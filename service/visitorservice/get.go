package visitorservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/schema"
)

func (s *Service) GetByID(ctx context.Context, id string) (*schema.Visitor, error) {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %w", err)
	}

	v, err := 	s.store.GetVisitor(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get visitor: %w", err)
	}

	var visitor schema.Visitor
	if err := json.Unmarshal(v.Data, &visitor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visitor data: %w", err)
	}
	visitor.ID = v.ID.String()
	return &visitor, nil
}
