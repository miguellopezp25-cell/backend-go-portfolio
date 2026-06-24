package api

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/pkg/response"
)

func (s *Server) GetByID(c *gin.Context) {
	id := c.Param("id")

	visitor, err := s.visitorService.GetByID(c.Request.Context(), id)
	if err != nil {
		if stderrors.Is(err, apperrors.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "visitor not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get visitor", err.Error())
		return
	}

	response.OK(c, visitor)
}
