package server

import (
	"github.com/go-redis/redis/v8"
	"go-framework/config"
	"go-framework/util/locker"
	"go-framework/util/mq/rocketmq"
	"go-framework/util/xlog"
	"go-framework/util/xsql/databese"
)

var Engine Server

type Server struct {
	Conf        config.Conf
	DBEngine    *databese.Engine
	RedisEngine map[string]*redis.Client
	Logger      *xlog.Log
	RedisLock   *locker.RedisLock
	MQClient    *rocketmq.Client
}
