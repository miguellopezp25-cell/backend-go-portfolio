package api

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Detectamos el origen de la petición
		origin := c.Request.Header.Get("Origin")

		// 2. Si viene de localhost o de tu dominio de producción, lo permitimos
		if origin == "http://localhost:3000" || origin == allowedOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 3. Configuramos el resto de los headers necesarios
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 4. ¡CRUCIAL! Si es un Preflight (OPTIONS), respondemos 204 y cortamos el flujo
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
