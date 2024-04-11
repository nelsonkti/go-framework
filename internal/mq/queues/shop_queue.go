package queues

import (
	job2 "go-framework/internal/mq/job"
	"go-framework/util/mq/queue"
)

var _ queue.Queue = (*ShopQueue)(nil)

type ShopQueue struct {
}

func (o *ShopQueue) Topic() string {
	return "test"
}

func (o *ShopQueue) GroupId() string {
	return "GID_test"
}

func (o *ShopQueue) Enqueue() []queue.Job {
	var jobs []queue.Job
	jobs = append(jobs, &job2.ShopJob{})

	return jobs
}
