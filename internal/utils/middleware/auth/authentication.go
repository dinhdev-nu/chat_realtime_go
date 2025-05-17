package auth

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/jwt"
	res "github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		token := c.Request.Header.Get("Authentication")
		userId := c.Request.Header.Get("Client-Id")

		if token == "" || userId == "" {

			token = c.GetString("token")
			userId = c.GetString("user_id")
			if token == "" || userId == "" {

				res.UnauthorizedError(c)
				c.Abort()
				return
			}
		}
		// verify token
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			res.UnauthorizedError(c)
			c.Abort()
			return
		}
		// check userId in token and header
		uuidToken := claims.Subject
		tokenId := uuidToken[:strings.Index(uuidToken, "token")]
		if tokenId != userId {
			res.UnauthorizedError(c)
			c.Abort()
			return
		}
		user, err := global.Rdb.Get(context.Background(), uuidToken).Result()
		if err != nil {
			res.BadRequestError(c, res.ServerErrorCode, res.CodeMessage[res.ServerErrorCode])
			c.Abort()
			return
		}
		userInfo := &model.GoDbUserInfo{}
		err = json.Unmarshal([]byte(user), userInfo)
		if err != nil {
			res.BadRequestError(c, res.ServerErrorCode, res.CodeMessage[res.ServerErrorCode])
			c.Abort()
			return
		}
		if id, _ := strconv.Atoi(userId); id != int(userInfo.UserID) {
			res.UnauthorizedError(c)
			c.Abort()
			return
		}

		c.Set("uuidToken", claims.Subject)
		c.Next()
	}
}
