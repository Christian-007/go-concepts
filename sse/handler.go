package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SseHandler struct {
	logger              *slog.Logger
	pubsubMessageBroker *InMemoryMessageBroker
	sseMessageBroker    *InMemoryMessageBroker
	publisher           message.Publisher
}

func NewSseHandler(logger *slog.Logger, pubsubMessageBroker *InMemoryMessageBroker, sseMessageBroker *InMemoryMessageBroker, publisher message.Publisher) SseHandler {
	return SseHandler{logger, pubsubMessageBroker, sseMessageBroker, publisher}
}

func (s SseHandler) SimulateRewards(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		http.Error(w, "Empty user ID parameter", http.StatusBadGateway)
		return
	}

	msg := Message{
		ID:     uuid.NewString(),
		Topic:  TodoCompletedTopic,
		Points: 10,
		UserID: userIDParam,
	}
	s.pubsubMessageBroker.Publish(TodoCompletedTopic, msg)

	s.logger.Info("published a message",
		slog.String("ID", msg.ID),
		slog.String("Topic", msg.Topic),
		slog.Int("Points", msg.Points),
		slog.String("UserID", msg.UserID),
	)
	w.WriteHeader(http.StatusOK)
}

func (s SseHandler) SimulateGooglePubSubRewards(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		http.Error(w, "Empty user ID parameter", http.StatusBadGateway)
		return
	}

	msgID := watermill.NewUUID()
	msg := message.NewMessage(msgID, []byte(userIDParam))
	err := s.publisher.Publish(TodoCompletedTopic, msg)
	if err != nil {
		s.logger.Error("[PubSub] failed to publish TodoCompleted", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.logger.Info("published a message",
		slog.String("ID", msgID),
		slog.String("UserID", userIDParam),
	)
	w.WriteHeader(http.StatusOK)
}

func (s SseHandler) GetSse(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	err := rc.SetWriteDeadline(time.Time{})
	if err != nil {
		http.Error(w, "Write deadline unsupported!", http.StatusInternalServerError)
		return
	}

	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		http.Error(w, "Empty user ID parameter", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	subscriber := s.sseMessageBroker.Subscribe(userIDParam)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			s.logger.Info("Client disconnected")
			return
		case msg := <-subscriber:
			fmt.Fprintf(w, "data: {\"points\": %d}\n\n", msg.Points)
			flusher.Flush()
		}
	}
}
