package mq

import (
	"go-framework/internal/mq/queues"
	"go-framework/util/mq/rocketmq"
)

func ConsumerHandler(client *rocketmq.Client) {
	rocketmq.ConsumerMessage(client, &queues.OrderQueue{}, rocketmq.WithConcurrency(1), rocketmq.WithRetryTimes(2)) // 订单队列消费者
	//rocketmq.ConsumerMessage(client, &queues.ShopQueue{}, rocketmq.WithConcurrency(3))                               // 商家队列消费者
}
