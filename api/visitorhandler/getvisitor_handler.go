package visitorhandler

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/pkg/response"
)

// GetByID extrae el UUID de la URL y lo pasa al servicio.
func (h *VisitorHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	visitor, err := h.svc.GetByID(c.Request.Context(), id)
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
