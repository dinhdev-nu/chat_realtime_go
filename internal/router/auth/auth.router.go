package auth

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{
	AuthControler *c.AuthController
}

func NewAuthRouter() *AuthRouter {
	
	authService:= service.DefaultAuthService()
	authController:= c.NewAuthController(authService)
	// nên sử dung Wire 

	return &AuthRouter{
		AuthControler: authController,
	}
}

func (ar *AuthRouter) InitRoutes(router *gin.RouterGroup) {

	
	
	// This group is for public route not need authentication
	authRouterGroupPublic := router.Group("auth")
	{
		authRouterGroupPublic.GET("/pong", func(ctx *gin.Context) {
			ctx.JSON(200, map[string]any{
				"message": "ping",
			})
		})
	}

	// This group is for private route need authentication
	authRouterGroupPrivate := router.Group("auth")
	{
		authRouterGroupPrivate.GET("/ping", ar.AuthControler.Ping)
	}
}


