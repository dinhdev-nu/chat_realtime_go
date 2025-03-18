package middlewares

import (
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		token:= ctx.GetHeader("Authorization")
		if token != "valid-token" {
			response.UnauthorizedError(ctx, 3001, "Invalid Token")
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}