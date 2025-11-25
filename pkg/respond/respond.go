package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorDetail represents a standardized error payload.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents the API error contract.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Error writes a standardized error response.
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
	c.Abort()
}

// ValidationError helper for binding/validation failures.
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "validation_error", message)
}

// JSON writes successful responses.
func JSON(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, payload)
}
