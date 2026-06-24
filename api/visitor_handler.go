package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

type createRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Country string `json:"country" binding:"required"`
	City    string `json:"city" binding:"required"`
}

func (s *Server) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	visitor, err := s.visitorService.Create(c.Request.Context(), visitorservice.VisitorRequest{
		Name:    req.Name,
		Email:   req.Email,
		Country: req.Country,
		City:    req.City,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create visitor", err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":      visitor.ID,
		"message": "visitor created",
	})
}
