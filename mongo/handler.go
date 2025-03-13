package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendResponse(w http.ResponseWriter, statusCode int, response any) {
	jsonRes, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}

type CollectionRes[Entity any] struct {
	Result []Entity `json:"result"`
}

type BookHandler struct {
	mongoClient *mongo.Client
	logger *slog.Logger
}

type Book struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

func NewBookHandler(mongoClient *mongo.Client, logger *slog.Logger) BookHandler {
	return BookHandler{
		mongoClient: mongoClient,
		logger: logger,
	}
}

func (b BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	collection := b.mongoClient.Database("db").Collection("books")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		b.logger.Error("Error collection.Find()", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []Book
	if err = cursor.All(ctx, &result); err != nil {
		b.logger.Error("Error cursor.All()", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SendResponse(w, http.StatusOK, CollectionRes[Book]{Result: result})
}

func (b BookHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		b.logger.Error("Error: Empty ID on URL Param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := b.mongoClient.Database("db").Collection("books")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Error("Error primitive.ObjectIDFromHex()", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result Book
	err = collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: objId}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		b.logger.Error("No documents was found with the ID: " + id)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	SendResponse(w, http.StatusOK, result)
}

func (b BookHandler) CreateOne(w http.ResponseWriter, r *http.Request) {
	var createOneBookReq Book
	err := json.NewDecoder(r.Body).Decode(&createOneBookReq)
	if err != nil {
		b.logger.Error("Error decoding JSON", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	collection := b.mongoClient.Database("db").Collection("books")
	
	result, err := collection.InsertOne(context.Background(), createOneBookReq)
	if err != nil {
		b.logger.Error("Error collection.InsertOne()", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b.logger.Info("Inserted document successfully", slog.Any("id", result.InsertedID))
	w.WriteHeader(http.StatusOK)
}
