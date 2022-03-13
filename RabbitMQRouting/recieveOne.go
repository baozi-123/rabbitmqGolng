package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rebbitmqOne := RabbitMQ.NewRabbitMQRouting("demoRouting", "One")
	rebbitmqOne.RecieveRouting()
}
