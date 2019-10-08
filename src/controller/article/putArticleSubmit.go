package article

import (
	"github.com/gin-gonic/gin"
	"log"
)

func PutArticleSubmit(ctx *gin.Context) {
	loger := log.New(gin.DefaultWriter, "", log.LstdFlags|log.Lshortfile)
	resp := make(map[string]string)
	var articleSubmit struct {
		Aid string `json:"aid"`
		Uid string `json:"uid"`
	}
	err := ctx.BindJSON(articleSubmit)
	if err != nil {
		resp["404"] = "无法获取到post中的aid和uid"
		loger.Println("[article_submit]", resp)
		ctx.JSON(200, resp)
		return
	}
	if articleSubmit.Aid == "" || articleSubmit.Uid == "" {
		resp["405"] = "aid与uid不能为空"
		loger.Println("[article_submit]", resp)
		ctx.JSON(200, resp)
		return
	}

}
