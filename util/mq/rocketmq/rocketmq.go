package rocketmq

import (
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/go-redis/redis/v8"
	"go-framework/util/helper"
	"go-framework/util/mq/queue"
	"go-framework/util/xlog"
)

type clientHandler func(client *Client)

type config struct {
	Endpoint  []string `json:"endpoint"`
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Namespace string   `json:"namespace"`
	Env       string   `json:"env"`
}

type Client struct {
	conf        *config
	Logger      *xlog.Log
	Producer    *Producer
	redisClient *redis.Client
	queues      map[string]queue.Queue
	Jobs        map[string]*QueueJob
	Decoder     Decoder
}

func NewClient(c interface{}, logger *xlog.Log, redisClient *redis.Client, fs ...clientHandler) (client *Client) {
	var conf *config
	err := helper.UnMarshalWithInterface(c, &conf)
	if err != nil {
		logger.Panicf("rocketmq config error: %v", err)
	}
	client = &Client{
		conf:        conf,
		Logger:      logger,
		redisClient: redisClient,
		queues:      make(map[string]queue.Queue),
		Jobs:        make(map[string]*QueueJob),
	}

	for _, f := range fs {
		f(client)
	}

	client.Producer = NewProducer(client)

	err = client.RegisterJob()
	if err != nil {
		client.Logger.Panicf("register job error: %v", err)
	}

	client.Decoder = NewJobDecoder(client)
	return
}

// ConsumerRun 启动消费者
func (c *Client) ConsumerRun(handler func(client *Client)) {
	handler(c)
}

// Client 创建阿里云客户端
func (c *Client) Client() mq_http_sdk.MQClient {
	client := mq_http_sdk.NewAliyunMQClient(c.conf.Endpoint[0], c.conf.AccessKey, c.conf.SecretKey, "")
	return client
}
