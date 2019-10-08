package tag

import (
	"articlebk/src/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func PostTagCname(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var tag struct {
		TagId   string `json:"tag_id"`
		NewName string `json:"new_name"`
	}
	var resp struct {
		Tagid   string `json:"tagid"`
		TagName string `json:"tag_name"`
		Code    string `json:"code"`
		Info    string `json:"info"`
	}
	//获取post数据
	err := ctx.BindJSON(&tag)
	if err != nil {
		resp.Tagid = ""
		resp.TagName = ""
		resp.Code = "404"
		resp.Info = "无法获取到post标签信息"
		loger.Println("[tag_cname_info_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判空
	tagId := tag.TagId
	tagNewName := tag.NewName
	if tagId == "" || tagNewName == "" {
		resp.Tagid = ""
		resp.TagName = ""
		resp.Code = "403"
		resp.Info = "无法获取到post标签tag_id或tag_name"
		loger.Println("[tag_cname_info_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判断是否存在
	if !sql.IsexistTagById(tagId) {
		resp.Tagid = tagId
		resp.TagName = ""
		resp.Code = "404"
		resp.Info = "数据库没有这个标签"
		loger.Println("[tag_cname_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	tname, err := sql.TagCnameById(tagId, tagNewName)
	if err != nil {
		resp.Tagid = tagId
		resp.TagName = ""
		resp.Code = "500"
		resp.Info = "数据库修改标签错误"
		loger.Println("[tag_cname_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.Tagid = tagId
	resp.TagName = tname
	resp.Code = "200"
	resp.Info = "标签修改成功"
	loger.Println("[tag_cname_app_info]", resp)
	ctx.JSON(200, resp)

}
