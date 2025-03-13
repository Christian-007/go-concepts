package main

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(mongoClient *mongo.Client, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	bookHandler := NewBookHandler(mongoClient, logger)

	r.Get("/books", bookHandler.GetAll)
	r.Get("/books/{id}", bookHandler.GetOne)
	r.Post("/books", bookHandler.CreateOne)

	return r
}