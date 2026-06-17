package handlers

import (
	"errors"
	"net/http"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/gin-gonic/gin"
)

// BaseController provides common functionality for all controllers
type BaseController struct{}

func (bc *BaseController) JSONSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, dto.NewSuccess(data))
}

func (bc *BaseController) JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.NewError(code, message))
}

func (bc *BaseController) JSONValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, dto.ValidationErr(message))
}

func (bc *BaseController) JSONInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, dto.InternalErr(message))
}

func (bc *BaseController) JSONNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, dto.NotFoundErr(resource))
}

func (bc *BaseController) JSONUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, dto.UnauthorizedErr())
}

// Maps domain errors to HTTP responses.
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
