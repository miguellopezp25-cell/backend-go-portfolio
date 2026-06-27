package api

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

// @Summary Update a visitor
// @Description Update an existing visitor's data
// @Tags visitors
// @Accept json
// @Produce json
// @Param id path string true "Visitor UUID"
// @Param visitor body createRequest true "Updated visitor data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /visitors/{id} [put]
func (s *Server) Update(c *gin.Context) {
	id := c.Param("id")

	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid visitor id", "must be a valid UUID")
		return
	}

	var req visitorservice.VisitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	if req.Name == "" || req.Email == "" || req.Country == "" || req.City == "" {
		response.Error(c, http.StatusBadRequest, "all fields are required", nil)
		return
	}

	visitor, err := s.visitorService.Update(c.Request.Context(), id, req)
	if err != nil {
		if stderrors.Is(err, apperrors.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "visitor not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update visitor", err.Error())
		return
	}

	response.OK(c, visitor)
}
