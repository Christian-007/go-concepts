package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		logger.Error("No .env file found")
		return
	}

	// Get MongoDB URI
	mongoDbUri := os.Getenv("MONGODB_URI")
	if mongoDbUri == "" {
		logger.Error("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		return
	}

	// Connect to MongoDB Atlas
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Connected to MongoDB Atlas!")

	// Initialize HTTP Server
	server := &http.Server{
		Addr: ":4000",
		Handler: Routes(client, logger),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting server on :4000")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
		return
	}
}
