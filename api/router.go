package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/miguel/go-back-portfolo/api/middleware"
)

func SetupRouter(srv *Server) *gin.Engine {
	r := gin.Default()

	allowOrigins := []string{
		"http://localhost:3000",
		"https://mlopezdev.up.railway.app",
	}
	if srv.cfg != nil && srv.cfg.Server.AllowedOrigins != "" {
		allowOrigins = append(allowOrigins, srv.cfg.Server.AllowedOrigins)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.RequestID())
	r.Use(middleware.RateLimit(100, 200))

	r.GET("/healthz", srv.HealthCheck)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/visitors", srv.List)
		v1.POST("/visitors", srv.Create)
		v1.GET("/visitors/:id", srv.GetByID)
		v1.PUT("/visitors/:id", srv.Update)
		v1.DELETE("/visitors/:id", srv.Delete)
	}

	return r
}
