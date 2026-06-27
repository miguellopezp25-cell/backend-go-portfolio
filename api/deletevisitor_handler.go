package api

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/pkg/response"
)

// @Summary Delete a visitor
// @Description Delete a visitor by their UUID
// @Tags visitors
// @Produce json
// @Param id path string true "Visitor UUID"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /visitors/{id} [delete]
func (s *Server) Delete(c *gin.Context) {
	id := c.Param("id")

	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid visitor id", "must be a valid UUID")
		return
	}

	if err := s.visitorService.Delete(c.Request.Context(), id); err != nil {
		if stderrors.Is(err, apperrors.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "visitor not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete visitor", err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
