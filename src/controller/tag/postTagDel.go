package tag

import (
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func PostTagDel(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var tag struct {
		TagId string `json:"tag_id"`
	}
	var resp struct {
		TagName string `json:"tag_name"`
		Code    string `json:"code"`
		Info    string `json:"info"`
	}
	//获取post数据
	err := ctx.BindJSON(&tag)
	if err != nil {
		resp.TagName = ""
		resp.Code = "404"
		resp.Info = "无法获取到post标签信息"
		loger.Println("[tag_del_info_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判空
	tagId := tag.TagId
	if tagId == "" {
		resp.TagName = ""
		resp.Code = "403"
		resp.Info = "无法获取到post标签tag_id"
		loger.Println("[tag_del_info_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判断是否存在
	if !sql.IsexistTagById(tagId) {
		resp.TagName = ""
		resp.Code = "404"
		resp.Info = "数据库没有这个标签"
		loger.Println("[tag_del_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//删除标签
	tname, err := sql.TagDelById(tagId)
	if err != nil {
		resp.TagName = ""
		resp.Code = "500"
		resp.Info = "数据库删除标签错误"
		loger.Println("[tag_del_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.TagName = tname
	resp.Code = "200"
	resp.Info = "标签删除成功"
	loger.Println("[tag_add_app_err]", resp)
	ctx.JSON(200, resp)
}
