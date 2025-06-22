package router

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/router/auth"
	"github.com/dinhdev-nu/realtime_auth_go/internal/router/chat"
	"github.com/dinhdev-nu/realtime_auth_go/internal/router/user"
	di "github.com/dinhdev-nu/realtime_auth_go/internal/wire"
	"github.com/gin-gonic/gin"
)

type RouterMain struct {
	Routers []Router
}

// Constructor cho RouterMain
// Thêm các router cho từng service vào đây
func newRouterMain() *RouterMain {
	container := di.NewContainer() // dependency injection
	return &RouterMain{
		Routers: []Router{
			auth.NewAuthRouter(container.AuthController),
			chat.NewChatRouter(container.ChatController),
			user.NewUserRouter(container.UserController),
		},
	}
}

// Hàm khởi tạo tất cả các route của RouterMain thế quy đinh interface Router
func (rm *RouterMain) initRoutes(api *gin.RouterGroup) {
	for _, r := range rm.Routers {
		r.InitRoutes(api)
	}
}

func InitRouter(r *gin.Engine) *gin.Engine {

	apiRoutes := r.Group("/v1/api")
	mainRouter := newRouterMain()
	mainRouter.initRoutes(apiRoutes)

	return r
}
