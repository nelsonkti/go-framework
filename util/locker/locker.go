package locker

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"time"
)

type RedisLock struct {
	Sync  *redsync.Redsync
	Tries redsync.Option
}

type Mutex struct {
	key   string
	redis *RedisLock
}

func NewRedisLock(redis *redis.Client) *RedisLock {
	pool := goredis.NewPool(redis)
	return &RedisLock{
		Sync:  redsync.New(pool),
		Tries: redsync.WithTries(1),
	}
}

func NewMutex(redisLocker *RedisLock) *Mutex {
	return &Mutex{redis: redisLocker}
}

// Lock acquires a distributed lock for the given key and duration.
func (m *Mutex) Lock(key string, duration time.Duration) error {
	m.key = key
	return m.redis.Sync.NewMutex(m.key, redsync.WithExpiry(duration), m.redis.Tries).Lock()
}

// UnLock 解锁
func (m *Mutex) UnLock() (bool, error) {
	return m.redis.Sync.NewMutex(m.key).Unlock()
}
