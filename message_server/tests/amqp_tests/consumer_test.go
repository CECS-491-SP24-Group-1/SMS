package amqptests

import (
	"fmt"
	"testing"
)

func TestConsumer(t *testing.T) {
	/* COMMON START */
	//Connect to LavinMQ
	connect()
	defer disconnect()

	//Declare a new queue
	q, err := QueueDecl(ch, QueueName)
	if err != nil {
		t.Fatalf("failed to declare a queue %s", err)
	}
	/* COMMON STOP */

	/*
		//Peek at the queue to get the number of messages that can be consumed
		messageCount, err := QueuePeek(ch, q.Name)
		if err != nil {
			t.Fatalf("Failed to peek at queue with error %s", err)
		}
	*/

	//Show the number of messages
	fmt.Printf("Found %d messages to consume from queue %s\n", q.Messages, q.Name)

	//Setup a queue consumer
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to register a consumer: %s", err)
	}

	//Consume from the queue as messages come in
	go func() {
		for msg := range msgs {
			fmt.Printf("Received: `%s`\n", string(msg.Body))
		}
	}()
}
