package router

import (
	"articlebk/src/controller"
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
		r1.POST("/login", controller.PostLogin)
		r1.POST("/register", controller.PostRegister)
	}
	r2 := router.Group("/user")
	{
		r2.DELETE("/userDelete", controller.DeleteUser)
		r2.PUT("/userUpdate/userPwdUpdate", controller.PutUserPwdUpdate)
		r2.PUT("/userUpdate/userRollerUpdate", controller.PutUserRollerUpdate)
		r2.GET("/userList", controller.GetUserList)

	}
	r3 := router.Group("/article")
	{
		r3.POST("/articleImageAdd", controller.PostArticleImageAdd)
		r3.POST("/articleAdd", controller.PostArticleAdd)
		r3.PATCH("/articleEdit", controller.PatchArticleEdit)
		r3.PUT("/articleSubmit", controller.PutArticleSubmit)
		r3.PUT("/articleRelease", controller.PutArticleRelease)
		r3.DELETE("/articleDelete", controller.DeleteArticleDel)
		r3.GET("/articleListByColumn", controller.GetArticleListByColumn)
		r3.GET("/articleListByTag", controller.GetArticleListByTag)
		r3.GET("/getFailedArticleList", controller.GetFailedArticleList)
		r3.GET("/getWillBeReleaseArticleList", controller.GetWillBeReleaseArticleList)
		r3.GET("/getReleasedArticleList", controller.GetReleasedArticleList)
	}
	r4 := router.Group("/column")
	{
		r4.POST("/columnAdd", controller.PostColumnAdd)
		r4.PUT("/columnCname", controller.PutColumnCname)
		r4.DELETE("/columnDel", controller.DeleteColumn)
		r4.GET("/columnList", controller.GetColumnList)
	}
	r5 := router.Group("/tag")
	{
		r5.POST("/tagAdd", controller.PostTagAdd)
		r5.POST("/tagDel", controller.PostTagDel)
		r5.POST("/TagCname", controller.PostTagCname)
	}
	return router
}
