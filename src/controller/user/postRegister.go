package user

import (
	"articlebk/src/utils"
	"articlebk/src/utils/dbtable"
	"articlebk/src/utils/sql"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strconv"
)

//用户注册模块
func PostRegister(ctx *gin.Context) {
	//loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	//定义返回信息的map

	//获取post提交的用户数据
	var userInfo struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		RollerId string `json:"roller_id"`
	}
	var user dbtable.User
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		loger.Println(utils.LOG_USER_REGISTER_ERR, utils.RESP_INFO_GETJSONERR, err)
		utils.Resp(ctx, utils.RESP_CODE_GETJSONERR, utils.RESP_INFO_GETJSONERR, user)
		return
	}
	//获取json中的用户信息
	userName := userInfo.UserName
	userPwd := userInfo.Password
	roller := userInfo.RollerId
	rollerId, _ := strconv.Atoi(roller)
	//判断post的数据是否为空
	if userName == "" || userPwd == "" || roller == "" {
		loger.Println(utils.LOG_USER_REGISTER_ERR, utils.RESP_INFO_JSON_DATANULL)
		utils.Resp(ctx, utils.RESP_CODE_JSON_DATANULL, utils.RESP_INFO_JSON_DATANULL, user)
		return
	}

	//给明文密码加密
	md5code := md5.New()
	md5code.Write([]byte(userPwd))
	userPwdHash := hex.EncodeToString(md5code.Sum(nil))
	//用户结构体赋值
	userinfo := &dbtable.User{
		Name:         userName,
		PasswordHash: userPwdHash,
		RollerId:     rollerId,
	}
	//判断用户是否存在
	if sql.UserIsExistByName(userName) {
		loger.Println(utils.LOG_USER_REGISTER_ERR, utils.LOG_USER_ISEXISTS)
		utils.Resp(ctx, utils.RESP_CODE_DATAISEXISTS, utils.RESP_INFO_DATAISEXISTS, user)
		return
	}
	//插入数据库,并接收返回结果
	uname, uid, rid, err := sql.RegisterUser(userinfo)
	if err != nil {
		loger.Println(utils.LOG_USER_REGISTER_ERR, utils.DB_CREATE_ERR, err)
		utils.Resp(ctx, utils.RESP_CODE_CERATE_ERR, utils.DB_CREATE_ERR, user)
		return
	}

	//返回结果
	user.Name = uname
	uidInt, _ := strconv.Atoi(uid)
	user.Id = uidInt
	ridInt, _ := strconv.Atoi(rid)
	user.RollerId = ridInt
	loger.Println(utils.LOG_USER_CREATE_OK)
	utils.Resp(ctx, utils.RESP_CODE_OK, utils.RESP_INFO_OK, user)
}
