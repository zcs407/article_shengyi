package article

import (
	. "articlebk/src/common"
	"articlebk/src/common/sql"

	"github.com/gin-gonic/gin"
	"os"
)

func DeleteArticleDel(ctx *gin.Context) {
	var resp struct {
		Code string `json:"code"`
		Info string `json:"info"`
	}
	var delArticle struct {
		ArticleId string `json:"article_id"`
		UserId    string `json:"user_id"`
	}
	err := ctx.BindJSON(&delArticle)
	if err != nil {
		resp.Code = "404"
		resp.Info = "无法获取post中删除文章的信息"
		Log.Error("article_del_err", resp, err)
		ctx.JSON(200, resp)
		return
	}
	articleId := delArticle.ArticleId
	userId := delArticle.UserId
	if articleId == "" || userId == "" {
		resp.Code = "403"
		resp.Info = "无法获取,article_id或user_id"
		Log.Error("article_del_err", resp)
		ctx.JSON(200, resp)
		return
	}
	//判断当前用户是否有权限
	if !sql.UserIsAdmin(userId) {
		resp.Code = "501"
		resp.Info = "当前用户无权删除"
		Log.Warn("article_del_err", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//删除数据库文章记录
	contentUrl, err := sql.ArticleDel(articleId)
	if err != nil {
		resp.Code = "500"
		resp.Info = "数据库无法删除该文章"
		Log.Error("article_del_db_err", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取该文章对应的图片列表
	images := sql.ArticleImageDelByAid(articleId)
	//依次删除该文章本地图片
	for _, image := range images {
		imgPath := image.ImagePath
		_ = os.Remove(imgPath)
	}
	_ = os.Remove(contentUrl)
	resp.Code = "200"
	resp.Info = "文章删除成功"
	Log.Info("article_del_info", resp)
	ctx.JSON(200, resp)
}
