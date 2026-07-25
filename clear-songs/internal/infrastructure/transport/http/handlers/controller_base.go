// Package handlers implements the HTTP handler layer of the application.
// Each controller embeds BaseController for standardised JSON responses
// and delegates business logic to application-layer use cases.
package handlers

import (
	"errors"
	"net/http"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
)

// BaseController provides common JSON response helpers that all HTTP
// controllers embed. It standardises the response envelope and error
// codes returned by every endpoint.
type BaseController struct{}

// JSONSuccess writes a 200 OK response wrapping data in the standard
// success envelope produced by dto.NewSuccess.
func (bc *BaseController) JSONSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, dto.NewSuccess(data))
}

// JSONError writes an error response with the given HTTP status code,
// machine-readable error code, and human-readable message.
func (bc *BaseController) JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.NewError(code, message))
}

// JSONValidationError writes a 400 Bad Request response using the
// standard validation-error envelope.
func (bc *BaseController) JSONValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, dto.ValidationErr(message))
}

// JSONInternalError writes a 500 Internal Server Error response using
// the standard internal-error envelope.
func (bc *BaseController) JSONInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, dto.InternalErr(message))
}

// JSONNotFound writes a 404 Not Found response for the named resource
// using the standard not-found envelope.
func (bc *BaseController) JSONNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, dto.NotFoundErr(resource))
}

// JSONUnauthorized writes a 401 Unauthorized response using the
// standard unauthorized-error envelope.
func (bc *BaseController) JSONUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, dto.UnauthorizedErr())
}

// HandleDomainError inspects a domain-layer error and maps it to the
// appropriate HTTP status code and response envelope. It recognises
// validation, not-found, unauthorized, and external-API errors; any
// unrecognised error is treated as an internal server error.
func (bc *BaseController) HandleDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, shared.ErrValidation):
		bc.JSONValidationError(c, err.Error())
	case errors.Is(err, shared.ErrNotFound):
		bc.JSONNotFound(c, "Resource")
	case errors.Is(err, shared.ErrUnauthorized):
		bc.JSONUnauthorized(c)
	case errors.Is(err, shared.ErrExternalAPI):
		c.JSON(http.StatusBadGateway, dto.NewError("EXTERNAL_API_ERROR", err.Error()))
	default:
		bc.JSONInternalError(c, "An unexpected error occurred")
	}
}
