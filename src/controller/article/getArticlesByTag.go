package article

import (
	. "articlebk/src/common"
	"articlebk/src/common/dbtable"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func GetArtileListByTag(ctx *gin.Context) {

	var tag struct {
		TagId string `json:"tag_id"`
	}
	var resp struct {
		TagId    string             `json:"tag_id"`
		Articles []*dbtable.Article `json:"articles"`
		Code     string             `json:"code"`
		Info     string             `json:"info"`
	}

	err := ctx.BindJSON(&tag)
	if err != nil {
		resp.TagId = ""
		resp.Articles = nil
		resp.Code = "404"
		resp.Info = "无法从post获取tag_id"

		log.Println("article_select_by_tag_err", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取tagid
	tagId := tag.TagId
	//判空
	if tagId == "" {
		resp.TagId = ""
		resp.Articles = nil
		resp.Code = "403"
		resp.Info = "tag_id不可为空"
		Log.Error("article_select_by_tag_err", resp)
		ctx.JSON(200, resp)
		return
	}
	//是否存在此标签
	if !sql.IsexistTagById(tagId) {
		resp.TagId = ""
		resp.Articles = nil
		resp.Code = "501"
		resp.Info = "该标签不存在"
		Log.Warn("[article_select_by_tag_err]", resp)
		ctx.JSON(200, resp)
		return
	}
	//数据库查询
	articles, err := sql.GetArticlesByTagId(tagId)
	if err != nil {
		resp.TagId = ""
		resp.Articles = nil
		resp.Code = "500"
		resp.Info = "数据库无法查询该标题的文章"
		Log.Error("[article_select_by_tag_err]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	resp.TagId = tagId
	resp.Articles = articles
	resp.Code = "200"
	resp.Info = "查询文章成功"
	Log.Info("article_select_by_tag_info", resp)
	ctx.JSON(200, resp)
}
