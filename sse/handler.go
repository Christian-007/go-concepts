package main

import (
	"log/slog"
	"net/http"
)

type SseHandler struct {
	logger *slog.Logger
}

func NewSseHandler(logger *slog.Logger) SseHandler {
	return SseHandler{logger}
}

func (s SseHandler) GetSse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
