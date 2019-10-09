package user

import (
	. "articlebk/src/common"
	"articlebk/src/common/database"
	"articlebk/src/common/database/sql"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

//用户注册模块
func PostRegister(ctx *gin.Context) {
	//获取post提交的用户数据
	var userInfo struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		RollerId string `json:"roller_id"`
	}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		Log.Error(LOG_USER_REGISTER_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取json中的用户信息
	userName := userInfo.UserName
	userPwd := userInfo.Password
	roller := userInfo.RollerId
	rollerId, _ := strconv.Atoi(roller)
	//判断post的数据是否为空
	if userName == "" || userPwd == "" || roller == "" {
		Log.Error(LOG_USER_REGISTER_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}

	//给明文密码加密
	md5code := md5.New()
	md5code.Write([]byte(userPwd))
	userPwdHash := hex.EncodeToString(md5code.Sum(nil))
	//用户结构体赋值
	userinfo := &database.User{
		Name:         userName,
		PasswordHash: userPwdHash,
		RollerId:     rollerId,
	}
	//判断用户是否存在
	if sql.UserIsExistByName(userName) {
		Log.Error(LOG_USER_REGISTER_ERR, LOG_USER_ISEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
		return
	}
	//插入数据库,并接收返回结果
	dbUser, err := sql.RegisterUser(userinfo)
	if err != nil {
		Log.Error(LOG_USER_REGISTER_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_CERATE_ERR, DB_CREATE_ERR, nil)
		return
	}
	//返回结果
	Log.Info(LOG_USER_CREATE_OK, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, dbUser)
}

func DeleteUser(ctx *gin.Context) {
	var userInfo struct {
		UserId string `json:"user_id"`
	}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		Log.Error(LOG_USER_DELETE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取用户信息
	uid := userInfo.UserId
	//判空
	if uid == "" {
		Log.Error(LOG_USER_DELETE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}

	//判断用户是否存在
	if !sql.UserIsExistByUid(uid) {
		Log.Error(LOG_USER_DELETE_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	err = sql.DeleteUser(uid)
	if err != nil {
		Log.Error(LOG_USER_DELETE_ERR, DB_DEL_ERR, err)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_CREATE_ERR, nil)
		return
	}
	Log.Info(LOG_USER_DEL_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}

//用户登录模块
func PostLogin(ctx *gin.Context) {
	var userinfo struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	err := ctx.BindJSON(&userinfo)
	if err != nil {
		Log.Error(LOG_USER_LOGIN_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取用户名密码
	username := userinfo.UserName
	userPwd := userinfo.Password
	//判断用户名密码是否为空
	if username == "" || userPwd == "" {
		Log.Error(LOG_USER_LOGIN_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//给密码加密
	md5code := md5.New()
	md5code.Write([]byte(userPwd))
	userPwdHash := hex.EncodeToString(md5code.Sum(nil))
	//检测用户是否存在
	if !sql.UserIsExistByName(username) {
		Log.Error(LOG_USER_LOGIN_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}

	//验证登录信息
	user, err := sql.UserLogin(username, userPwdHash)
	if err != nil {
		Log.Error(LOG_USER_LOGIN_ERR, DB_DEL_ERR, err)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_CREATE_ERR, nil)
		return
	}
	//登录后保存session
	s := sessions.Default(ctx)
	s.Set("name", username)
	err = s.Save()
	if err != nil {

	}
	//登录成功
	Log.Info(LOG_USER_LOGIN_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, user)
}

func GetUserList(ctx *gin.Context) {
	var userInfo struct {
		UserId   string `json:"user_id"`
		RollerId string `json:"roller_id"`
	}
	//获取post中的数据
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		Log.Error(LOG_USER_LIST_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取值
	uid := userInfo.UserId
	rid := userInfo.RollerId
	//判空
	if uid == "" || rid == "" {
		Log.Error(LOG_USER_LIST_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断用户是否为管理员
	ridint, _ := strconv.Atoi(rid)
	if ridint != 0 {
		Log.Error(LOG_USER_LIST_ERR, RESP_INFO_NOPERMISSION)
		Resp(ctx, RESP_CODE_NOPERMISSION, RESP_INFO_NOPERMISSION, nil)
		return
	}
	//用户是否存在
	if !sql.IsExistByUidRid(uid, rid) {
		Log.Error(LOG_USER_LIST_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}

	//查询用户列表
	users, err := sql.UserListGet()
	if err != nil {
		Log.Error(LOG_USER_LIST_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_SELECT_ERR, nil)
		return
	}
	Log.Info(LOG_USER_LIST_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_SELECT, RESP_INFO_OK, users)
}
func PutUserRollerUpdate(ctx *gin.Context) {
	var userInfo struct {
		UserId      string `json:"user_id"`
		UserName    string `json:"user_name"`
		OldrollerId string `json:"old_roller_id"`
		NewrollerId string `json:"new_roller_id"`
	}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		Log.Error(LOG_USER_UPDATE_ROLLER_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
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
		Log.Error(LOG_USER_UPDATE_ROLLER_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断用户是否存在
	if !sql.UserIsExistByUid(userid) {
		Log.Error(LOG_USER_UPDATE_ROLLER_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	//更新用户角色
	user, err := sql.UserRollerUpdate(userid, username, newrollerid)
	if err != nil {
		Log.Error(LOG_USER_UPDATE_ROLLER_ERR, DB_UPDATE_ERR, err)
		Resp(ctx, RESP_CODE_UPDATE_ERR, DB_UPDATE_ERR, nil)
		return
	}
	//正确返回值
	Log.Info(LOG_USER_UPDATE_ROLLER_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, user)
}
func PutUserPwdUpdate(ctx *gin.Context) {
	var userInfo struct {
		UserId      string `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	//获取post中的用户信息
	err := ctx.BindJSON(&userInfo)
	if err != nil {

		return
	}
	//赋值判空
	uid := userInfo.UserId
	opwd := userInfo.OldPassword
	npwd := userInfo.NewPassword
	if uid == "" || opwd == "" || npwd == "" {

		return
	}
	//判断用户是否存在
	if !sql.UserIsExistByUid(uid) {

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

		return
	}
	err = sql.UpdateUserPwd(uid, npwdh5)
	if err != nil {

		return
	}
	resp.UserId = uid
	resp.Code = "200"
	resp.Info = "用户密码更新完成"
	loger.Println("[user_pwd_update_err]", resp, err)
	ctx.JSON(200, resp)
}
