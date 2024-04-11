package rocketmq

import (
	"context"
	"fmt"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
	"github.com/panjf2000/ants/v2"
	"go-framework/util/helper"
	"go-framework/util/locker"
	"go-framework/util/mq/queue"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	RetryTimeKeyFormat = "retry-time:%s:%s"
	RockTimeKeyFormat  = "rock:%s:%s"
)

var NotAskBuffer []mq_http_sdk.ConsumeMessageEntry

type Consumer struct {
	client           *Client
	consumer         mq_http_sdk.MQConsumer
	queue            queue.Queue
	askBuffer        []mq_http_sdk.ConsumeMessageEntry
	pool             *ants.Pool
	askBufferLock    sync.Mutex
	batchAskInterval time.Duration
	concurrency      int
	retryTimes       int64
}

type ConsumerOption func(consumer *Consumer)

func WithConcurrency(concurrency int) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.concurrency = concurrency
	}
}

func WithRetryTimes(retryTimes int64) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.retryTimes = retryTimes
	}
}

func ConsumerMessage(client *Client, queue queue.Queue, opts ...ConsumerOption) {
	consumer := &Consumer{
		client:           client,
		queue:            queue,
		batchAskInterval: time.Millisecond * 200,
	}
	for _, opt := range opts {
		opt(consumer)
	}

	consumer.pool, _ = ants.NewPool(consumer.concurrency)
	consumer.consumer = client.Client().GetConsumer(client.conf.Namespace, queue.Topic(), queue.GroupId(), "")

	respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
	errChan := make(chan error)

	go consumer.pullMessage(respChan, errChan)

	consumer.StartProcessing(respChan, errChan)
	return
}

// StartProcessing 启动消息处理循环
func (c *Consumer) StartProcessing(respChan chan mq_http_sdk.ConsumeMessageResponse, errChan chan error) {
	var wg sync.WaitGroup

	go c.batchAskTimer()       // Start the batch ask timer in a separate goroutine
	go c.errorHandler(errChan) // Start the error handler in a separate goroutine

	for resp := range respChan {
		for _, msg := range resp.Messages {
			wg.Add(1)
			msgCopy := msg
			_ = c.pool.Submit(func() {
				defer wg.Done()
				c.processMessage(msgCopy)
			})
		}
	}

	wg.Wait()
	c.sendBatchAsk() // Send remaining asks if any
}

// pullMessage 获取消息
func (c *Consumer) pullMessage(respChan chan mq_http_sdk.ConsumeMessageResponse, errChan chan error) {
	for {
		c.consumer.ConsumeMessage(respChan, errChan, 16, 30)
		time.Sleep(time.Millisecond * 500)
	}

	//job := &job.OrderJob{}
	//_, msgd, _ := c.client.Decoder.Marshal(job, []byte("hello world 000"))
	//for {
	//	notAskBuffer := NotAskBuffer
	//	if len(notAskBuffer) == 0 {
	//		respChan <- mq_http_sdk.ConsumeMessageResponse{Messages: generateMessages(1, msgd)}
	//	} else {
	//		respChan <- mq_http_sdk.ConsumeMessageResponse{Messages: notAskBuffer}
	//	}
	//
	//	time.Sleep(15 * time.Second)
	//}
}

//func generateMessages(count int, msgd []byte) []mq_http_sdk.ConsumeMessageEntry {
//	var messages []mq_http_sdk.ConsumeMessageEntry
//	properties := make(map[string]string)
//	properties["groupId"] = "GID_test_dev38"
//
//	for i := 0; i < count; i++ {
//
//		message := mq_http_sdk.ConsumeMessageEntry{}
//		message.MessageId = xid.New().String()
//		message.ReceiptHandle = message.MessageId
//		message.Properties = properties
//		message.MessageBody = string(msgd)
//		messages = append(messages, message)
//	}
//	return messages
//}

