package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/driver/postgres"

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
	// Database connection
	dial := "host=postgresql user=postgresql password=U9Ni8JJp3LnJYBCR dbname=user port=5432 TimeZone=Asia/Bangkok"
	//dial := "host=localhost user=postgresql password=U9Ni8JJp3LnJYBCR dbname=user port=5432 TimeZone=Asia/Bangkok"
	ConDB, err := gorm.Open(postgres.Open(dial), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ConDB.AutoMigrate(&User{})

	// RabbitMQ connection
	connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@rabbitmq-management-alpine:5672/")
	//connRabbitMQ, err := amqp.Dial("amqp://rabbitmq:mypassword@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// Open a new channel.
	channel, err := connRabbitMQ.Channel()
	if err != nil {
		log.Println(err)
	}
	defer channel.Close()

	// Start delivering queued messages.
	messages, err := channel.Consume(
		"publisher", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Println(err)
	}

	// Open a channel to receive messages.
	forever := make(chan bool)

	user := User{}

	go func() {
		for message := range messages {
			// For example, just show received message in console.
			/* log.Printf("Received message: %s\n", message.Body) */

			json.Unmarshal(message.Body, &user)
			ConDB.Create(&user)
		}
	}()

	// Close the channel.
	<-forever
}
