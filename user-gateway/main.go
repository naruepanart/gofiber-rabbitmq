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
	// Create a new Fiber instance.
	app := fiber.New()

	// Create a new RabbitMQ connection.
	connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@rabbitmq-management-alpine:5672/")
	if err != nil {
		log.Println(err)
	}
	defer connRabbitMQ.Close()

	ch, err := connRabbitMQ.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"publisher", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Println(err)
	}

	app.Get("/send", func(c *fiber.Ctx) error {
		u := User{}
		u.Name = shortuuid.New()
		out, _ := json.Marshal(u)

		// Attempt to publish a message to the queue.
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
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
