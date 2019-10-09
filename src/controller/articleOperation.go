package controller

import (
	. "articlebk/src/common"
	"articlebk/src/common/database"
	"articlebk/src/common/database/sql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetArticleListByTag(ctx *gin.Context) {

	var tag struct {
		TagId string `json:"tag_id"`
	}
	err := ctx.BindJSON(&tag)
	if err != nil {
		Log.Error(LOG_ARTICLE_LIST_BYTAG_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	tagId := tag.TagId
	//判空
	if tagId == "" {
		Log.Error(LOG_ARTICLE_LIST_BYTAG_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否存在此标签
	if !sql.IsexistTagById(tagId) {
		Log.Error(LOG_ARTICLE_LIST_BYTAG_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	//数据库查询
	articles, err := sql.GetArticlesByTagId(tagId)
	if err != nil {
		Log.Error(LOG_ARTICLE_LIST_BYTAG_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_LISTBYTAG_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
}
func GetArticleListByColumn(ctx *gin.Context) {
	//sid := ctx.Query("sid")
	//database.ArticleSelectByColumn(sid)
	//TODO
}
func GetFailedArticleList(ctx *gin.Context) {
	//TODO
}
func GetReleasedArticleList(ctx *gin.Context) {
	//TODO
}
func GetWillBeReleaseArticleList(ctx *gin.Context) {
	//TODO
}
func PatchArticleEdit(ctx *gin.Context) {
	err := ctx.BindJSON("")
	if err != nil {

	}
	//TODO
}

func PostArticleAdd(ctx *gin.Context) {
	var article struct {
		Uid      string   `json:"uid"`
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		ColumnId string   `json:"column_id"`
		Tags     []string `json:"tags"`
	}
	//获取post的文章信息
	err := ctx.BindJSON(&article)
	if err != nil {
		Log.Error(LOG_ARTICLE_ADD_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取json中文章的主题/内容/专题id/标签的id数组(一个文章属于多个标签)
	uid := article.Uid
	uidInt, _ := strconv.Atoi(uid)
	content := article.Content
	title := article.Title
	columnId := article.ColumnId
	sid, _ := strconv.Atoi(columnId)
	tags := article.Tags
	//创建本地文件用于保存文章
	contentFile, err := ioutil.TempFile("/Users/ander/go/src/articlebk/articleData/articleFile/", "article-*.txt")
	if err != nil {
		Log.Error(LOG_ARTICLE_ADD_ERR, "创建文章文件出错!!", err)
		Resp(ctx, RESP_CODE_CERATE_ERR, "创建文章文件出错!!", nil)
	}
	contentPath := contentFile.Name()
	defer contentFile.Close()
	//创建Url提供用户访问
	fileSplit := strings.Split(contentFile.Name(), "/")
	filename := string(fileSplit[len(fileSplit)-1])
	contentUrl := "http://127.0.0.1:8888/articleFile/" + filename
	//给article表结构赋值
	articleTable := database.Article{
		UserId:      uidInt,
		Title:       title,
		ContentUrl:  contentUrl,
		ContentPath: contentPath,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Status:      0,
		ColumnId:    sid,
	}

	//判空
	if uid == "" || article.Title == "" || content == "" || len(tags) == 0 || columnId == "" {
		Log.Error(LOG_ARTICLE_ADD_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否已存在
	if sql.IsexistArticle(article.Title) {
		Log.Error(LOG_ARTICLE_ADD_ERR, RESP_INFO_DATAISEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISEXISTS, nil)
		return
	}
	//查看专栏是否存在
	if !sql.IsExistColumnBySid(columnId) {
		Log.Error(LOG_ARTICLE_ADD_ERR, "专栏不存在")
		Resp(ctx, RESP_CODE_DATAISEXISTS, "专栏不存在", nil)
		return
	}
	//查看标签是否存在
	for _, tag := range tags {
		if !sql.IsexistTagById(tag) {
			Log.Error(LOG_ARTICLE_ADD_ERR, "标签不存在")
			Resp(ctx, RESP_CODE_DATAISEXISTS, "标签不存在", nil)
			return
		}
	}

	art, err := sql.ArticleAdd(articleTable, tags)
	if err != nil {
		Log.Error(LOG_ARTICLE_ADD_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_CERATE_ERR, DB_CREATE_ERR, nil)
		return
	}
	_, err = contentFile.WriteString(content)
	if err != nil {
		Log.Error(LOG_ARTICLE_ADD_ERR, "无法创建文件, 无法保存数据内容", err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_CREATE_ERR, nil)
		return
	}
	Log.Info(LOG_ARTICLE_ADD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, art)
}

func DeleteArticleDel(ctx *gin.Context) {
	var delArticle struct {
		ArticleId string `json:"article_id"`
		UserId    string `json:"user_id"`
	}
	err := ctx.BindJSON(&delArticle)
	if err != nil {
		Log.Error(LOG_ARTICLE_DELETE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	articleId := delArticle.ArticleId
	aid, _ := strconv.Atoi(articleId)
	userId := delArticle.UserId
	uid, _ := strconv.Atoi(userId)
	if articleId == "" || userId == "" {
		Log.Error(LOG_ARTICLE_DELETE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断当前用户是否有权限
	if !sql.UserIsAdmin(userId) {
		//是否为作者本人
		if !sql.UserHasArticle(aid, uid) {
			Log.Error(LOG_ARTICLE_DELETE_ERR, RESP_INFO_NOPERMISSION)
			Resp(ctx, RESP_CODE_NOPERMISSION, RESP_INFO_NOPERMISSION, nil)
			return
		}
	}
	//删除数据库文章记录
	contentUrl, err := sql.ArticleDel(articleId)
	if err != nil {
		Log.Error(LOG_ARTICLE_DELETE_ERR, DB_DEL_ERR, err)
		Resp(ctx, RESP_CODE_DEL_ERR, DB_DEL_ERR, nil)
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
	Log.Info(LOG_ARTICLE_DEL_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articleId)
}

func PutArticleRelease(ctx *gin.Context) {
	//TODO
}

func PutArticleSubmit(ctx *gin.Context) {
	var articleSubmit struct {
		Aid string `json:"aid"`
		Uid string `json:"uid"`
	}
	//获取json中的uid和aid
	err := ctx.BindJSON(articleSubmit)
	if err != nil {
		Log.Error(LOG_ARTICLE_SUBMIT_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判断uid与aid是否为空
	userId := articleSubmit.Uid
	uid, _ := strconv.Atoi(userId)
	articleId := articleSubmit.Aid
	aid, _ := strconv.Atoi(articleId)
	if userId == "" || articleId == "" {
		Log.Error(LOG_ARTICLE_SUBMIT_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断文章和用户是否为绑定关系
	if !sql.UserHasArticle(aid, uid) {
		Log.Error(LOG_ARTICLE_SUBMIT_ERR, RESP_INFO_NOPERMISSION)
		Resp(ctx, RESP_CODE_NOPERMISSION, RESP_INFO_NOPERMISSION, nil)
		return
	}
	if err := sql.ArticleSubmit(aid, uid); err != nil {
		Log.Error(LOG_ARTICLE_SUBMIT_ERR, DB_ARTICLE_SUBMIT_ERR, err)
		Resp(ctx, RESP_CODE_SUBMIT_ERR, DB_ARTICLE_SUBMIT_ERR, nil)
		return
	}
	Log.Info(LOG_ARTICLE_SUBMIT_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}
