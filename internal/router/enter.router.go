package router

import "github.com/gin-gonic/gin"

type Router interface {
	// InitRoutes là phương thức khởi tạo tất cả các route của router
	// Hành vi chung các route của router run 
	InitRoutes(api *gin.RouterGroup)
}