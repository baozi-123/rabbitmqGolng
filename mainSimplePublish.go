package main

import (
	"fmt"
	RabbitMQ "rabbitmq/rabbitmq"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMqSimple("demoSimple")
	rabbitmq.PublishSimple("第二次测试")
	fmt.Println("发送成功")
}
