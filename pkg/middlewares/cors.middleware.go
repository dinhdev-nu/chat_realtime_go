package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")                                                       //chấp nhận các domain client khác nhau
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")                        // chấp nhận các phương	thức
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Client-Id, Authentication") // chấp nhận các header khác nhau
		// c.Header("Access-Control-Allow-Credentials", "true")                                               // cho phép gửi cookie từ client

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200) // trả về status 204 nếu là preflight request
			return
		}

		c.Next()
	}
}
