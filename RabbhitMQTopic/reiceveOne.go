package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rebbitmqOne := RabbitMQ.NewRabbitMQTopic("demoTopic", "demo.*")
	rebbitmqOne.RecieveTopic()
}
