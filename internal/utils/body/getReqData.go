package body

import (
	"github.com/gin-gonic/gin"
)

// func Generics func giúp nhận và trả về một kiểu dữ liệu bất kỳ
// T là kiểu dữ liệu bất kỳ
func GetPayLoadFromRequestBody[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
