package api

import (
	"github.com/gin-gonic/gin"
	"github.com/miguel/go-back-portfolo/api/visitorhandler"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

// SetupRouter configura las rutas y middlewares de Gin. Se mantiene separado
// de los handlers para poder testear rutas de forma aislada si es necesario.
func SetupRouter(svc *visitorservice.Service) *gin.Engine {
	r := gin.Default()

	vh := visitorhandler.NewVisitorHandler(svc)

	r.GET("/healthz", HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/visitors", vh.Create)
		v1.GET("/visitors/:id", vh.GetByID)
	}

	return r
}
