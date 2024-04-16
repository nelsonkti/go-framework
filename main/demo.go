package main

import (
	"fmt"
	"github.com/judwhite/go-svc"
	"go-framework/config"
	"go-framework/internal/router"
	"go-framework/internal/server"
	"go-framework/util/app"
	"go-framework/util/xconfig"
	"go-framework/util/xconfig/file"
	"go-framework/util/xlog"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type demoProgram struct {
	once sync.Once
}

func main() {
	p := &demoProgram{}
	if err := svc.Run(p, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL); err != nil {
		fmt.Println(err)
	}

}

// svc 服务运行框架 程序启动时执行Init+Start, 服务终止时执行Stop
func (p *demoProgram) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *demoProgram) Start() error {
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

	svc := server.NewSvcContext(c, logger)

	a := app.New()
	router.Register2(a, svc)

	a.Run(":8081")

	if err != nil {
		panic(err)
	}
	return nil
}

func (p *demoProgram) Stop() error {
	p.once.Do(func() {
	})
	return nil
}
