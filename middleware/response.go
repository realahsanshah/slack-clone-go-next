package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error description"`
	Code    int         `json:"code" example:"200"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, data interface{}, message string, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    code,
	}
	c.JSON(code, response)
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	if statusCode > 499 {
		log.Println("Internal server error: ", err)
	}

	response := APIResponse{
		Success: false,
		Message: message,
		Code:    statusCode,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
	c.Abort()
}

// ResponseMiddleware handles global response formatting and error recovery
func ResponseMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Handle panics and convert them to 500 errors
		defer func() {
			if err := recover(); err != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
			}
		}()

		c.Next()

		// Handle cases where no response was sent but there was an error
		if c.Writer.Status() >= 400 && !c.Writer.Written() {
			switch c.Writer.Status() {
			case http.StatusNotFound:
				ErrorResponse(c, http.StatusNotFound, "Route not found", nil)
			case http.StatusMethodNotAllowed:
				ErrorResponse(c, http.StatusMethodNotAllowed, "Method not allowed", nil)
			default:
				ErrorResponse(c, c.Writer.Status(), "An error occurred", nil)
			}
		}
	})
}
