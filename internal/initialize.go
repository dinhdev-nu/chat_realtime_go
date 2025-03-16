package internal

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/router"
	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {

	r:= *router.InitRouter()

	
	return &r
}