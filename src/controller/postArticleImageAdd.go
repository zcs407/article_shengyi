package controller

import (
	. "articlebk/src/common"
	"articlebk/src/common/database/sql"
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
	//从post的json中获取图片文件数组
	imageForm, _ := ctx.MultipartForm()
	images := imageForm.File["images[]"]
	imageUrls := make([]string, 0)
	//图片切片不可为空,至少要有一个图片
	if len(images) == 0 {
		Log.Error(LOG_IMAGE_ADD_ERR, RESP_INFO_JSON_DATANULL)
		Resp(ctx, RESP_CODE_JSON_DATANULL, RESP_INFO_JSON_DATANULL, nil)
		return
	}
	//每张图片不能大于10兆
	for _, image := range images {
		if image.Size > 10000000 {
			Log.Error(LOG_IMAGE_ADD_ERR, LOG_IAMGE_SIZE_ERR)
			Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, LOG_IAMGE_SIZE_ERR)
			return
		}
		//判断图片类型是否为jpg,png或jpeg格式
		ext := path.Ext(image.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			Log.Error(RESP_INFO_GETJSONERR, LOG_IAMGE_EXT_ERR)
			Resp(ctx, RESP_CODE_GETJSONERR, RESP_INFO_GETJSONERR, LOG_IAMGE_EXT_ERR)
			return
		}
		imageFile, err := ioutil.TempFile("./articleData/images/", "article-img-*"+ext)
		//定义外网访问的URL
		imagePath := "./" + imageFile.Name()
		imgSplit := strings.Split(imageFile.Name(), "/")
		imgName := string(imgSplit[len(imgSplit)-1])
		imageUrl := "http://127.0.0.1:8888/images/" + imgName
		//插入数据库,如果失败则不保存文件
		err = sql.ArticleImageAdd(imageUrl, imagePath, aid)
		if err != nil {
			Log.Error(LOG_IMAGE_ADD_ERR, DB_CREATE_ERR, err)
			Resp(ctx, RESP_CODE_CERATE_ERR, RESP_INFO_CREATE_ERR, nil)
			return
		}
		err = ctx.SaveUploadedFile(image, imageFile.Name())
		if err != nil {
			Log.Error(LOG_IMAGE_ADD_ERR, DB_CREATE_ERR, err)
			Resp(ctx, RESP_CODE_CERATE_ERR, RESP_INFO_CREATE_ERR, nil)
			return
		}
		imageUrls = append(imageUrls, imageUrl)
	}
	Log.Info(LOG_IMAGE_ADD_SUCCESS, RESP_INFO_OK)
	Resp(ctx, RESP_CODE_OK, RESP_INFO_OK, nil)
}
