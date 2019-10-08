package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetArticlelistBySpecial(ctx *gin.Context) {
	sid := ctx.Query("sid")
	fmt.Println("sid", sid)
	//dbtable.ArticleSelectBySpecial(sid)
	//TODO

}
