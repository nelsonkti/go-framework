package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/judwhite/go-svc"
	"go-framework/config"
	"go-framework/internal/router"
	"go-framework/internal/server"
	validator2 "go-framework/util/validator"
	"go-framework/util/xconfig"
	"go-framework/util/xconfig/file"
	"go-framework/util/xlog"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type logicProgram struct {
	once       sync.Once
	svcContext *server.SvcContext
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

	p.svcContext = server.NewSvcContext(c, logger)

	go func() {
		newApp(c, p.svcContext)
	}()

	return nil
}

func newApp(c config.Conf, s *server.SvcContext) {
	// 创建并配置验证器
	r := gin.Default()
	router.Register(r, s)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("chinese", validator2.ChineseValidation)
	}
	err := r.Run(c.Server.Http.Addr)
	if err != nil {
		panic(err)
	}
}

func (p *logicProgram) Stop() error {
	p.once.Do(func() {
		defer p.svcContext.DBEngine.Close()
	})
	return nil
}
