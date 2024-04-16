package server

import (
	"go-framework/config"
	"go-framework/internal/mq"
	"go-framework/internal/provider"
	"go-framework/util/mq/rocketmq"
	"go-framework/util/xlog"
	"go-framework/util/xredis"
	"go-framework/util/xsql"
	"go-framework/util/xsql/databese"
)

var Engine *SvcContext

type SvcContext struct {
	Conf        config.Conf
	DBEngine    *databese.Engine
	RedisClient *xredis.RedisClient
	Logger      *xlog.Log
	MQClient    *rocketmq.Client
	Container   *provider.Container
}

func NewSvcContext(c config.Conf, logger *xlog.Log) *SvcContext {
	redisClient := xredis.NewClient(c.Redis)
	rocketmqClient := rocketmq.NewClient(c.MQ, logger, redisClient.Default(), mq.RegisterQueue)
	rocketmqClient.ConsumerRun(mq.ConsumerHandler)

	return &SvcContext{
		Conf:        c,
		Logger:      logger,
		DBEngine:    xsql.NewClient(c.DB),
		RedisClient: xredis.NewClient(c.Redis),
		MQClient:    rocketmqClient,
		Container:   provider.Register(),
	}
}
