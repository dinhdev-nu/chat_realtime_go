package user

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/middleware/auth"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	uc *c.UserController
}

func NewUserRouter(uc *c.UserController) *UserRouter {
	return &UserRouter{
		uc: uc,
	}
}

func (ur *UserRouter) InitRoutes(router *gin.RouterGroup) {

	chatRouter := router.Group("user")
	chatRouter.Use(auth.AuthMiddleware())
	{
		chatRouter.GET("/info/:username", ur.uc.GetUserInfoByName)
		chatRouter.GET("/search", ur.uc.SearchUserByName)
	}
}
