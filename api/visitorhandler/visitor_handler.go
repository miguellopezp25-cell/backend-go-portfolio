package visitorhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

// createRequest valida los campos HTTP antes de pasarlos al service.
type createRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Country string `json:"country" binding:"required"`
	City    string `json:"city" binding:"required"`
}

// VisitorHandler inyecta el servicio en lugar de crearlo internamente, lo que
// facilita testear el handler con distintos escenarios del servicio.
type VisitorHandler struct {
	svc *visitorservice.Service
}

func NewVisitorHandler(svc *visitorservice.Service) *VisitorHandler {
	return &VisitorHandler{svc: svc}
}

// Create valida el request con ShouldBindJSON (usa las tags binding del modelo)
// y retorna 201 con el UUID generado.
func (h *VisitorHandler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	visitor, err := h.svc.Create(c.Request.Context(), visitorservice.VisitorRequest{
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
