package main

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func Routes(logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	sseHandler := NewSseHandler(logger)

	r.Get("/api/sse", sseHandler.GetSse)

	return r
}