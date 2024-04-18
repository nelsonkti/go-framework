package front

import (
	"github.com/gin-gonic/gin"
	"go-framework/internal/controller/user_controller"
	"go-framework/internal/server"
)

func FrontRegister(app *gin.RouterGroup, ctx *server.SvcContext) {
	app.GET("/kkk", user_controller.GetUserInfo(ctx))
}