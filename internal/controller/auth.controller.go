package controller

import (
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	res "github.com/dinhdev-nu/realtime_auth_go/pkg/response"
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
	res.SuccessResponse(ctx, ac.AuthService.Ping())
	// res.BadRequestError(ctx, 4001, "Error")
}