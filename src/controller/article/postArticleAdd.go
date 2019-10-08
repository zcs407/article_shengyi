package article

import (
	"articlebk/src/common/dbtable"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func PostArticleAdd(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	var article struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		SpecialId string   `json:"special_id"`
		Tags      []string `json:"tags"`
	}
	var resp struct {
		Title      string `json:"title"`
		Aid        string `json:"aid"`
		ContentUrl string `json:"content_url"`
		Code       string `json:"code"`
		Info       string `json:"info"`
	}
	//获取post的文章信息
	err := ctx.BindJSON(&article)
	if err != nil {
		resp.Title = ""
		resp.ContentUrl = ""
		resp.Code = "404"
		resp.Info = "无法获取到post文章信息"
		loger.Println("[article_add]", resp, err)
		ctx.JSON(200, resp)
		return
	}
	//获取json中文章的主题/内容/专题id/标签的id数组(一个文章属于多个标签)
	content := article.Content
	title := article.Title
	specialId := article.SpecialId
	sid, _ := strconv.Atoi(specialId)
	tags := article.Tags
	//创建本地文件用于保存文章
	contentfile, err := ioutil.TempFile("./articleData/articleFile/", "article-*.txt")
	contentPath := contentfile.Name()
	defer contentfile.Close()
	//创建Url提供用户访问
	imagesplit := strings.Split(contentfile.Name(), "/")
	filename := string(imagesplit[len(imagesplit)-1])
	contentUrl := "http://127.0.0.1:8888/articleFile/" + filename
	//给article表结构赋值
	article_table := dbtable.Article{
		Title:       title,
		ContentUrl:  contentUrl,
		ContentPath: "./" + contentPath,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Status:      0,
		SpecialId:   sid,
	}

	//判空
	if article.Title == "" || content == "" || len(tags) == 0 || specialId == "" {
		resp.Title = ""
		resp.ContentUrl = ""
		resp.Code = "404"
		resp.Info = "文章标题或内容为空,请填写"
		loger.Println("[article_add]", resp)
		ctx.JSON(200, resp)
		return
	}
	if sql.IsexistArticle(article.Title) {
		resp.Title = article.Title
		resp.ContentUrl = ""
		resp.Code = "404"
		resp.Info = "文章已存在"
		loger.Println("[article_add]", resp)
		ctx.JSON(200, resp)
		return
	}
	aid, err := sql.ArticleAdd(article_table, tags)
	if err != nil {
		resp.Title = article.Title
		resp.ContentUrl = contentUrl
		resp.Code = "404"
		resp.Info = "无法获取到post文章信息"
		loger.Println("[article_add]", resp)
		ctx.JSON(200, resp)
		return
	}
	_, err = contentfile.WriteString(content)
	if err != nil {
		loger.Println("[article_add]", "无法创建文件, 无法保存数据内容", err)
		return
	}
	resp.Title = article.Title
	resp.Aid = aid
	resp.ContentUrl = contentUrl
	resp.Code = "200"
	resp.Info = "文章创建成功"
	loger.Println("[article_add]", resp)
	ctx.JSON(200, resp)
}
