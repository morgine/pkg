package admin

import "github.com/gin-gonic/gin"

type MessageSender interface {
	SendData(ctx *gin.Context, data interface{})
	SendMsgSuccess(ctx *gin.Context, msg string)
	SendMsgWarning(ctx *gin.Context, msg string)
	SendMsgError(ctx *gin.Context, msg string)
	SendError(ctx *gin.Context, err error)
}
