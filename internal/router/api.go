package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-framework/internal/mq/job"
	"go-framework/internal/server"
	"net/http"
	"time"
)

func Register(app *gin.Engine) {
	app.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	api := app.Group("/api")
	api.GET("/hello", func(req *gin.Context) {

		msg := fmt.Sprintf("Hello world-%d", 66)
		server.Engine.MQClient.Producer.SendJobMessage(&job.OrderJob{}, []byte(msg))
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
