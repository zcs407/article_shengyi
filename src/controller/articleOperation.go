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

//新增文章
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
	cid, _ := strconv.Atoi(columnId)
	tags := article.Tags
	//创建本地文件用于保存文章
	contentFile, err := ioutil.TempFile(Settings.FileServer.TextPath, "article-*.txt")
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
		ColumnId:    cid,
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

//删除文章
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
	if len(articleId) == 0 || len(userId) == 0 {
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
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}

//获取编辑信息
func GetArticleForEdit(ctx *gin.Context) {
	respBody := make(map[string]interface{})
	var getArticle struct {
		UserId    string `json:"user_id"`
		ArticleId string `json:"article_id"`
	}
	err := ctx.BindJSON(&getArticle)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	userId := getArticle.UserId
	uid, _ := strconv.Atoi(userId)
	articleId := getArticle.ArticleId
	aid, _ := strconv.Atoi(articleId)
	if len(userId) == 0 || len(articleId) == 0 {
		Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断当前用户是否有权限
	if !sql.UserIsAdmin(userId) {
		//是否为作者本人
		if !sql.UserHasArticle(aid, uid) {
			Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, RESP_INFO_NOPERMISSION)
			Resp(ctx, RESP_CODE_NOPERMISSION, RESP_INFO_NOPERMISSION, nil)
			return
		}
	}
	//根据aid获取文章信息
	article, err := sql.GetArticleByAid(aid)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
		return
	}
	//获取所有专栏
	columns, err := sql.GetColumnList()
	if err != nil {
		Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
		return
	}
	//获取所有标签
	tags, err := sql.GetTagList()
	if err != nil {
		Log.Error(LOG_ARTICLE_GETFOREDIT_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
		return
	}
	//读取文件内容并赋值
	content, err := ioutil.ReadFile(article.ContentPath)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETFOREDIT_SUCCESS, RESP_INFO_READFILE_ERR, err)
		Resp(ctx, RESP_CODE_READFILE_ERR, RESP_INFO_READFILE_ERR, nil)
		return
	}
	//给返回体赋值
	respBody["article"] = article
	respBody["article.content"] = string(content)
	respBody["columns"] = columns
	respBody["tags"] = tags
	Log.Error(LOG_ARTICLE_GETFOREDIT_SUCCESS, DB_SELECT_ERR)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, respBody)
}

