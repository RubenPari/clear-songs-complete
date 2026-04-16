package logging

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const RequestIDKey = "request_id"

func InitFromEnv() *zap.Logger {
	level := parseLevel(strings.TrimSpace(os.Getenv("LOG_LEVEL")), zap.InfoLevel)
	stacktraceLevel := parseLevel(strings.TrimSpace(os.Getenv("LOG_STACKTRACE_LEVEL")), zap.ErrorLevel)
	format := strings.TrimSpace(strings.ToLower(os.Getenv("LOG_FORMAT")))

	cfg := zap.NewProductionConfig()
	if format == "console" {
		cfg.Encoding = "console"
	} else {
		cfg.Encoding = "json"
	}
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.DisableStacktrace = false
	cfg.EncoderConfig.TimeKey = "ts"

	logger, err := cfg.Build(zap.AddStacktrace(stacktraceLevel))
	if err != nil {
		logger = zap.NewExample()
	}

	zap.ReplaceGlobals(logger)
	return logger
}

func SafeSync(logger *zap.Logger) {
	if logger == nil {
		return
	}
	_ = logger.Sync()
}

func LoggerFromGinContext(c *gin.Context) *zap.Logger {
	if c == nil || c.Request == nil {
		return zap.L()
	}

	requestID := c.GetString(RequestIDKey)
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	return zap.L().With(
		zap.String("request_id", requestID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.FullPath()),
	)
}

func parseLevel(raw string, fallback zapcore.Level) zapcore.Level {
	if raw == "" {
		return fallback
	}

	var parsed zapcore.Level
	if err := parsed.Set(raw); err != nil {
		return fallback
	}
	return parsed
}
