package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//url 格式 amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://root:root@127.0.0.1:5672/local"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//链接信息
	Mqurl string
}

//创建RabbitMQ结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	//创建rabbitmq链接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建链接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

//断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单step:1.创建简单模式下RabbitMQ实例
func NewRabbitMqSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

//简单模式step: .2简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，如果纯在则跳过创建
	//保证队列纯在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否堵塞
		false,
		//额外熟悉
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据exhange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//简单模式step: .3简单模式下消费代码
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，如果纯在则跳过创建
	//保证队列纯在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外熟悉
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//2.接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启动协程处理消息
	go func() {
		for d := range msgs {
			//fmt.Println(d.Body)
			log.Printf("%s", d.Body)
		}
	}()
	log.Printf("等待接受消息")
	<-forever
}

//订阅模式创建rabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	//创建RabbitMQ实例
	return NewRabbitMQ("", exchangeName, "")
}

//订阅模式生产
func (r *RabbitMQ) PubLishPub(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"fanout",   //类型 设置成广播类型
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败")

	//2。发送消息
	r.channel.Publish(
		r.Exchange,
		"",
		//如果为true，根据exhange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

//订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"fanout",   //类型 设置成广播类型
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败")
	//2。试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		true,
		//是否阻塞
		false,
		//额外熟悉
		nil,
	)
	r.failOnErr(err, "创建队列失败")
	//3。绑定队列到 exchange （交换机中）
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub(订阅模式)模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)

	//4.消费消息
	mesages, err := r.channel.Consume(
		q.Name,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启动协程处理消息
	go func() {
		for d := range mesages {
			//fmt.Println(d.Body)
			log.Printf("%s", d.Body)
		}
	}()
	log.Printf("等待接受消息")
	<-forever
}

//路由模式创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

//路由模式下生产消息
func (r *RabbitMQ) PubLishRouting(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"direct",   //类型 设置成direct
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败")
	//2。发送消息
	r.channel.Publish(
		r.Exchange,
		r.Key,
		//如果为true，根据exhange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"direct",   //类型 设置成direct
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败")
	//2。试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		true,
		//是否阻塞
		false,
		//额外熟悉
		nil,
	)
	r.failOnErr(err, "创建队列失败")
	//3。绑定队列到 exchange （交换机中）
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定KEY
		r.Key,
		r.Exchange,
		false,
		nil)
	//4.消费消息
	mesages, err := r.channel.Consume(
		q.Name,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启动协程处理消息
	go func() {
		for d := range mesages {
			//fmt.Println(d.Body)
			log.Printf("%s", d.Body)
		}
	}()
	log.Printf("等待接受消息")
	<-forever
}

//话题模式创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

//话题模式下生产消息
func (r *RabbitMQ) PubLishTopic(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"topic",    //类型 设置成topic
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败11")
	//2。发送消息
	r.channel.Publish(
		r.Exchange,
		r.Key,
		//如果为true，根据exhange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//话题模式接受消息
//要注意key规则
//其中的*用于匹配一个单词，#用于匹配多个单词（可以是零个）
//匹配demo.*表示匹配 demo.hello, 但是demo.hello.one需要用demo.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"topic",    //类型 设置成topic
		true,       //是否持久化
		false,      //是否自动删除
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false, //是否阻塞
		nil,   //额外属性
	)
	r.failOnErr(err, "创建交换机失败")
	//2。试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		true,
		//是否阻塞
		false,
		//额外熟悉
		nil,
	)
	r.failOnErr(err, "创建队列失败")
	//3。绑定队列到 exchange （交换机中）
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定KEY
		r.Key,
		r.Exchange,
		false,
		nil)
	//4.消费消息
	mesages, err := r.channel.Consume(
		q.Name,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启动协程处理消息
	go func() {
		for d := range mesages {
			//fmt.Println(d.Body)
			log.Printf("%s", d.Body)
		}
	}()
	log.Printf("等待接受消息")
	<-forever
}
