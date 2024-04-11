package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
	"go-framework/config"
	"go-framework/internal/mq"
	"go-framework/internal/router"
	"go-framework/util/mq/rocketmq"
	"go-framework/util/xconfig"
	"go-framework/util/xconfig/file"
	"go-framework/util/xlog"
	"go-framework/util/xredis"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type logicProgram struct {
	once sync.Once
}

func main() {
	p := &logicProgram{}
	if err := svc.Run(p, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL); err != nil {
		fmt.Println(err)
	}

}

// svc 服务运行框架 程序启动时执行Init+Start, 服务终止时执行Stop
func (p *logicProgram) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *logicProgram) Start() error {
	var c config.Conf
	path := "/Users/fzy/workspace/Go/src/rrzuji/go-framework"
	err := xconfig.New(&c, file.NewConfig(path+"/config.yaml"))
	if err != nil {
		panic(err)
	}

	logger, err := xlog.NewLogger(path + "/log/test.log")
	if err != nil {
		panic(err)
	}

	redisClient, err := xredis.NewClient(c.Redis)

	rocketmqClient := rocketmq.NewClient(c.MQ, logger, redisClient["default"], mq.RegisterQueue)
	rocketmqClient.ConsumerRun(mq.ConsumerHandler)

	r := gin.Default()
	router.Register(r)
	r.Run()
	return nil
}

func (p *logicProgram) Stop() error {
	p.once.Do(func() {
	})
	return nil
}