//提交编辑信息
func PatchArticleEdit(ctx *gin.Context) {
	var updateArticle struct {
		Uid       string   `json:"uid"`
		ArticleId string   `json:"article_id"`
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		ColumnId  string   `json:"column_id"`
		Tags      []string `json:"tags"`
	}
	//获取post的文章信息
	err := ctx.BindJSON(&updateArticle)
	if err != nil {
		Log.Error(LOG_ARTICLE_EDIT_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取json中文章的主题/内容/专题id/标签的id数组(一个文章属于多个标签)
	userId := updateArticle.Uid
	uid, _ := strconv.Atoi(userId)
	content := updateArticle.Content
	title := updateArticle.Title
	columnId := updateArticle.ColumnId
	cid, _ := strconv.Atoi(columnId)
	tags := updateArticle.Tags
	articleId := updateArticle.ArticleId
	aid, _ := strconv.Atoi(articleId)
	//判空
	if len(userId) == 0 || len(title) == 0 || len(content) == 0 || len(tags) == 0 || len(columnId) == 0 || len(articleId) == 0 {
		Log.Error(LOG_ARTICLE_EDIT_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//如果是管理员则无需检查
	if !sql.UserIsAdmin(userId) {
		//查看用户是否有该文章
		if !sql.UserHasArticle(aid, uid) {
			Log.Error(LOG_ARTICLE_EDIT_ERR, RESP_INFO_DATAISNOTEXISTS)
			Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
			return
		}
	}
	//查看专栏是否存在
	if !sql.IsExistColumnBySid(columnId) {
		Log.Error(LOG_ARTICLE_EDIT_ERR, "专栏不存在")
		Resp(ctx, RESP_CODE_DATAISEXISTS, "专栏不存在", nil)
		return
	}
	//查看标签是否存在
	for _, tag := range tags {
		if !sql.IsexistTagById(tag) {
			Log.Error(LOG_ARTICLE_EDIT_ERR, "标签不存在")
			Resp(ctx, RESP_CODE_DATAISEXISTS, "标签不存在", nil)
			return
		}
	}
	//获取该文章信息
	article, err := sql.GetArticleByAid(aid)
	if err != nil {
		Log.Error(LOG_ARTICLE_EDIT_ERR, DB_SELECT_ERR, err)
		Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
		return
	}
	//给新的文章赋值
	newArticle := database.Article{
		Id:          article.Id,
		UserId:      article.UserId,
		Title:       title,
		ContentUrl:  article.ContentUrl,
		ContentPath: article.ContentPath,
		CreateTime:  article.CreateTime,
		UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Images:      article.Images,
		Status:      0,
		ColumnId:    cid,
	}
	art, err := sql.ArticleUpdate(newArticle, tags, aid)
	if err != nil {
		Log.Error(LOG_ARTICLE_EDIT_ERR, DB_CREATE_ERR, err)
		Resp(ctx, RESP_CODE_CERATE_ERR, DB_CREATE_ERR, nil)
		return
	}
	contentFile, err := os.Open(article.ContentPath)
	if err != nil {
		Log.Error(LOG_ARTICLE_EDIT_ERR, RESP_INFO_READFILE_ERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_READFILE_ERR, nil)
	}
	_, err = contentFile.WriteString(content)
	if err != nil {
		Log.Error(LOG_ARTICLE_EDIT_ERR, "无法创建文件, 无法保存数据内容", err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_CREATE_ERR, nil)
		return
	}
	Log.Info(LOG_ARTICLE_EDIT_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, art)
}

//提交文章
func PutArticleSubmit(ctx *gin.Context) {
	var articleSubmit struct {
		ArticleId string `json:"article_id"`
		UserId    string `json:"user_id"`
	}
	//获取json中的uid和aid
	err := ctx.BindJSON(&articleSubmit)
	if err != nil {
		Log.Error(LOG_ARTICLE_SUBMIT_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判断uid与aid是否为空
	userId := articleSubmit.UserId
	uid, _ := strconv.Atoi(userId)
	articleId := articleSubmit.ArticleId
	aid, _ := strconv.Atoi(articleId)
	if len(userId) == 0 || len(articleId) == 0 {
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

//按标签获取文章
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

//按专栏获取文章
func GetArticleListByColumn(ctx *gin.Context) {
	var byColumn struct {
		ColumnId string `json:"column_id"`
	}
	err := ctx.BindJSON(&byColumn)
	if err != nil {
		Log.Error(LOG_ARTICLE_LIST_BYCOLUMN_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	columnId := byColumn.ColumnId
	//判空
	if columnId == "" {
		Log.Error(LOG_ARTICLE_LIST_BYCOLUMN_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否存在此标签
	if !sql.IsexistTagById(columnId) {
		Log.Error(LOG_ARTICLE_LIST_BYCOLUMN_ERR, RESP_INFO_DATAISNOTEXISTS)
		Resp(ctx, RESP_CODE_DATAISEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	//数据库查询
	cid, _ := strconv.Atoi(columnId)
	columns, err := sql.GetArticlesByColumnID(cid)
	if err != nil {
		Log.Error(LOG_ARTICLE_LIST_BYCOLUMN_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_LIST_BYCOLUMN_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, columns)
}

//发布文章
func PutArticleRelease(ctx *gin.Context) {
	var articleRelease struct {
		ArticleId string `json:"article_id"`
		UserId    string `json:"user_id"`
	}
	//获取json中的uid和aid
	err := ctx.BindJSON(&articleRelease)
	if err != nil {
		Log.Error(LOG_ARTICLE_RELEASE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//判断uid与aid是否为空
	userId := articleRelease.UserId
	uid, _ := strconv.Atoi(userId)
	articleId := articleRelease.ArticleId
	aid, _ := strconv.Atoi(articleId)
	if len(userId) == 0 || len(articleId) == 0 {
		Log.Error(LOG_ARTICLE_RELEASE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//判断是否为管理员
	if !sql.UserIsAdmin(userId) {
		Log.Error(LOG_ARTICLE_RELEASE_ERR, RESP_INFO_NOPERMISSION)
		Resp(ctx, RESP_CODE_NOPERMISSION, RESP_INFO_NOPERMISSION, nil)
		return
	}
	if err := sql.ArticleRelease(aid, uid); err != nil {
		Log.Error(LOG_ARTICLE_RELEASE_ERR, DB_ARTICLE_SUBMIT_ERR, err)
		Resp(ctx, RESP_CODE_SUBMIT_ERR, DB_ARTICLE_SUBMIT_ERR, nil)
		return
	}
	Log.Info(LOG_ARTICLE_RELEASE_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}

//获取未提交的文章
func GetWillBeSubmitArticleList(ctx *gin.Context) {
	var user struct {
		UserId string `json:"user_id"`
	}
	err := ctx.BindJSON(&user)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETWILLBESUB_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	uid := user.UserId
	//判空
	if len(uid) == 0 {
		Log.Error(LOG_ARTICLE_GETWILLBESUB_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否为管理员
	if !sql.UserIsAdmin(uid) {
		//普通用户只会获取自己发布的文章
		uidInt, _ := strconv.Atoi(uid)
		articles, err := sql.ArticleWillBeSubmitByUid(uidInt)
		if err != nil {
			Log.Error(LOG_ARTICLE_GETWILLBESUB_ERR, DB_SELECT_ERR, err)
			Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
			return
		}
		Log.Info(LOG_ARTICLE_GETWILLBESUB_ERR, RESP_INFO_OK)
		Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
		return
	}
	//管理员查询所有发布的文章
	articles, err := sql.ArticleWillBeSubmit()
	if err != nil {
		Log.Error(LOG_ARTICLE_GETWILLBESUB_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_GETWILLBESUB_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
}

//获取发布失败的文章
func GetFailedArticleList(ctx *gin.Context) {
	var user struct {
		UserId string `json:"user_id"`
	}
	err := ctx.BindJSON(&user)
	if err != nil {
		Log.Error(LOG_ARTICLE_RELEASEFAILD_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	uid := user.UserId
	//判空
	if len(uid) == 0 {
		Log.Error(LOG_ARTICLE_RELEASEFAILD_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否为管理员
	if !sql.UserIsAdmin(uid) {
		//普通用户只会获取自己发布的文章
		uidInt, _ := strconv.Atoi(uid)
		articles, err := sql.ArticleReleaseFailedByUid(uidInt)
		if err != nil {
			Log.Error(LOG_ARTICLE_RELEASEFAILD_ERR, DB_SELECT_ERR, err)
			Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
			return
		}
		Log.Info(LOG_ARTICLE_RELEASEFAILD_ERR, RESP_INFO_OK)
		Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
		return
	}
	//管理员查询所有发布的文章
	articles, err := sql.ArticleReleaseFailed()
	if err != nil {
		Log.Error(LOG_ARTICLE_RELEASEFAILD_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_RELEASEFAILD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
}

//获取已发布文章
func GetReleasedArticleList(ctx *gin.Context) {
	var user struct {
		UserId string `json:"user_id"`
	}
	err := ctx.BindJSON(&user)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETRELEASELIST_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	uid := user.UserId
	//判空
	if len(uid) == 0 {
		Log.Error(LOG_ARTICLE_GETRELEASELIST_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否为管理员
	if !sql.UserIsAdmin(uid) {
		//普通用户只会获取自己发布的文章
		uidInt, _ := strconv.Atoi(uid)
		articles, err := sql.ArticleReleasedByUid(uidInt)
		if err != nil {
			Log.Error(LOG_ARTICLE_GETRELEASELIST_ERR, DB_SELECT_ERR, err)
			Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
			return
		}
		Log.Info(LOG_ARTICLE_GETRELEASELIST_ERR, RESP_INFO_OK)
		Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
		return
	}
	//管理员查询所有发布的文章
	articles, err := sql.ArticleReleased()
	if err != nil {
		Log.Error(LOG_ARTICLE_GETRELEASELIST_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_GETRELEASELIST_SUCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
}

//获取将发布的文章(已提交)
func GetWillBeReleaseArticleList(ctx *gin.Context) {
	var user struct {
		UserId string `json:"user_id"`
	}
	err := ctx.BindJSON(&user)
	if err != nil {
		Log.Error(LOG_ARTICLE_GETWILLBERELEASE_ERR, RESP_INFO_GETJSONERR, err)
		Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, nil)
		return
	}
	//获取tagid
	uid := user.UserId
	//判空
	if len(uid) == 0 {
		Log.Error(LOG_ARTICLE_GETWILLBERELEASE_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//是否为管理员
	if !sql.UserIsAdmin(uid) {
		//普通用户只会获取自己提交的文章
		uidInt, _ := strconv.Atoi(uid)
		articles, err := sql.ArticleSubmitedByUid(uidInt)
		if err != nil {
			Log.Error(LOG_ARTICLE_GETWILLBERELEASE_SUCCESS, DB_SELECT_ERR, err)
			Resp(ctx, RESP_CODE_SELECT, DB_SELECT_ERR, nil)
			return
		}
		Log.Info(LOG_ARTICLE_GETWILLBERELEASE_SUCCESS, RESP_INFO_OK)
		Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
		return
	}
	//管理员获取所有已提交的文章,用来审批
	articles, err := sql.ArticleSubmited()
	if err != nil {
		Log.Error(LOG_ARTICLE_GETWILLBERELEASE_ERR, RESP_INFO_DATAISNOTEXISTS, err)
		Resp(ctx, RESP_CODE_DATAISNOTEXISTS, RESP_INFO_DATAISNOTEXISTS, nil)
		return
	}
	Log.Info(LOG_ARTICLE_GETWILLBERELEASE_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, articles)
}
