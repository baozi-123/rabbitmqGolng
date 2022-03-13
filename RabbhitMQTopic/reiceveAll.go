package main

import RabbitMQ "rabbitmq/rabbitmq"

func main() {
	rebbitmqOne := RabbitMQ.NewRabbitMQTopic("demoTopic", "#")
	rebbitmqOne.RecieveTopic()
}
