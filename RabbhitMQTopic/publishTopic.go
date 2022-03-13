package main

import (
	"fmt"
	RabbitMQ "rabbitmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rebbitmqOne := RabbitMQ.NewRabbitMQTopic("demoTopic", "demo.one")
	rebbitmqTwo := RabbitMQ.NewRabbitMQTopic("demoTopic", "demo.two")

	for i := 0; i < 10; i++ {
		rebbitmqOne.PubLishTopic("One" + strconv.Itoa(i))
		rebbitmqTwo.PubLishTopic("Two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
