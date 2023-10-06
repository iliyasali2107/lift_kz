package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mado/pkg/logger"
)

func logMiddleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := []zap.Field{ //zap.Field objects. These fields are key-value pairs that will be included in the log entry.

			zap.String("endpoint", c.Request.URL.Path),      //endpoint URL
			zap.String("method", c.Request.Method),          //HTTP method
			zap.String("remote_addr", c.Request.RemoteAddr), //remote address of the client
		}
		//this function adds the logger and the fields to the Gin context (c).
		// The With method of the Zap logger is used to add fields to the logger instance.
		logger.RequestWithLogger(c, l.With(fields...))

		c.Next()
	}
}
