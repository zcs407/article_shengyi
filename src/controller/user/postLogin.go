package user

import (
	"articlebk/src/utils/sql"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

//用户登录模块
func PostLogin(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var userinfo struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	var resp struct {
		UserName string `json:"uer_name"`
		UserId   string `json:"user_id"`
		RollerId string `json:"roller_id"`
		Code     string `json:"code"`
		Info     string `json:"info"`
	}
	err := ctx.BindJSON(&userinfo)
	if err != nil {
		resp.UserName = ""
		resp.UserId = ""
		resp.RollerId = ""
		resp.Code = "404"
		resp.Info = "无法获取post中用户信息"
		loger.Println("[user_login_err]", "获取post的用户登录信息失败", err)
		ctx.JSON(200, resp)
		return
	}
	//获取用户名密码
	username := userinfo.UserName
	userPwd := userinfo.Password
	//判断用户名密码是否为空
	if username == "" || userPwd == "" {
		resp.UserName = ""
		resp.UserId = ""
		resp.RollerId = ""
		resp.Code = "403"
		resp.Info = "用户名密码不得为空"
		loger.Println("[user_login_err]", "用户名密码不能为空")
		ctx.JSON(200, resp)
	}
	//给密码加密
	md5code := md5.New()
	md5code.Write([]byte(userPwd))
	userPwdHash := hex.EncodeToString(md5code.Sum(nil))
	//检测用户是否存在
	if !sql.UserIsExistByName(username) {
		resp.UserName = ""
		resp.UserId = ""
		resp.RollerId = ""
		resp.Code = "402"
		resp.Info = "用户登录失败,用户不存在"
		loger.Println("[user_login_err]", "用户登录,数据库验证失败")
		ctx.JSON(200, resp)
		return
	}

	//验证登录信息
	uname, uid, rid, err := sql.UserLogin(username, userPwdHash)
	if err != nil {
		resp.UserName = ""
		resp.UserId = ""
		resp.RollerId = ""
		resp.Code = "500"
		resp.Info = "用户登录,数据库验证失败"
		loger.Println("[user_login_err]", "用户登录,数据库验证失败")
		ctx.JSON(200, resp)
		return
	}
	//登录后保存session
	s := sessions.Default(ctx)
	s.Set("name", username)
	err = s.Save()
	if err != nil {
		loger.Println("[user_login_err]", "无法保存session,请检查", err)
	}
	//登录成功
	resp.UserName = uname
	resp.UserId = uid
	resp.RollerId = rid
	resp.Code = "200"
	resp.Info = "用户登录成功"
	loger.Println("[user_login_info]", resp)
	ctx.JSON(200, resp)

}
