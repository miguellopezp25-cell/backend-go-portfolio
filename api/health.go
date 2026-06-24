package api

import (
	"github.com/gin-gonic/gin"
	"github.com/miguel/go-back-portfolo/pkg/response"
)

// HealthCheck es un endpoint sencillo para que balanceadores de carga y
// orquestadores verifiquen que el servicio responde.
func HealthCheck(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}
