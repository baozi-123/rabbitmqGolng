package main

import (
	"fmt"
	RabbitMQ "rabbitmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMqSimple("demoSimple")
	rabbitmq.PublishSimple("第二次测试")
	fmt.Println("发送成功")

	for i := 0; i < 100; i++ {
		rabbitmq.PublishSimple("Work测试" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
