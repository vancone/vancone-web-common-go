package mq

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

var RabbitMqConfig Config

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Init(viper *viper.Viper, handler func(delivery amqp091.Delivery)) {
	err := viper.UnmarshalKey("rabbitmq", &RabbitMqConfig)
	if err != nil {
		log.Println("viper unmarshal err:", err)
		return
	}

	// 连接 RabbitMQ 服务器
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		RabbitMqConfig.Username,
		RabbitMqConfig.Password,
		RabbitMqConfig.Host,
		RabbitMqConfig.Port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明队列，如果队列不存在则会创建
	q, err := ch.QueueDeclare(
		fmt.Sprintf("datable-agent-sync-%s", viper.GetString("agent.id")), // 队列名称
		false, // 持久化
		false, // 自动删除（当最后一个消费者断开连接后）
		false, // 排他性
		false, // 非阻塞
		nil,   // 额外参数
	)
	failOnError(err, "Failed to declare a queue")

	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标签
		true,   // 自动确认
		false,  // 非排他性
		false,  // 不等待
		false,  // 非阻塞
		nil,    // 额外参数
	)
	failOnError(err, "Failed to register a consumer")

	// 用于等待中断信号的通道
	forever := make(chan bool)

	// 启动一个 goroutine 处理接收到的消息
	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	log.Printf(" [*] Waiting for message")

	// 等待中断信号优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Exiting...")
	forever <- true

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
