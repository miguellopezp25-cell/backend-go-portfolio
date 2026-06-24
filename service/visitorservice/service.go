package visitorservice

import (
	"github.com/miguel/go-back-portfolo/schema/db"
)

type VisitorRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{store: store}
}
