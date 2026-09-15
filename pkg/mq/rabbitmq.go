package mq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/vancone/vancone-web-common-go/pkg/logger"
)

var RabbitMqConfig Config

type MessageHandler func(d amqp091.Delivery)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Init(viper *viper.Viper) (*amqp091.Channel, error) {
	var err error
	err = viper.UnmarshalKey("rabbitmq", &RabbitMqConfig)
	if err != nil {
		logger.Errorf("viper unmarshal err: %v", err)
		return nil, err
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", RabbitMqConfig.Username, RabbitMqConfig.Password, RabbitMqConfig.Host, RabbitMqConfig.Port)

	conn, err := amqp091.Dial(url)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ:", err.Error())
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel:", err.Error())
		return nil, err
	}

	logger.Infof("Connected to RabbitMQ (%s:%d)", RabbitMqConfig.Host, RabbitMqConfig.Port)
	return channel, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Errorf("%s: %s", msg, err)
	}
}

func Subscribe(channel *amqp091.Channel, queueName string, handler MessageHandler) error {
	if channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to declare a queue:", err.Error())
		return err
	}

	msgs, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to register a consumer:", err.Error())
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	logger.Infof("Subscribed to queue: %s", queueName)
	return nil
}

func SubscribeWithTTL(channel *amqp091.Channel, queueName string, ttl int64, handler MessageHandler) error {
	if channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	args := make(amqp091.Table)
	if ttl > 0 {
		args["x-expires"] = ttl
	}

	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		logger.Error("Failed to declare a queue:", err.Error())
		return err
	}

	msgs, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to register a consumer:", err.Error())
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	logger.Infof("Subscribed to queue: %s with TTL: %dms", queueName, ttl)
	return nil
}

func Publish(channel *amqp091.Channel, queueName string, message []byte) error {
	if channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Errorf("Failed to declare a queue: %v", err)
		return err
	}

	err = channel.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		logger.Errorf("Failed to publish a message: %v", err)
		return err
	}

	logger.Debugf("Published message to queue: %s", queueName)
	return nil
}

func PublishString(channel *amqp091.Channel, queueName string, message string) error {
	return Publish(channel, queueName, []byte(message))
}

func PublishJSON(channel *amqp091.Channel, queueName string, message []byte) error {
	if channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Errorf("Failed to declare a queue: %v", err)
		return err
	}

	err = channel.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		logger.Errorf("Failed to publish a message: %v", err)
		return err
	}

	logger.Debugf("Published JSON message to queue: %s", queueName)
	return nil
}

func Close(channel *amqp091.Channel, conn *amqp091.Connection) {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
	logger.Info("RabbitMQ connection closed")
}
