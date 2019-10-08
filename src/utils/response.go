package utils

import "github.com/gin-gonic/gin"

func Resp(ctx *gin.Context, code, info string, body interface{}) {
	var resp struct {
		Code string      `json:"code"`
		Info string      `json:"info"`
		Body interface{} `json:"body"`
	}
	resp.Code = code
	resp.Info = info
	resp.Body = body
	ctx.JSON(200, resp)
}
