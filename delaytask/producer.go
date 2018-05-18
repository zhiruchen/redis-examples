package delaytask

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

// Producer is delay task producer
type Producer interface {
	// ProduceTask enqueue a delay task
	ProduceTask(*redis.Client, *DTask) error
}

// DefaultProducer is default task producer
type DefaultProducer struct{}

// ProduceTask create a delay task(redis sorted set)
func (p *DefaultProducer) ProduceTask(client *redis.Client, task *DTask) error {
	pipe := client.Pipeline()
	defer pipe.Close()

	key := fmt.Sprintf(DelayTaskKey, task.ID)
	if err := pipe.HMSet(key, structs.Map(task)).Err(); err != nil {
		return err
	}

	z := redis.Z{Member: task.ID, Score: float64(task.ExecuteTime)}
	if err := pipe.ZAdd(DelayTaskListKey, z).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return pipe.Discard()
	}

	return nil
}
