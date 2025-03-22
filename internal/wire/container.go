package wire

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
)

type Container struct {
	AuthController *c.AuthController
}

func NewContainer() *Container {

	authRepo := repo.NewAuthRepo()
	authService := service.NewAuthService(authRepo)
	authController := c.NewAuthController(authService)

	return &Container{
		AuthController: authController,
	}

}
