package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

type createRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=100" example:"Miguel"`
	Email   string `json:"email" binding:"required,email" example:"miguel@example.com"`
	Country string `json:"country" binding:"required,min=1,max=100" example:"Mexico"`
	City    string `json:"city" binding:"required,min=1,max=100" example:"CDMX"`
}

// @Summary Create a visitor
// @Description Create a new visitor record
// @Tags visitors
// @Accept json
// @Produce json
// @Param visitor body createRequest true "Visitor data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /visitors [post]
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
