package amqptests

import (
	"github.com/rabbitmq/amqp091-go"
)

var QueueName = "7a246403-da6d-4bd3-8800-c217b389b93d"

func QueueDecl(ch *amqp091.Channel, name string) (amqp091.Queue, error) {
	queue, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return queue, err
}

func QueuePeek(ch *amqp091.Channel, name string) (int, error) {
	//Declare the queue using passive declaration
	qPassive, err := ch.QueueDeclarePassive(
		name,  // name of the queue
		false, // durable
		true,  // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return -1, err
	}

	//Get the number of messages in the queue
	return qPassive.Messages, nil
}

//-- Unused
/*
func QueueBind(ch *amqp091.Channel, q amqp091.Queue, name string) error {
	return ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		name,   // exchange
		false,
		nil,
	)
}

func FanExch(ch *amqp091.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,     // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}
*/