// processMessage 处理消息
func (c *Consumer) processMessage(message mq_http_sdk.ConsumeMessageEntry) {
	defer helper.RecoverPanic(c.client.Logger)

	groupName := c.client.GetGroupName(c.queue.Topic())
	if message.Properties["groupId"] != groupName {
		return
	}

	rockKey := c.GetRockKey(c.queue.Topic(), message.MessageId)
	retryTimesKey := c.GetRetryTimesKey(c.queue.Topic(), message.MessageId)

	// 为300秒，超时会导致重复消费，Http协议，该时间不支持配置修改，
	// 放置最前面，放置其他调用类等情况的致命异常
	mutex := locker.NewMutex(c.client.redisLock)
	err := mutex.Lock(rockKey, time.Second*600)
	if err != nil {
		c.client.Logger.Errorf("key: %s 消费id：%s，消息重复消费: %+v", rockKey, message.MessageId, err)
		return
	}
	defer mutex.UnLock()

	times := c.client.redisClient.Incr(context.Background(), retryTimesKey).Val()
	c.client.redisClient.Expire(context.Background(), retryTimesKey, time.Second*600)

	fmt.Println(message.MessageId, "应答次数， retryTimes", times)
	isAsk := true
	if times > c.retryTimes+1 {
		c.addAskBuffer(message, &isAsk)
	}

	msgBody := message.MessageBody

	if c.client.Decoder.Check(msgBody) {
		task, msgBodyByte, err := c.client.Decoder.UnMarshal(msgBody)
		if err != nil {
			c.client.Logger.Errorf("消息反序列化失败: %+v", err)
		}

		if c.retryTimes == 0 {
			fmt.Println(message.MessageId, "直接应答， retryTimes", c.retryTimes)
			c.addAskBuffer(message, &isAsk)
		}

		// 捕获panic
		var isError bool
		c.taskExecute(task, msgBodyByte, &isError)
		fmt.Println(message.MessageId, "出现报错， isError", isError)
		if isError {
			fmt.Println(message.MessageId, "出现报错， isAsk", isAsk)
			isAsk = false
		}
	}

	c.addAskBuffer(message, &isAsk)
	fmt.Println(message.MessageId, "准备结束， isAsk", isAsk)
}

func (c *Consumer) taskExecute(task queue.Job, msgBodyByte []byte, isError *bool) {
	defer func() {
		if err := recover(); err != nil {
			*isError = true
			c.client.Logger.Errorf("消息执行失败: %+v", err)

			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			c.client.Logger.Errorf("%s", buf[:n])
		}
	}()
	err := task.Execute(msgBodyByte)
	if err != nil {
		*isError = true
		c.client.Logger.Errorf("%s 消息消费失败,参数：%s, 执行失败: %+v", task.Name(), string(msgBodyByte), err)
	}
}

func (c *Consumer) addAskBuffer(message mq_http_sdk.ConsumeMessageEntry, isAsk *bool) {
	if !*isAsk {
		NotAskBuffer = append(NotAskBuffer, message)
		fmt.Println(message.MessageId, "不进行 isAsk")
		return
	}
	defer c.askBufferLock.Unlock()
	c.askBufferLock.Lock()
	c.askBuffer = append(c.askBuffer, message)
	fmt.Println(message.MessageId, "进行 isAsk")
	*isAsk = false
}

func (c *Consumer) batchAskTimer() {
	timer := time.NewTicker(c.batchAskInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.sendBatchAsk()
		}
	}
}

func (c *Consumer) sendBatchAsk() {
	c.askBufferLock.Lock()
	defer c.askBufferLock.Unlock()

	askBufferLen := len(c.askBuffer)
	if askBufferLen > 16 {
		askBufferLen = 16
	}

	if askBufferLen == 0 {
		return
	}

	newAskBuffer := c.askBuffer[:askBufferLen] // Clear the buffer

	var messageId []string
	var receiptHandle []string
	for _, message := range newAskBuffer {
		messageId = append(messageId, message.MessageId)
		receiptHandle = append(receiptHandle, message.ReceiptHandle)
	}
	fmt.Println("出现应答：", receiptHandle)
	//ackErr := c.consumer.AckMessage(receiptHandle)
	//if ackErr != nil {
	//	// 某些消息的句柄可能超时，会导致消息消费状态确认不成功。
	//	if errAckItems, ok := ackErr.(errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem); ok {
	//		c.client.Logger.Errorf("消息id: %+v,确认消费失败: %+v", messageId, errAckItems)
	//
	//	} else {
	//		c.client.Logger.Errorf("消息id: %+v,确认消费失败: %+v", messageId, ackErr)
	//	}
	//	time.Sleep(time.Second * 3)
	//	return
	//}

	c.client.Logger.Infof("消息成功: %+v", messageId)
	c.askBuffer = c.askBuffer[:askBufferLen]
}

func (c *Consumer) errorHandler(errChan chan error) {
	for err := range errChan {
		// 消费出现异常
		if !strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
			c.client.Logger.Errorf("消费消息失败: %+v", err)
		}
	}
}

func (c *Consumer) GetRetryTimesKey(topic string, message string) string {
	return fmt.Sprintf(RetryTimeKeyFormat, topic, message)
}

func (c *Consumer) GetRockKey(topic string, message string) string {
	return fmt.Sprintf(RockTimeKeyFormat, topic, message)
}
