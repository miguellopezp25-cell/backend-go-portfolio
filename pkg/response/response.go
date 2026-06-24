// Package response proporciona una estructura consistente para todas las
// respuestas JSON de la API. Toda respuesta tiene {success, data?, error?,
// message?} para que el cliente siempre sepa interpretar el resultado.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse es el envelope JSON estándar. Success=true indica operación
// exitosa; Success=false indica error con Message describiendo el problema.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, detail interface{}) {
	resp := APIResponse{
		Success: false,
		Message: message,
	}
	if detail != nil {
		if s, ok := detail.(string); ok {
			resp.Error = s
		}
	}
	c.JSON(status, resp)
}
