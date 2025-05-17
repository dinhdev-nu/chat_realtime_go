package middlewares

import (
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

func Authorzation() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get header
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("key")
		}
		// get db

		// compare token with db

		if token != "valid-token" {
			response.ForbiddenError(c)
			c.Abort()
			return
		}
		c.Next()
	}

}
