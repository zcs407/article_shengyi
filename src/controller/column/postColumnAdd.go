package column

import (
	"articlebk/src/common/dbtable"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func PostColumnAdd(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var special struct {
		SpecleName string `json:"specle_name"`
		ParentId   string `json:"parent_id"`
	}
	var resp struct {
		SpecleName string `json:"specle_name"`
		ParentId   string `json:"parent_id"`
		Sid        string `json:"sid"`
		Code       string `json:"code"`
		Info       string `json:"info"`
	}
	if err := ctx.BindJSON(&special); err != nil {
		resp.SpecleName = ""
		resp.ParentId = ""
		resp.Sid = ""
		resp.Code = "404"
		resp.Info = "无法从post获取到专题内容"
		loger.Println("[special_add]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	specialName := special.SpecleName
	specialParentId := special.ParentId
	if specialName == "" {
		resp.SpecleName = ""
		resp.ParentId = ""
		resp.Sid = ""
		resp.Code = "403"
		resp.Info = "无法从post获取到专题名称"
		loger.Println("[special_add]", resp)
		ctx.JSON(200, resp)
		return
	}
	//如果没有父级id,则属于一级菜单,默认父级id是0
	if specialParentId == "" {
		specialParentId = "0"
	}
	specialPid, _ := strconv.Atoi(specialParentId)
	specialinfo := dbtable.Special{
		SpecialName: specialName,
		Pid:         specialPid,
	}
	if sql.IsexistSpecial(specialName) {
		resp.SpecleName = ""
		resp.ParentId = ""
		resp.Sid = ""
		resp.Code = "505"
		resp.Info = "专题已存在"
		loger.Println("[special_add]", resp)
		ctx.JSON(200, resp)
		return
	}
	//如果pid非0且不存在,则无法添加
	if specialParentId != "0" && !sql.IsExistSpecialBySid(specialParentId) {
		resp.SpecleName = ""
		resp.ParentId = ""
		resp.Sid = ""
		resp.Code = "501"
		resp.Info = "指定的pid,父级专题不存在"
		loger.Println("[special_add]", resp)
		ctx.JSON(200, resp)
		return
	}

	sid, err := sql.SpecialAdd(&specialinfo)
	if err != nil {
		resp.SpecleName = ""
		resp.ParentId = ""
		resp.Sid = ""
		resp.Code = "501"
		resp.Info = "数据错误,无法保存专题信息"
		loger.Println("[special_add]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.SpecleName = specialName
	resp.ParentId = specialParentId
	resp.Sid = sid
	resp.Code = "200"
	resp.Info = "专题创建成功"
	loger.Println("[special_add]", resp)
	ctx.JSON(200, resp)
}
