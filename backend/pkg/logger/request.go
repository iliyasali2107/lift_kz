package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	requestLoggerKey = "logger"
)

// RequestWithLogger adds logger to request.
// It sets the logger l as a value in the Gin context with the key "logger".
// This allows other parts of the application to access the logger using the same context.
func RequestWithLogger(c *gin.Context, l *zap.Logger) {
	c.Set(requestLoggerKey, l)
}

// FromRequest returns logger from request.
func FromRequest(c *gin.Context) *zap.Logger {
	v, ok := c.Get(requestLoggerKey)
	if ok {
		if l, ok := v.(*zap.Logger); ok {
			return l
		}
	}

	return zap.L()
}
