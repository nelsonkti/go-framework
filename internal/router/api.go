package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-framework/internal/controller/user_controller"
	"go-framework/internal/mq/job"
	"go-framework/internal/server"
	"go-framework/util/app"
	"net/http"
)

func Register(app *gin.Engine, ctx *server.SvcContext) {
	app.GET("/kkk", user_controller.GetUserInfo(ctx))
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	api := app.Group("/api")
	api.GET("/hello", func(req *gin.Context) {

		msg := fmt.Sprintf("Hello world-%d", 66)
		err := server.Engine.MQClient.Producer.SendJobMessage(&job.OrderJob{}, []byte(msg))
		fmt.Println(err)
		//msg = fmt.Sprintf("Hello world-%d", 999)
		//server.Engine.MQClient.Producer.SendJobMessage(&job.ShopJob{}, []byte(msg))
		//for i := 0; i < 1; i++ {
		//	msg := fmt.Sprintf("Hello world-%d", 66)
		//	server.Engine.MQClient.Producer.SendJobMessage(&job.OrderJob{}, []byte(msg))
		//	//msg = fmt.Sprintf("Hello world-%d", 999)
		//	//server.Engine.MQClient.Producer.SendJobMessage(&job.ShopJob{}, []byte(msg))
		//
		//	//err = client.Producer.SendMessage("test", []byte(fmt.Sprintf("test: %s", time.Now().Format(time.DateTime))))
		//}
		fmt.Println("2313")
		req.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
}

func Register2(app *app.App, ctx *server.SvcContext) {
	app.GET("/kkk2", user_controller.GetUserInfo2(ctx))
}
