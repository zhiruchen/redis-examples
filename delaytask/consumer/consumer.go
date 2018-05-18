package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"

	"github.com/zhiruchen/redis-examples/delaytask"
)

// Consumer is the consumer interface
type Consumer interface {
	// Consume handle task from redis sorted set
	Consume(*redis.Client)
}

// DefaultConsumer default consumer
type DefaultConsumer struct {
	interval time.Duration
}

// NewDefaultConsumer new default consumer
func NewDefaultConsumer(d time.Duration) *DefaultConsumer {
	return &DefaultConsumer{interval: d}
}

// Consume default consume the delay task
func (c *DefaultConsumer) Consume(client *redis.Client) {
	ticker := time.NewTicker(c.interval)
	done := make(chan struct{})
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.handleTask(client); err != nil {
					log.Println("DefaultConsumer Consumer handleTask error: ", err)
				}
			case <-quit:
				ticker.Stop()
				done <- struct{}{}
				return
			}
		}
	}()

	<-done
}

func (c *DefaultConsumer) handleTask(client *redis.Client) error {
	key := delaytask.DelayTaskListKey
	zs, err := client.ZRangeWithScores(key, 0, 0).Result()
	if err != nil {
		return err
	}

	if len(zs) == 0 {
		return nil
	}

	z := zs[0]
	executeTime := z.Score

	now := time.Now().Unix()
	if now == int64(executeTime) {
		memb := z.Member.(string)
		key := fmt.Sprintf(delaytask.DelayTaskKey, memb)
		payload := client.HGet(key, "Payload").Val()
		log.Printf("payload: %s\n", payload)

		if err := client.ZRem(delaytask.DelayTaskListKey, memb).Err(); err != nil {
			log.Println("redis zrem error: ", err)
		}
	}

	return nil
}
