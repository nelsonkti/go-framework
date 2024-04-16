package user_controller

import (
	"github.com/gin-gonic/gin"
	"go-framework/internal/server"
	"go-framework/util/app"
	"net/http"
)

func GetUserInfo(ctx *server.SvcContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx.Container.OrderService.Get()
		c.String(http.StatusOK, "Welcome Gin Server")
	}
}

func GetUserInfo2(ctx *server.SvcContext) app.HandlerFunc {
	return func(c *app.Context) {
		c.ResponseWriter.Write([]byte("Hello, ddd!"))
	}
}
