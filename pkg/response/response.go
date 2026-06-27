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

// PaginatedResponse extiende APIResponse con metadatos de paginación.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func Paginated(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total / int64(pageSize))
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
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
