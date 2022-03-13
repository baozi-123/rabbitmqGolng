package main

import (
	RabbitMQ "rabbitmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("demoPubSub")
	for i := 0; i < 100; i++ {
		rabbitmq.PubLishPub("订阅消息生产第" + strconv.Itoa(i))
		time.Sleep(1 + time.Second)
	}
}
