package user_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-framework/internal/server"
	"go-framework/util/app"
	"go-framework/util/xerror"
	"net/http"
)

func GetUserInfo(ctx *server.SvcContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ctx.Container.OrderService.Get()
		if err != nil {
			fmt.Println("Error tips")
			if xerror.IsStatus(err, 3002) {
				fmt.Println("Error 3002")
			}
			err1 := xerror.UnmarshalError(err)
			fmt.Println(err1)
		}
		type student struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var stu student
		stu.Name = "fuzengyao"
		stu.Age = 18
		c.JSON(http.StatusOK, stu)
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
