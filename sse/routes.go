package main

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(logger *slog.Logger, pubsubMessageBroker *InMemoryMessageBroker, sseMessageBroker *InMemoryMessageBroker, publisher message.Publisher) *chi.Mux {
	r := chi.NewRouter()
	sseHandler := NewSseHandler(logger, pubsubMessageBroker, sseMessageBroker, publisher)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Key"},
	}))
	r.Post("/api/simulate-rewards/{userID}", sseHandler.SimulateRewards)
	r.Post("/api/pubsub-simulate-rewards/{userID}", sseHandler.SimulateGooglePubSubRewards)
	r.Get("/api/sse/{userID}", sseHandler.GetSse)

	return r
}
