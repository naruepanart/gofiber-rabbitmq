package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := shortuuid.New()
	u.ID = uuid
	return
}

func main() {
	// RabbitMQ connection
	// connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@rabbitmq-management-alpine:5672/")
	connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	app.Get("/send", func(c *fiber.Ctx) error {
		u := User{}
		u.Name = shortuuid.New()
		out, _ := json.Marshal(u)

		// Open a new channel.
		ch, err := connRabbitMQ.Channel()
		if err != nil {
			log.Println(err)
		}
		defer ch.Close()

		// With the instance and declare Queues that we can publish and subscribe to.
		_, err = ch.QueueDeclare(
			"publisher",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println(err)
		}

		// Attempt to publish a message to the queue.
		err = ch.Publish(
			"",          // exchange
			"publisher", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        out,
			},
		)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
