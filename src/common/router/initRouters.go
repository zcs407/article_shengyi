package router

import (
	"articlebk/src/controller/article"
	"articlebk/src/controller/article/image"
	"articlebk/src/controller/column"
	"articlebk/src/controller/tag"
	"articlebk/src/controller/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24,
	})
	router.Use(sessions.Sessions("name", store))
	r1 := router.Group("/")
	{
		r1.POST("/login", user.PostLogin)
		r1.POST("/register", user.PostRegister)
	}
	r2 := router.Group("/user")
	{
		r2.DELETE("/userDelete", user.DeleteUser)
		r2.PUT("/userUpdate/userPwdUpdate", user.PutUserPwdUpdate)
		r2.PUT("/userUpdate/userRollerUpdate", user.PutUserRollerUpdate)
		r2.GET("/userList", user.GetUserList)

	}
	r3 := router.Group("/article")
	{
		r3.POST("/articleImageAdd", image.PostArticleImageAdd)
		r3.POST("/articleAdd", article.PostArticleAdd)
		r3.PATCH("/articleEdit", article.PatchArticleEdit)
		r3.PUT("/articleSubmit", article.PutArticleSubmit)
		r3.PUT("/articleRelease", article.PutArticleRelease)
		r3.DELETE("/articleDelete", article.DeleteArticleDel)
		r3.GET("/articleListBySpecial", article.GetArticlelistBySpecial)
		r3.GET("/articleListByTag", article.GetArticleListByTag)
		r3.GET("/getFailedArticleList", article.GetFailedArticleList)
		r3.GET("/getWillBeReleaseArticleList", article.GetWillBeReleaseArticleList)
		r3.GET("/getReleasedArticleList", article.GetReleasedArticleList)
	}
	r4 := router.Group("/column")
	{
		r4.POST("/columnAdd", column.PostColumnAdd)
		r4.PUT("/columnCname", column.PutColumncname)
		r4.DELETE("/columnDel", column.DeleteColumn)
		r4.GET("/columnList", column.GetColumnList)
	}
	r5 := router.Group("/tag")
	{
		r5.POST("/tagAdd", tag.PostTagAdd)
		r5.POST("/tagDel", tag.PostTagDel)
		r5.POST("/TagCname", tag.PostTagCname)
	}
	return router
}
