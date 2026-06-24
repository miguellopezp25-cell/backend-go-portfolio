package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *Server) *gin.Engine {
	r := gin.Default()

	// Configuración explícita y robusta de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",            // Tu local de Next.js
			"https://mlopezdev.up.railway.app", // ¡TU DOMINIO REAL DE FRONTEND!
			srv.cfg.Server.AllowedOrigins,      // Por si las dudas dejamos tu variable
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/healthz", HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/visitors", srv.List)
		v1.POST("/visitors", srv.Create)
		v1.GET("/visitors/:id", srv.GetByID)
	}

	return r
}
