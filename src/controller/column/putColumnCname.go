package column

import (
	"articlebk/src/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func PutColumncname(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var specialInfo struct {
		SpecialId      string `json:"special_id"`
		SpecialNewName string `json:"special_new_name"`
	}
	var resp struct {
		Sid  string `json:"sid"`
		Code string `json:"code"`
		Info string `json:"info"`
	}
	err := ctx.BindJSON(&specialInfo)
	if err != nil {
		resp.Sid = ""
		resp.Code = "404"
		resp.Info = "无法获取post中的专题信息"
		loger.Println("[special_cname_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取专题信息判空
	sid := specialInfo.SpecialId
	newName := specialInfo.SpecialNewName
	if sid == "" || newName == "" {
		resp.Sid = ""
		resp.Code = "403"
		resp.Info = "post中的special_id或special_new_name不能为空"
		loger.Println("[special_cname_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断专题是否存在
	if !sql.IsExistSpecialBySid(sid) {
		resp.Sid = sid
		resp.Code = "501"
		resp.Info = "该专题不存在,不可修改名称"
		loger.Println("[special_cname_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//修改专题名
	if err := sql.SpecialCname(sid, newName); err != nil {
		resp.Sid = sid
		resp.Code = "500"
		resp.Info = "数据库错误,专题无法重命名"
		loger.Println("[special_cname_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.Sid = sid
	resp.Code = "200"
	resp.Info = "专题重命名成功"
	loger.Println("[special_cname_err]", resp)
	ctx.JSON(200, resp)
	return
}
