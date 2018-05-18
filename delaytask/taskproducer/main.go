package main

import (
	"log"
	"time"

	"github.com/zhiruchen/redis-examples/db"
	task "github.com/zhiruchen/redis-examples/delaytask"
)

func main() {
	ts := []*task.DTask{
		task.NewDTask("test task payload", 2*time.Second),
		task.NewDTask("test task payload", 6*time.Second),
		task.NewDTask("test task payload", 10*time.Second),
	}

	if err := db.NewRedisClient(); err != nil {
		log.Println(err)
		return
	}

	p := &task.DefaultProducer{}
	for _, t := range ts {
		p.ProduceTask(db.RedisClient, t)
	}
}
