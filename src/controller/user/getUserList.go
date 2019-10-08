package user

import (
	"articlebk/src/common/dbtable"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func GetUserList(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var userInfo struct {
		UserId   string `json:"user_id"`
		RollerId string `json:"roller_id"`
	}

	var resp struct {
		UserList []dbtable.User `json:"user_list"`
		Code     string         `json:"code"`
		Info     string         `json:"info"`
	}
	//获取post中的数据
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		resp.UserList = nil
		resp.Code = "404"
		resp.Info = "无法从post获取用户信息"
		loger.Println("[user_get_userList_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取值
	uid := userInfo.UserId
	rid := userInfo.RollerId
	//判空
	if uid == "" || rid == "" {
		resp.UserList = nil
		resp.Code = "404"
		resp.Info = "user_id或roller_id不能为空"
		loger.Println("[user_get_userList_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断用户是否为管理员
	ridint, _ := strconv.Atoi(rid)
	if ridint != 0 {
		resp.UserList = nil
		resp.Code = "403"
		resp.Info = "非管理员不能查询用户列表"
		loger.Println("[user_get_userList_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//用户是否存在
	if !sql.IsExistByUidRid(uid, rid) {
		resp.UserList = nil
		resp.Code = "504"
		resp.Info = "你提供的用户不存在,不能获取用户列表"
		loger.Println("[user_get_userList_err]", resp)
		ctx.JSON(200, resp)
		return
	}

	//查询用户列表
	users, err := sql.UserListGet()
	if err != nil {
		resp.UserList = nil
		resp.Code = "404"
		resp.Info = "无法从post获取用户信息"
		loger.Println("[user_get_userList_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	resp.UserList = users
	resp.Code = "200"
	resp.Info = "用户列表查询成功"
	loger.Println("[user_get_userList_err]", resp)
	ctx.JSON(200, resp)
}
