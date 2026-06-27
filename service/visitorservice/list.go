package visitorservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/miguel/go-back-portfolo/schema"
	"github.com/miguel/go-back-portfolo/schema/db"
)

func (s *Service) List(ctx context.Context, page, pageSize int) ([]schema.Visitor, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := int32((page - 1) * pageSize)
	limit := int32(pageSize)

	total, err := s.store.CountVisitors(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count visitors: %w", err)
	}

	records, err := s.store.ListVisitors(ctx, db.ListVisitorsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list visitors: %w", err)
	}

	visitors := make([]schema.Visitor, 0, len(records))
	for _, v := range records {
		var visitor schema.Visitor
		if err := json.Unmarshal(v.Data, &visitor); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal visitor data: %w", err)
		}
		visitor.ID = v.ID.String()
		visitor.CreatedAt = v.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		visitors = append(visitors, visitor)
	}

	return visitors, total, nil
}
