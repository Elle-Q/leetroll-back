package router

import (
	"github.com/gin-gonic/gin"
	"leetroll/admin/api"
	"leetroll/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCheckRouter)
}

func registerCheckRouter(g *gin.RouterGroup) {

	_UserApi := api.User{}
	user := g.Group("/user").Use(middleware.AuthJWTMiddleware())
	{
		user.GET("/:id", _UserApi.GetUser)
		user.GET("/list", _UserApi.List)
		//r.POST("", api.Insert)
		//r.PUT("", api.Update)
		//r.DELETE("", api.Delete)
	}

	_CatApi := api.Category{}
	cat := g.Group("/cat").Use(middleware.AuthJWTMiddleware())
	{
		cat.GET("/:id", _UserApi.GetUser)
		cat.GET("/list", _CatApi.List)
		cat.GET("/list-name", _CatApi.ListName)
		cat.POST("/update", _CatApi.Update)
		cat.POST("/delete", _CatApi.Delete)
	}

	_ItemApi := api.Item{}
	item := g.Group("/item").Use(middleware.AuthJWTMiddleware())
	{
		item.GET("/list", _ItemApi.List)
		item.POST("/update", _ItemApi.Update)
		item.POST("/upload", _ItemApi.Upload)
		item.POST("/chapter/upload", _ItemApi.Upload)
		item.POST("/delete", _ItemApi.Delete)
		item.GET("/:ID", _ItemApi.Get)
		item.GET("/files/:ID", _ItemApi.GetFilesByItemId)
	}

	_FileApi := api.File{}
	file := g.Group("/file").Use(middleware.AuthJWTMiddleware())
	{
		file.POST("/delete", _FileApi.DeleteFile)
	}

	_ChapterApi := api.Chapter{}
	chapter := g.Group("/chapter").Use(middleware.AuthJWTMiddleware())
	{
		chapter.POST("/upload", _ChapterApi.Upload)
		chapter.POST("/file/delete", _ChapterApi.FileDelete)

	}
}
