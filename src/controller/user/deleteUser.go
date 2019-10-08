package user

import (
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func DeleteUser(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var userInfo struct {
		UserId string `json:"user_id"`
	}
	var resp struct {
		Code string `json:"code"`
		Info string `json:"info"`
	}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		resp.Code = "404"
		resp.Info = "无法获取post中的user_id"
		loger.Println("[user_delete]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取用户信息
	uid := userInfo.UserId
	//判空
	if uid == "" {
		resp.Code = "403"
		resp.Info = "post中user_id不可为空"
		loger.Println("[user_delete]", resp, err)
		ctx.JSON(200, resp)
		return
	}

	//判断用户是否存在
	if !sql.UserIsExistByUid(uid) {
		resp.Code = "501"
		resp.Info = "该用户不存在"
		loger.Println("[user_delete]", resp)
		ctx.JSON(200, resp)
		return
	}
	err = sql.DeleteUser(uid)
	if err != nil {
		resp.Code = "500"
		resp.Info = "数据库删除用户失败"
		loger.Println("[user_delete]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.Code = "200"
	resp.Info = "用户删除成功"
	loger.Println("[user_delete]", resp, err)
	ctx.JSON(200, resp)
}
