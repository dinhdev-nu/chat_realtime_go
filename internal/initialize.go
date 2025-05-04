package internal

import (
	"fmt"

	c "github.com/dinhdev-nu/realtime_auth_go/config"
	g "github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/router"
	m "github.com/dinhdev-nu/realtime_auth_go/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func Run() {

	// Config
	c.LoadConfig() // fmt.Print(global.Config.Server.Port) load config from environment
	c.InitLogger() // g.Log.Info("Server is starting...")
	c.InitMysql()
	c.InitRedis()

	// Init gin router
	var r *gin.Engine
	if g.Config.Server.Mode == "dev" {
		r = gin.Default() // log, recovery, cors default của gin
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New() // no log, no recovery
	}

	// Middlewares
	r.Use(m.Cors())
	r.Use(m.LoggerMidleware())
	r.Use(m.ErrorMiddleware())
	r.Use(m.RateLimitMiddleware())
	// r.Use(m.Authorzation())

	server := router.InitRouter(r)

	http := fmt.Sprintf("%s:%s", g.Config.Server.Host, g.Config.Server.Port)
	server.Run(http) // run server

}
