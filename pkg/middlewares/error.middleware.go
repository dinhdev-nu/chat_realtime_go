package middlewares

import (
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware is a middleware to handle error internal crash
func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.ServerError(ctx)
				ctx.Abort()
				return
			}
		}()
		ctx.Next()
	}
}
