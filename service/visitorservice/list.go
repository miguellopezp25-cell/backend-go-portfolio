package visitorservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/miguel/go-back-portfolo/schema"
)

func (s *Service) List(ctx context.Context) ([]schema.Visitor, error) {
	records, err := s.store.ListVisitors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list visitors: %w", err)
	}

	visitors := make([]schema.Visitor, 0, len(records))
	for _, v := range records {
		var visitor schema.Visitor
		if err := json.Unmarshal(v.Data, &visitor); err != nil {
			return nil, fmt.Errorf("failed to unmarshal visitor data: %w", err)
		}
		visitor.ID = v.ID.String()
		visitors = append(visitors, visitor)
	}

	return visitors, nil
}
