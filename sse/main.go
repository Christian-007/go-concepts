package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const TodoCompletedTopic = "todo.completed"
const ProjectID = "test-project"

func main() {
	fmt.Println("SSE")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
	}

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize message broker by userID (for optimizing SSE)
	sseMessageBroker := NewInMemoryMessageBroker()

	// Initialize Google PubSub using Watermill
	watermillLogger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			ProjectID: ProjectID,
		},
		watermillLogger,
	)
	if err != nil {
		log.Fatalf("error creating a new Subscriber: %v", err.Error())
	}

	// Initialize Google PubSub subscriber
	messages, err := subscriber.Subscribe(context.Background(), TodoCompletedTopic)
	if err != nil {
		log.Fatalf("error subscribing: %v", err.Error())
	}

	go func() {
		for msg := range messages {
			userID := string(msg.Payload)
			log.Printf("received message: %s, payload: %s", msg.UUID, userID)
			sseMessageBroker.Publish(userID, Message{
				ID:     msg.UUID,
				Topic:  TodoCompletedTopic,
				Points: 10,
				UserID: userID,
			})
			msg.Ack()
		}
	}()

	// Initialize Google PubSub publisher
	publisher, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: ProjectID,
		},
		watermillLogger,
	)
	if err != nil {
		log.Fatalf("error creating a new Publisher: %v", err.Error())
	}

	// Initialize message broker by topic (mimicking PubSub)
	pubSubMessageBroker := NewInMemoryMessageBroker()

	pubSubSubscriber := pubSubMessageBroker.Subscribe(TodoCompletedTopic)
	go func() {
		for msg := range pubSubSubscriber {
			fmt.Println("[PubSubSubscriber] Received:", msg)
			sseMessageBroker.Publish(msg.UserID, msg)
		}
	}()

	go func() {
		pubSubMessageBroker.Publish(TodoCompletedTopic, Message{
			ID:     uuid.NewString(),
			Topic:  TodoCompletedTopic,
			Points: 10,
			UserID: "abc123",
		})
	}()

	go func() {
		pubSubMessageBroker.Publish(TodoCompletedTopic, Message{
			ID:     uuid.NewString(),
			Topic:  TodoCompletedTopic,
			Points: 30,
			UserID: "abc123",
		})
	}()

	sseSubscriber := sseMessageBroker.Subscribe("abc123")
	go func() {
		for msg := range sseSubscriber {
			fmt.Println("[SseSubscriber] Received: ", msg)
		}
	}()

	// Initialize HTTP Server
	server := &http.Server{
		Addr:         ":4000",
		Handler:      Routes(logger, pubSubMessageBroker, sseMessageBroker, publisher),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting server on :4000")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
		return
	}
}
