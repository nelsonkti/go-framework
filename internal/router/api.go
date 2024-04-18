package router

import (
	"github.com/gin-gonic/gin"
	"go-framework/internal/controller/user_controller"
	"go-framework/internal/router/admin"
	"go-framework/internal/router/front"
	"go-framework/internal/server"
)

type StudentRequest struct {
	Name string `json:"name" binding:"required,chinese" label:"文章ID"`
}

func Register(app *gin.Engine, ctx *server.SvcContext) {
	app.GET("/kkk", user_controller.GetUserInfo(ctx))

	api := app.Group("/api")
	admin.AdminRegister(api, ctx) // 后台路由
	front.FrontRegister(api, ctx) // 后台路由

	//api.GET("/hello", func(c *gin.Context) {
	//	var studentRequest StudentRequest
	//	if err := c.ShouldBind(&studentRequest); err != nil {
	//		fmt.Println(err)
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//		return
	//	}
	//
	//	c.JSON(200, gin.H{
	//		"message": "Hello World",
	//	})
	//})
}
