package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard JSON envelope for all responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a 200 OK response.
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Paginated sends a paginated list response.
func Paginated(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
	})
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
	})
}

// NotFound sends a 404 Not Found response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
	})
}

// InternalError sends a 500 Internal Server Error response.
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: message,
	})
}

// UnprocessableEntity sends a 422 Validation Error response.
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
