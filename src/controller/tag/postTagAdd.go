package tag

import (
	"articlebk/src/common/database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func PostTagAdd(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var tag struct {
		TagName string `json:"tag_name"`
	}
	var resp struct {
		Tagid string `json:"tagid"`
		Code  string `json:"code"`
		Info  string `json:"info"`
	}
	//获取post中的标签名称
	err := ctx.BindJSON(&tag)
	if err != nil {
		resp.Tagid = ""
		resp.Code = "404"
		resp.Info = "无法获取post中的标签信息"
		loger.Println("[tag_add_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判空
	tagname := tag.TagName
	if tagname == "" {
		resp.Tagid = ""
		resp.Code = "403"
		resp.Info = "标签名不可为空"
		loger.Println("[tag_add_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判断是否存在
	if sql.IsexistTag(tagname) {
		resp.Tagid = ""
		resp.Code = "403"
		resp.Info = "标题已存在"
		loger.Println("[tag_add_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//判断数据库返状态
	tid, err := sql.TagAdd(tagname)
	if err != nil {
		resp.Tagid = ""
		resp.Code = "500"
		resp.Info = "数据库创建标签失败"
		loger.Println("[tag_add_app_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//返回标签id
	resp.Tagid = tid
	resp.Code = "200"
	resp.Info = "标签创建成功"
	loger.Println("[tag_add_app_info]", resp)
	ctx.JSON(200, resp)

}
