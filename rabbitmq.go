package main

import (
	"context"
	"log"
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
)

func enqueueImage(data []byte, filename string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect rabbimq server")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"workers", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	_ = q

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "image/jpeg",
			Body:        []byte(data),
            Headers: amqp.Table {
                "filename": filename,
                "id": -1,
            },
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", "image")
}
