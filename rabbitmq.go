package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func enqueueImage(data []byte, filename string) []byte {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect rabbimq server")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"workers", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	_ = q

	// Declare a callback queue
	callbackQueue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a callback queue")

	// Set up consumer for the callback queue
	msgs, err := ch.Consume(
		callbackQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:   "image/jpeg",
			Body:          []byte(data),
			ReplyTo:       callbackQueue.Name,
			CorrelationId: "unique-correlation-id",
			Headers: amqp.Table{
				"filename": filename,
				"id":       -1,
			},
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", "image")

	// Wait for a response from the worker
	for d := range msgs {
		if d.CorrelationId == "unique-correlation-id" {
            log.Printf("Done")
            return d.Body
		}
	}
    panic("error")
}
