package middleware

import (
	"strings"

	"github.com/RubenPari/clear-songs/internal/infrastructure/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// requestIDHeader is the HTTP header used to propagate request IDs.
const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware ensures each request has a unique request ID. It reads the
// X-Request-ID header if present, otherwise generates a new UUID. The ID is stored
// in the Gin context and echoed back in the response header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(logging.RequestIDKey, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}
