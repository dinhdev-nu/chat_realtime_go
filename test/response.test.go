package test

import "github.com/gin-gonic/gin"

func Test(ctx *gin.Context) {

	a, err := 1, 2
	if err != 2 {
		ctx.JSON(200, gin.H{
			"code":    4001,
			"message": "Error",
		})
	}
	ctx.JSON(200, gin.H{
		"code":    2001,
		"message": "Success",
		"data":    a,
	})

}
