package visitorservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/miguel/go-back-portfolo/schema"
)

func (s *Service) Create(ctx context.Context, req VisitorRequest) (*schema.Visitor, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal visitor data: %w", err)
	}

	v, err := 	s.store.CreateVisitor(ctx, json.RawMessage(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create visitor: %w", err)
	}

	visitor := schema.Visitor{}
	err = json.Unmarshal(v.Data, &visitor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal created visitor data: %w", err)
	}
	visitor.ID = v.ID.String()

	return &visitor, nil
}
