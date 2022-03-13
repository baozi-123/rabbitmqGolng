package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rebbitmqTwo := RabbitMQ.NewRabbitMQRouting("demoRouting", "Two")
	rebbitmqTwo.RecieveRouting()
}
