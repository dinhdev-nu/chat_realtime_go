package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}



// Sử lý các error cụ thể
func BadRequestError(ctx *gin.Context, code int, mes string){ // 400
	JSON(ctx, http.StatusBadRequest, code, mes, nil)
}

func UnauthorizedError(ctx *gin.Context, code int, mes string){ // 401
	JSON(ctx, http.StatusUnauthorized, InvalidToken, CodeMessage[InvalidToken], nil)
}

func ServerError(ctx *gin.Context){ // 500
	JSON(ctx, http.StatusInternalServerError, ServerErrorCode, CodeMessage[ServerErrorCode], nil)
}


func Error(ctx *gin.Context, httpCode int, code int, mes string){
	JSON(ctx, httpCode, code, mes, nil)
}

// SuccessRes ponse return success response
func SuccessResponse(ctx *gin.Context, data interface{}){
	JSON(ctx, http.StatusOK, SuccessCode, CodeMessage[SuccessCode], data)
}



// JSON return json response
func JSON(ctx *gin.Context, httpCode int, Code int, mes string, data interface{}){
	ctx.JSON(httpCode, &Response{
		Code: Code,
		Message: mes,
		Data: data,
	})
}