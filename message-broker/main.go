package main

import (
	"context"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/joho/godotenv"
)

const topic = "example.topic"
const projectID = "test-project"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}

	logger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			ProjectID: projectID,
		},
		logger,
	)
	if err != nil {
		log.Fatalf("error creating a new Subscriber: %v", err.Error())
	}

	messages, err := subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		log.Fatalf("error subscribing: %v", err.Error())
	}

	go process(messages)
	
	publisher, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: projectID,
		},
		logger,
	)
	if err != nil {
		log.Fatalf("error creating a new Publisher: %v", err.Error())
	}

	publishMessages(publisher)
}

func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))
		if err := publisher.Publish(topic, msg); err != nil {
			log.Fatalf("error publishing a topic: %v", err.Error())
		}

		time.Sleep(time.Second)
	}
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		msg.Ack()
	}
}
