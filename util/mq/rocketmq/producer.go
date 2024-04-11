package rocketmq

import (
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"go-framework/util/mq/queue"
	"time"
)

type Producer struct {
	client *Client
}

func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// SendJobMessage 发送任务消息
func (p *Producer) SendJobMessage(job queue.Job, msg []byte) error {
	topic, msgd, err := p.client.Decoder.Marshal(job, msg)
	if err != nil {
		p.client.Logger.Errorf("marshal job error, %v", err)
		return err
	}

	return p.SendMessage(topic, msgd)
}

// SendJobDelayMessage 发送延时任务消息
func (p *Producer) SendJobDelayMessage(job queue.Job, msg []byte, duration time.Duration) error {
	topic, msgd, err := p.client.Decoder.Marshal(job, msg)
	if err != nil {
		p.client.Logger.Errorf("marshal job error, %v", err)
		return err
	}

	return p.SendDelayMessage(topic, msgd, duration)
}

// SendMessage 发送消息
func (p *Producer) SendMessage(topic string, msg []byte) error {
	msgRequest := p.publishMessageRequest(topic, msg)

	return p.publishMessage(msgRequest, topic)
}

// SendDelayMessage 发送延时消息
func (p *Producer) SendDelayMessage(topic string, msg []byte, duration time.Duration) error {
	msgRequest := p.publishMessageRequest(topic, msg)
	msgRequest.StartDeliverTime = time.Now().Add(duration).Unix() * 1000

	err := p.publishMessage(msgRequest, topic)
	return err
}

// publishMessage 发送消息
func (p *Producer) publishMessage(msgRequest mq_http_sdk.PublishMessageRequest, topic string) error {
	producer := p.client.Client().GetProducer(p.client.conf.Namespace, topic)
	res, err := producer.PublishMessage(msgRequest)

	p.log(msgRequest, res, err)
	return err
}

// publishMessageRequest 封装消息
func (p *Producer) publishMessageRequest(topic string, msg []byte) mq_http_sdk.PublishMessageRequest {
	msgRequest := mq_http_sdk.PublishMessageRequest{
		MessageBody: string(msg),         //消息内容。
		MessageTag:  "",                  // 消息标签。
		Properties:  map[string]string{}, // 消息属性。
	}

	msgRequest.MessageKey = topic
	msgRequest.Properties["groupId"] = p.client.GetGroupName(topic)
	return msgRequest
}

// log 日志
func (p *Producer) log(msgRequest mq_http_sdk.PublishMessageRequest, msgResponse mq_http_sdk.PublishMessageResponse, err error) {
	p.client.Logger.Errorf("%+v, %+v, %+v", msgRequest, msgResponse, err)
}
