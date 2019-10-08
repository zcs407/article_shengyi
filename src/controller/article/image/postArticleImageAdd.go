package image

import (
	. "articlebk/src/common"
	"articlebk/src/common/sql"
	"github.com/gin-gonic/gin"
	"io/ioutil"

	"math/rand"
	"path"
	"strings"
	"time"
)

func PostArticleImageAdd(ctx *gin.Context) {

	aid := ctx.Query("aid")

	rand.Seed(time.Now().UnixNano())
	//定义返回信息的对象
	var resp struct {
		Code       string   `json:"code"`
		Info       string   `json:"info"`
		ImageUrls  []string `json:"image_urls"`
		ImagePaths []string `json:"image_path"`
	}
	//从post的json中获取图片文件数组
	imageForm, _ := ctx.MultipartForm()
	images := imageForm.File["images[]"]
	imageUrls := make([]string, 0)
	//图片切片不可为空,至少要有一个图片
	if len(images) == 0 {
		resp.Code = "404"
		resp.ImageUrls = imageUrls
		resp.Info = "没有获取到post的图片"
		ctx.JSON(200, resp)
		Log.Error("[article_add_image]", resp)
		return
	}
	//每张图片不能大于10兆
	for _, image := range images {
		if image.Size > 10000000 {
			resp.Code = "300"
			resp.ImageUrls = imageUrls
			resp.Info = "图片不能超过10兆"
			ctx.JSON(200, resp)
			Log.Error("[article_add_image]", resp)
			return
		}
		//判断图片类型是否为jpg,png或jpeg格式
		ext := path.Ext(image.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			resp.Code = "300"
			resp.ImageUrls = imageUrls
			resp.Info = "图片格式不支持,目前支持的格式为:jpg,png或jpeg"
			ctx.JSON(200, resp)
			Log.Error("[article_add_image]", resp)
			return
		}
		imagefile, err := ioutil.TempFile("./articleData/images/", "article-img-*"+ext)
		//定义外网访问的URL

		imagePath := "./" + imagefile.Name()
		imgsplit := strings.Split(imagefile.Name(), "/")
		imgname := string(imgsplit[len(imgsplit)-1])
		imageUrl := "http://127.0.0.1:8888/images/" + imgname
		Log.Info("[article_add_image]", "数据库应该添加的图片路径有:", imageUrl)
		//插入数据库,如果失败则不保存文件
		err = sql.ArticleImageAdd(imageUrl, imagePath, aid)
		if err != nil {
			resp.Code = "501"
			resp.ImageUrls = imageUrls
			resp.Info = "数据库保存图片错误"
			ctx.JSON(200, resp)
			Log.Error("[article_add_image]", resp)
			return
		}
		err = ctx.SaveUploadedFile(image, imagefile.Name())
		if err != nil {
			resp.Code = "300"
			resp.ImageUrls = imageUrls
			resp.Info = "图片已存在"
			ctx.JSON(200, resp)
			Log.Warn("[article_add_image]", resp)
			return
		}
		imageUrls = append(imageUrls, imageUrl)
	}

	resp.Code = "200"
	resp.ImageUrls = imageUrls
	resp.Info = "添加图片成功"
	Log.Info("[article_add_image]", resp)
	ctx.JSON(200, resp)
}
