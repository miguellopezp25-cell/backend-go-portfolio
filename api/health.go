package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miguel/go-back-portfolo/pkg/response"
)

// @Summary Health check
// @Description Returns service health status including database connectivity
// @Tags system
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 503 {object} response.APIResponse
// @Router /healthz [get]
func (s *Server) HealthCheck(c *gin.Context) {
	dbStatus := "up"
	if s.pool != nil {
		if err := s.pool.Ping(c.Request.Context()); err != nil {
			dbStatus = "down"
		}
	}

	if dbStatus == "down" {
		response.Error(c, http.StatusServiceUnavailable, "service degraded", gin.H{"database": dbStatus})
		return
	}

	response.OK(c, gin.H{
		"status":   "ok",
		"database": dbStatus,
	})
}
