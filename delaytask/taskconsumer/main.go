package main

import (
	"log"
	"time"

	"github.com/zhiruchen/redis-examples/db"
	"github.com/zhiruchen/redis-examples/delaytask/consumer"
)

func main() {
	if err := db.NewRedisClient(); err != nil {
		log.Fatalln("new redis error: ", err)
	}

	log.Println("start consume...")
	consumer.NewDefaultConsumer(1 * time.Second).Consume(db.RedisClient)
}
