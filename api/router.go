package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *Server) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/visitors", srv.List)
		v1.POST("/visitors", srv.Create)
		v1.GET("/visitors/:id", srv.GetByID)
	}



	return r
}
