package router

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/router/auth"
	"github.com/gin-gonic/gin"
)

type RouterMain struct {
	Routers []Router
}

// Constructor cho RouterMain
// Thêm các router cho từng service vào đây
func NewRouterMain() *RouterMain {
	return &RouterMain{
		Routers: []Router{
			auth.NewAuthRouter(),
		},
	}	
}


// Hàn khỏi tại tất cả các route của RouterMain thế quy đinh interface Router
func (rm *RouterMain) InitRoutes(api *gin.RouterGroup) {
	for _, r := range rm.Routers {
		r.InitRoutes(api)
	}
}

func InitRouter() *gin.Engine {
	r := gin.Default() // init gin with default log middleware

	apiRoutes := r.Group("/v1/api")
	mainRouter := NewRouterMain()
	mainRouter.InitRoutes(apiRoutes)
	
	return r
}
