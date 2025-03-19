package internal

import (
	"fmt"

	c "github.com/dinhdev-nu/realtime_auth_go/config"
	g "github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/router"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func Run() {

	// Config 
	c.LoadConfig(); // fmt.Print(global.Config.Server.Port)

	r:= gin.Default() // Init gin router

	// Init middlewares
	// r.Use(middlewares.ErrorMiddleware())
	r.Use(middlewares.Cors())


	server:= router.InitRouter(r)
 
	http:= fmt.Sprintf("%s:%s", g.Config.Server.Host, g.Config.Server.Port)
	server.Run(http) // run server

}