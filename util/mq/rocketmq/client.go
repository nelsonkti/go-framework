package rocketmq

import "fmt"

func (c *Client) GetGroupName(topic string) string {
	if c.queues[topic] == nil || c.queues[topic].Topic() == "" {
		return ""
	}

	groupId := c.queues[topic].GroupId()
	if c.conf.Env == "" || c.conf.Env == "production" {
		return groupId
	}

	return fmt.Sprintf("%s_%s", groupId, c.conf.Env)
}
