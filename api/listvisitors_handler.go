package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/schema"
)

// @Summary List all visitors
// @Description Retrieve a paginated list of visitors
// @Tags visitors
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} response.PaginatedResponse
// @Router /visitors [get]
func (s *Server) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	visitors, total, err := s.visitorService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list visitors", err.Error())
		return
	}

	if visitors == nil {
		visitors = []schema.Visitor{}
	}

	response.Paginated(c, visitors, total, page, pageSize)
}
