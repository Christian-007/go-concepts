package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const TodoCompletedTopic = "todo.completed"

func main() {
	fmt.Println("SSE")

	// Initialize message broker by topic (mimicking PubSub)
	pubSubMessageBroker := NewInMemoryMessageBroker()

	// Initialize message broker by userID (for optimizing SSE)
	sseMessageBroker := NewInMemoryMessageBroker()

	pubSubSubscriber := pubSubMessageBroker.Subscribe(TodoCompletedTopic)
	go func() {
		for msg := range pubSubSubscriber {
			fmt.Println("[PubSubSubscriber] Received:", msg)
			sseMessageBroker.Publish(msg.UserID, msg)
		}
	}()

	go func() {
		pubSubMessageBroker.Publish(TodoCompletedTopic, Message{
			ID: uuid.NewString(),
			Topic: TodoCompletedTopic,
			Points: 10,
			UserID: "abc123",
		})
	}()

	go func() {
		pubSubMessageBroker.Publish(TodoCompletedTopic, Message{
			ID: uuid.NewString(),
			Topic: TodoCompletedTopic,
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

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize HTTP Server
	server := &http.Server{
		Addr: ":4000",
		Handler: Routes(logger),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting server on :4000")
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
		return
	}
}
