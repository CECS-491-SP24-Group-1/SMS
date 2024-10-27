package amqptests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func TestProducer(t *testing.T) {
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

	//Send x messages to the queue
	messageCount := 10
	prefix := uuid.New()
	for i := 0; i < messageCount; i++ {
		//Create a context for the message
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//Create the message body
		body := fmt.Sprintf("<%s> Message %d", prefix, i)

		//Publish the message
		err := ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)

		//Fail on errors
		if err != nil {
			t.Fatalf("Failed to publish a message: %s", err)
		}

		//Log the message
		fmt.Printf("Sent: `%s`\n", body)
	}
}
