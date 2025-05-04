package middlewares

import (
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		global.Log.Info("Resqest log::: ",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("error", c.Errors.String()),
		)
	}
}
