package utils

import (
	"net/http"
	"slack-clone-go-next/middleware"

	"github.com/gin-gonic/gin"
)

// Success responses
func OK(c *gin.Context, data interface{}, message string) {
	response := middleware.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func Created(c *gin.Context, data interface{}, message string) {
	response := middleware.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusCreated, response)
}

// Error responses
func BadRequest(c *gin.Context, message string, err error) {
	middleware.ErrorResponse(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	middleware.ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	middleware.ErrorResponse(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	middleware.ErrorResponse(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c *gin.Context, message string, err error) {
	middleware.ErrorResponse(c, http.StatusInternalServerError, message, err)
}

func UnprocessableEntity(c *gin.Context, message string, err error) {
	middleware.ErrorResponse(c, http.StatusUnprocessableEntity, message, err)
}
