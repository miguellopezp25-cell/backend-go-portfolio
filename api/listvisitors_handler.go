package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/schema"
)

func (s *Server) List(c *gin.Context) {
	visitors, err := s.visitorService.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list visitors", err.Error())
		return
	}

	if visitors == nil {
		visitors = []schema.Visitor{}
	}

	response.OK(c, visitors)
}
