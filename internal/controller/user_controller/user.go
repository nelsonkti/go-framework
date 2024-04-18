package user_controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-framework/internal/server"
	"go-framework/util/app"
	"go-framework/util/xerror"
	"go-framework/util/xhttp"
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

		}

		err4 := xerror.IsError(err)
		fmt.Println("err4", err4)

		err3 := errors.New("错误测试")
		err9 := xerror.IsError(err3)
		fmt.Println("err9", err9)
		type student struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var stu student
		stu.Name = "fuzengyao"
		stu.Age = 18

		c.JSON(http.StatusOK, xhttp.Error(err))
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
