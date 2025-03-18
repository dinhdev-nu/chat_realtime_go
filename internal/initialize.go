package internal

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/router"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {

	r:= gin.Default()

	// Init middlewares
	// r.Use(middlewares.ErrorMiddleware())
	r.Use(middlewares.Cors())

	return router.InitRouter(r)


}