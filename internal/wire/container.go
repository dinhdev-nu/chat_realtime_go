package wire

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	sa "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	sc "github.com/dinhdev-nu/realtime_auth_go/internal/service/chat"
	su "github.com/dinhdev-nu/realtime_auth_go/internal/service/user"
)

type Container struct {
	AuthController *c.AuthController
	ChatController *c.ChatController
	UserController *c.UserController
}

func NewContainer() *Container {

	// auth
	authRepo := repo.NewAuthRepo()
	authService := sa.NewAuthService(authRepo)
	authController := c.NewAuthController(authService)

	// user
	userRepo := repo.NewUserRepo()
	userService := su.NewUserService(userRepo)
	userController := c.NewUserController(userService)

	// chat
	chatRepo := repo.NewChatRepo()
	chatService := sc.NewChatService(chatRepo, authRepo, userRepo)
	chatController := c.NewChatController(chatService)

	return &Container{
		AuthController: authController,
		ChatController: chatController,
		UserController: userController,
	}

}
