package auth

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/middleware/auth"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	AuthControler *c.AuthController
}

func NewAuthRouter(ac *c.AuthController) *AuthRouter {
	//dependency injection thủ công
	// nên sử dung Wire
	return &AuthRouter{
		AuthControler: ac,
	}
}

func (ar *AuthRouter) InitRoutes(router *gin.RouterGroup) {

	// This group is for public route not need authentication
	authRouterGroupPublic := router.Group("auth")
	{
		authRouterGroupPublic.POST("/register", ar.AuthControler.Register)
		authRouterGroupPublic.POST("/send-otp", ar.AuthControler.SendOtp)
		authRouterGroupPublic.POST("/verify-otp", ar.AuthControler.VerifyOtp)
		authRouterGroupPublic.POST("/sign-up", ar.AuthControler.SignUp)
		authRouterGroupPublic.POST("/login", ar.AuthControler.Login)
	}

	// This group is for private route need authentication
	authRouterGroupPrivate := router.Group("auth")
	authRouterGroupPrivate.Use(auth.AuthMiddleware())
	{
		authRouterGroupPrivate.POST("/logout", ar.AuthControler.Logout)
	}
}
