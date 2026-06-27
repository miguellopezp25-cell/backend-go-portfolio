package visitorservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/schema"
	"github.com/miguel/go-back-portfolo/schema/db"
)

func (s *Service) Update(ctx context.Context, id string, req VisitorRequest) (*schema.Visitor, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal visitor data: %w", err)
	}

	v, err := s.store.UpdateVisitor(ctx, db.UpdateVisitorParams{
		ID:   uid,
		Data: data,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update visitor: %w", err)
	}

	var visitor schema.Visitor
	if err := json.Unmarshal(v.Data, &visitor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visitor data: %w", err)
	}
	visitor.ID = v.ID.String()

	return &visitor, nil
}
