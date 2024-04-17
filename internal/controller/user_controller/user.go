package user_controller

import (
	"fmt"
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

func GetUserInfo2(svc *server.SvcContext) app.HandlerFunc {
	return func(c *app.Context) {
		fmt.Println(c.Param("name"))
		type student struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var stu student
		stu.Name = "fuzengyao"
		stu.Age = 18
		c.JSON(200, stu)
	}
}
