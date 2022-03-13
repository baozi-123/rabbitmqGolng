package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMqSimple("demoSimple")
	rabbitmq.ConsumeSimple()
}
