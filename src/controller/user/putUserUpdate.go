package user

import (
	"articlebk/src/common/sql"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

var (
	resp struct {
		UserName string `json:"uer_name"`
		UserId   string `json:"user_id"`
		rollerId string `json:"roller_id"`
		Code     string `json:"code"`
		Info     string `json:"info"`
	}
	loger = log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
)

func PutUserRollerUpdate(ctx *gin.Context) {
	var userInfo struct {
		UserId      string `json:"user_id"`
		UserName    string `json:"user_name"`
		OldrollerId string `json:"old_roller_id"`
		NewrollerId string `json:"new_roller_id"`
	}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		resp.UserName = ""
		resp.UserId = ""
		resp.rollerId = ""
		resp.Code = "404"
		resp.Info = "无法获取post更新用户角色的信息"
		loger.Println("[user_update_roller_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取用户信息
	userid := userInfo.UserId
	username := userInfo.UserName
	oldrollerid := userInfo.OldrollerId
	newrollerid := userInfo.NewrollerId
	rollernum, _ := strconv.Atoi(newrollerid)
	//判空
	if newrollerid == "" || username == "" || userid == "" || oldrollerid == "" || rollernum > 1 || rollernum < 0 {
		resp.UserName = ""
		resp.UserId = ""
		resp.rollerId = ""
		resp.Code = "403"
		resp.Info = "用户信息不能为空"
		loger.Println("[user_update_roller_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断用户是否存在
	if !sql.UserIsExistByUid(userid) {
		resp.UserName = ""
		resp.UserId = ""
		resp.rollerId = ""
		resp.Code = "403"
		resp.Info = "用户不存在,无法修改信息"
		loger.Println("[user_update_roller_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//更新用户角色
	uname, uid, urid, err := sql.UserRollerUpdate(userid, username, newrollerid)
	if err != nil {
		resp.UserName = ""
		resp.UserId = ""
		resp.rollerId = ""
		resp.Code = "500"
		resp.Info = "数据库变更用户角色失败"
		loger.Println("[user_update_roller_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//正确返回值
	resp.UserName = uname
	resp.UserId = uid
	resp.rollerId = urid
	resp.Code = "200"
	resp.Info = "用户更新角色成功"
	loger.Println("[user_update_roller_info]", resp)
	ctx.JSON(200, resp)
}
func PutUserPwdUpdate(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var resp struct {
		UserId string `json:"user_id"`
		Code   string `json:"code"`
		Info   string `json:"info"`
	}
	var userInfo struct {
		UserId      string `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	//获取post中的用户信息
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		resp.UserId = ""
		resp.Code = "404"
		resp.Info = "无法获取到post中用户信息"
		loger.Println("[user_pwd_update_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//赋值判空
	uid := userInfo.UserId
	opwd := userInfo.OldPassword
	npwd := userInfo.NewPassword
	if uid == "" || opwd == "" || npwd == "" {
		resp.UserId = ""
		resp.Code = "403"
		resp.Info = "无法获取到post中用户user_id,old_password,new_password"
		loger.Println("[user_pwd_update_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断用户是否存在
	if !sql.UserIsExistByUid(uid) {
		resp.UserId = uid
		resp.Code = "501"
		resp.Info = "该用户不存在"
		loger.Println("[user_pwd_update_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//给明文密码加密
	md5new := md5.New()
	dm5old := md5.New()
	md5new.Write([]byte(npwd))
	npwdh5 := hex.EncodeToString(md5new.Sum(nil))
	dm5old.Write([]byte(opwd))
	opwdh5 := hex.EncodeToString(dm5old.Sum(nil))
	//验证用户老密码是否正确
	if !sql.VerifyUserPwd(uid, opwdh5) {
		resp.UserId = uid
		resp.Code = "501"
		resp.Info = "用户密码更新失败,请检查输入的老密码是否正确"
		loger.Println("[user_pwd_update_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	err = sql.UpdateUserPwd(uid, npwdh5)
	if err != nil {
		resp.UserId = uid
		resp.Code = "500"
		resp.Info = "用户密码更新失败,数据库无法更新用户密码"
		loger.Println("[user_pwd_update_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.UserId = uid
	resp.Code = "200"
	resp.Info = "用户密码更新完成"
	loger.Println("[user_pwd_update_err]", resp, err)
	ctx.JSON(200, resp)
}
