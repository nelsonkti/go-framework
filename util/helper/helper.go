package helper

import (
	"fmt"
	"go-framework/util/xlog"
	"runtime"
)

// RecoverPanic 恢复panic
func RecoverPanic(logger *xlog.Log) {
	err := recover()
	if err != nil {
		fmt.Println("进来了")
		logger.Error(err)

		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		logger.Errorf("%s", buf[:n])
	}
}
