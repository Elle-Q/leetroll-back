package router

import (
	"github.com/gin-gonic/gin"
	"gomo/app/api"
	"gomo/common/actions"
)

func init()  {
	routerCheckRole = append(routerCheckRole, registerCheckRouter)
}

func registerCheckRouter(g *gin.RouterGroup) {

	_UserApi := api.User{}
	user := g.Group("/user").Use(actions.PermissionAction())
	{
		user.GET("/:id", _UserApi.GetUser)
		user.GET("/register", _UserApi.AddUser)
		//r.GET("/:id", api.Get)
		//r.POST("", api.Insert)
		//r.PUT("", api.Update)
		//r.DELETE("", api.Delete)
	}


}