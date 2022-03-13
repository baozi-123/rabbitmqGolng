package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("demoPubSub")
	rabbitmq.RecieveSub()
}
