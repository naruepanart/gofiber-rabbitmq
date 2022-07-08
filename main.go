package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid"
	"github.com/streadway/amqp"
)

func main() {
	// Create a new RabbitMQ connection.
	connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	// Add route.
	app.Get("/send", func(c *fiber.Ctx) error {
		// Let's start by opening a channel to our RabbitMQ instance
		// over the connection we have already established
		ch, err := connRabbitMQ.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		// With this channel open, we can then start to interact.
		// With the instance and declare Queues that we can publish and subscribe to.
		_, err = ch.QueueDeclare(
			"TestQueue",
			true,
			true,
			false,
			false,
			nil,
		)
		// Handle any errors if we were unable to create the queue.
		if err != nil {
			return err
		}
		order := shortuuid.New()
		// Attempt to publish a message to the queue.
		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(order),
			},
		)
		if err != nil {
			return err
		}

		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
