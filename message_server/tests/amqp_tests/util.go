package amqptests

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"wraith.me/message_server/pkg/amqp"
)

var conn *amqp091.Connection
var ch *amqp091.Channel

var conf = amqp.AMQPConfig{
	Host:     "127.0.0.1",
	Port:     5672,
	Username: "guest",
	Password: "guest",
}

func connect() {
	var err error
	conn, err := amqp.GetInstance().Connect(&conf)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
}

func disconnect() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
