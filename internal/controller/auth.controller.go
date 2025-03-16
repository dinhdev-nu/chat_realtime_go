package controller

import (
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		AuthService: service.DefaultAuthService(),
	}
}
	
func (ac *AuthController) Ping(ctx *gin.Context) {
	ctx.JSON(200, ac.AuthService.Ping())
}