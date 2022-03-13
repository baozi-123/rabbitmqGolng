package main

import (
	"fmt"
	RabbitMQ "rabbitmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rebbitmqOne := RabbitMQ.NewRabbitMQRouting("demoRouting", "One")
	rebbitmqTwo := RabbitMQ.NewRabbitMQRouting("demoRouting", "Two")

	for i := 0; i < 10; i++ {
		rebbitmqOne.PubLishRouting("One" + strconv.Itoa(i))
		rebbitmqTwo.PubLishRouting("Two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
